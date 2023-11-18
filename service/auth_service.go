package service

import (
	"encoding/json"
	"fmt"
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
		UserID:   uuid.NewString(),
		Name:     req.Name,
		Password: string(hash),
		Role:     req.Role,
	}

	if err = s.db.Save(&data).Error; err != nil {
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

	token, err := utils.CreateToken(data.UserID, data.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.SetCookie("user-name", data.Name, 86400*100, "/", "", false, false)
	c.SetCookie("token", token, 86400*100, "/", "", false, false)
	c.Set("token", token)

	c.JSON(200, gin.H{
		"error": false,
		"msg":   "success",
		"token": token,
	})
}

func (s *AuthService) Me(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	res, err := utils.ExtractTokenID(token)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var tokenModel models.TokenModel
	b, _ := json.Marshal(res)
	if err = json.Unmarshal(b, &tokenModel); err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(200, tokenModel)
}

func (s *AuthService) Logout(c *gin.Context) {
	s.RemoveCookie("user-name", c)
	s.RemoveCookie("token", c)

	c.Redirect(http.StatusFound, "/login")
}

func (s *AuthService) RemoveCookie(name string, c *gin.Context) {
	expiredCookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(c.Writer, expiredCookie)
}
