package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

func (s *AuthService) RegisterEmployee(c *gin.Context) {
	c.HTML(http.StatusOK, "create_employee", gin.H{
		"title": "Tambahkan Karyawan",
	})
}

func (s *AuthService) Register(c *gin.Context) {
	req := models.RegisterModel {
		Email: c.PostForm("email"),
		Name: c.PostForm("name"),
		Password: c.PostForm("password"),
		Role: c.PostForm("role"),
	}

	hash, err := utils.HashedPassword(req.Password)
	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/employee?error=%s", err.Error()))
		return
	}

	data := models.Auth{
		UserID:   uuid.NewString(),
		Name:     req.Name,
		Email: 	  req.Email,
		Password: string(hash),
		Role:     req.Role,
	}

	if err = s.db.Save(&data).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/employee?error=%s", err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/list/employee")
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
	if err := s.db.Where("email = ?", req.Email).Find(&data).Error; err != nil {
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

	var data models.Auth
	if err := s.db.Where("user_id = ?", tokenModel.UserID).Find(&data).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	c.JSON(200, data)
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

func (s *AuthService) GetAllUser(c *gin.Context) {
	data := make([]models.Auth, 0)
	if err := s.db.Where("role != ?", "admin").Find(&data).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	c.JSON(200, data)
}

func (s *AuthService) DeletelUser(c *gin.Context) {
	var data models.Permit
	var dataSimper models.Simper
	var dataUser models.Auth
	id := c.Query("id")

	if err := s.db.Where("user_id = ?", id).Find(&dataUser).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/employee?error=%s", err.Error()))
		return
	}

	if err := s.db.Where("user_id = ?", id).Delete(&dataUser).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/employee?error=%s", err.Error()))
		return
	}

	if err := s.db.Where("name = ?", dataUser.Name).First(&data).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/employee?error=%s", err.Error()))
		return
	}

	if err := s.db.Where("name = ?", dataUser.Name).Delete(&data).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/employee?error=%s", err.Error()))
		return
	}

	if err := s.db.Where("permit_id = ?", data.PermitID).Delete(&dataSimper).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/employee?error=%s", err.Error()))
		return
	}

	if err := os.Remove(data.Image); err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/employee?error=%s", err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/list/employee")
}