package handler

import (
	"bot/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CreatePassword handles the creation of a new password
// @Summary Create Password
// @Description Create a new password
// @Tags Password
// @Accept json
// @Produce json
// @Param password body models.Password true "Password data"
// @Success 201 {object} map[string]string "Password created successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Failed to create password"
// @Router /password [post]
func (h *HTTPHandler) CreatePassword(c *gin.Context) {
	var password models.Password
	if err := c.BindJSON(&password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.service.PrService.CreatePassword(password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create password"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Password created successfully"})
}

// GetAllPasswordsByPhone retrieves all passwords by phone number
// @Summary Get All Passwords by Phone
// @Description Retrieve all passwords associated with a phone number
// @Tags Password
// @Produce json
// @Param phone path string true "Phone number"
// @Success 200 {array} models.Password "List of passwords"
// @Failure 400 {object} map[string]string "Phone number is required"
// @Failure 500 {object} map[string]string "Failed to fetch passwords"
// @Router /password/{phone} [get]
func (h *HTTPHandler) GetAllPasswordsByPhone(c *gin.Context) {
	phone := c.Param("phone")
	if phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is required"})
		return
	}

	passwords, err := h.service.PrService.GetAllPasswordsByPhone(phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch passwords"})
		return
	}

	c.JSON(http.StatusOK, passwords)
}

// GetByName retrieves passwords by phone and site
// @Summary Get Passwords by Site
// @Description Retrieve passwords by phone number and site name
// @Tags Password
// @Produce json
// @Param phone query string true "Phone number"
// @Param site query string true "Site name"
// @Success 200 {array} models.Password "List of passwords matching criteria"
// @Failure 400 {object} map[string]string "Phone and site are required"
// @Failure 500 {object} map[string]string "Failed to fetch passwords"
// @Router /password [get]
func (h *HTTPHandler) GetByName(c *gin.Context) {
	phone := strings.TrimSpace(c.Query("phone"))
	site := strings.TrimSpace(c.Query("site"))

	if phone == "" || site == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone and site are required"})
		return
	}
	passwords, err := h.service.PrService.GetByName(phone, site)
	if err != nil {
		if err.Error() == "no passwords found for the given phone number and site name" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch passwords"})
		return
	}

	c.JSON(http.StatusOK, passwords)
}