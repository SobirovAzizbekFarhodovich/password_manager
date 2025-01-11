package handler

import (
	"bot/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getUserID() string {
	return "user-uuid-example"
}

func (h *HTTPHandler) CreatePassword(c *gin.Context) {
	userID := getUserID()
	var password models.Password
	if err := c.BindJSON(&password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := h.service.PrService.CreatePassword(userID, password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create password"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Password created successfully"})
}

func (h *HTTPHandler) GetAllPasswordsByUserID(c *gin.Context) {
	userID := getUserID()
	passwords, err := h.service.PrService.GetAllPasswordsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch passwords"})
		return
	}
	c.JSON(http.StatusOK, passwords)
}

func (h *HTTPHandler) GetByName(c *gin.Context) {
	userID := getUserID()
	site := c.Query("site")
	if site == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Site name is required"})
		return
	}
	passwords, err := h.service.PrService.GetByName(userID, site)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch passwords"})
		return
	}

	c.JSON(http.StatusOK, passwords)
}
