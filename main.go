package main

import (
	"context"
	"crypto/rand"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	PhoneNumber string   `bson:"phone_number"`
	Passwords   []Record `bson:"passwords"`
}

type Record struct {
	Name     string `bson:"name"`
	Password string `bson:"password"`
}

var (
	client            *mongo.Client
	userCollection    *mongo.Collection
	verificationCodes = make(map[int64]string)
	registeredUsers   = make(map[int64]string)
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGO_URI")
	token := os.Getenv("TELEGRAM_TOKEN")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	userCollection = client.Database("password_manager_bot").Collection("users")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		text := update.Message.Text

		if _, ok := registeredUsers[chatID]; !ok {
			switch {
			case text == "/start":
				startCommand(bot, chatID)
			case strings.HasPrefix(text, "+"):
				registerUser(bot, chatID, text)
			case len(text) == 6 && isNumeric(text):
				verifyUser(bot, chatID, text)
			default:
				askForRegistration(bot, chatID)
			}
		} else {
			switch {
			case text == "Add Password":
				requestAddPassword(bot, chatID)
			case strings.HasPrefix(text, "Add"):
				addPassword(bot, chatID, text)
			case text == "Retrieve Password":
				retrievePassword(bot, chatID)
			default:
				bot.Send(tgbotapi.NewMessage(chatID, "Unknown command. Use the buttons below to navigate."))
			}
		}
	}
}

func generateCode() string {
	const digits = "0123456789"
	result := ""
	for i := 0; i < 6; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		result += string(digits[n.Int64()])
	}
	return result
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func startCommand(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Welcome! Please send your phone number starting with '+'.")
	bot.Send(msg)
}

func askForRegistration(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Please send your phone number to register.")
	bot.Send(msg)
}

func registerUser(bot *tgbotapi.BotAPI, chatID int64, phone string) {
	code := generateCode()
	verificationCodes[chatID] = code

	log.Printf("Sending SMS to %s: Your verification code is %s", phone, code) // Simulating SMS

	_, err := userCollection.InsertOne(context.TODO(), bson.M{
		"phone_number": phone,
		"passwords":   []Record{},
	})
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Error registering user. Please try again."))
		return
	}

	msg := tgbotapi.NewMessage(chatID, "A verification code has been sent to your phone. Please enter the 6-digit code.")
	bot.Send(msg)
}

func verifyUser(bot *tgbotapi.BotAPI, chatID int64, code string) {
	if verificationCodes[chatID] == code {
		delete(verificationCodes, chatID)
		registeredUsers[chatID] = code
		msg := tgbotapi.NewMessage(chatID, "Verification successful! You can now manage your passwords.")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Add Password"),
				tgbotapi.NewKeyboardButton("Retrieve Password"),
			),
		)
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatID, "Invalid verification code. Please try again.")
		bot.Send(msg)
	}
}

func requestAddPassword(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Please send the password details in the format: Add <name> <password>")
	bot.Send(msg)
}

func addPassword(bot *tgbotapi.BotAPI, chatID int64, text string) {
	parts := strings.SplitN(text, " ", 3)
	if len(parts) < 3 {
		bot.Send(tgbotapi.NewMessage(chatID, "Invalid format. Use: Add <name> <password>"))
		return
	}

	name := parts[1]
	password := parts[2]

	_, err := userCollection.UpdateOne(
		context.TODO(),
		bson.M{"phone_number": registeredUsers[chatID]},
		bson.M{"$push": bson.M{"passwords": Record{Name: name, Password: password}}},
	)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Error adding password. Please try again."))
		return
	}

	bot.Send(tgbotapi.NewMessage(chatID, "Password added successfully."))
}

func retrievePassword(bot *tgbotapi.BotAPI, chatID int64) {
	var user User
	err := userCollection.FindOne(context.TODO(), bson.M{"phone_number": registeredUsers[chatID]}).Decode(&user)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "No passwords found."))
		return
	}

	var response strings.Builder
	response.WriteString("Your saved passwords:\n")
	for _, record := range user.Passwords {
		response.WriteString(record.Name + "\n")
	}

	bot.Send(tgbotapi.NewMessage(chatID, response.String()))
}
