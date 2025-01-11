CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    telegram_id BIGINT UNIQUE NOT NULL,            -- Telegram ID uchun UNIQUE cheklov
    first_name VARCHAR(50),                        -- Ism
    last_name VARCHAR(50),                         -- Familiya
    username VARCHAR(50),                          -- Foydalanuvchi nomi
    phone_number VARCHAR(15) UNIQUE CHECK (phone_number ~ '^\+[1-9][0-9]{7,14}$') -- Telefon raqam formati tekshiruvi
);
CREATE TABLE passwords (
    id SERIAL PRIMARY KEY,             -- Birlamchi kalit
    user_id UUID NOT NULL,             -- Foydalanuvchini aniqlash uchun UUID
    site VARCHAR(255) NOT NULL,        -- Sayt nomi
    password VARCHAR(255) NOT NULL,    -- Parol
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE -- Foydalanuvchilar jadvaliga xorijiy kalit
);
