package service

import (
	"net/http"
	"time"

	"github.com/aca/permit/models"
	"github.com/aca/permit/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) LoginHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{
		"title": "Permit",
	})
}

func (s *AuthService) Register(c *gin.Context) {
	var req models.LoginModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	hash, err := utils.HashedPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	data := models.Auth{
		UserID: uuid.NewString(),
		Name: req.Name,
		Password: string(hash),
	}

	if err := s.db.Save(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *AuthService) Login(c *gin.Context) {
	var req models.LoginModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	var data models.Auth
	if err := s.db.Where("name = ?", req.Name).Find(&data).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	if err := utils.VerifyPassword(data.Password, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   "Password Doesn't Match",
		})
		return
	}

	token, err := utils.CreateToken(data.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.SetCookie("token", token, 86400*100, "/", "", false, true)
	c.Set("token", token)

	c.JSON(200, gin.H{
		"error": false,
		"msg":   "success",
		"token": token,
	})
}

func (s *AuthService) Logout(c *gin.Context) {
	cookieName := "token"

    expiredCookie := &http.Cookie{
        Name:     cookieName,
        Value:    "",
        Expires:  time.Unix(0, 0),
        MaxAge:   -1,
        HttpOnly: true,
    }

    http.SetCookie(c.Writer, expiredCookie)

	c.Redirect(http.StatusFound, "/login")
}