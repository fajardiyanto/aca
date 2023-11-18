package service

import (
	"fmt"
	"net/http"

	"github.com/aca/permit/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PositionService struct {
	db *gorm.DB
}

func NewPositionService(db *gorm.DB) *PositionService {
	return &PositionService{db: db}
}

func (s *PositionService) ListPositionHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "list_position", gin.H{
		"title": "List Position",
	})
}

func (s *PositionService) CreatePosition(c *gin.Context) {
	data := models.Position{
		Name: c.PostForm("name"),
	}

	if err := s.db.Save(&data).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/position?error=%s", err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/list/position")
}

func (s *PositionService) ListPosition(c *gin.Context) {
	data := make([]models.Position, 0)
	if err := s.db.Find(&data).Error; err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *PositionService) DeletePosition(c *gin.Context) {
	var data models.Position
	id := c.Query("id")
	if err := s.db.Where("id = ?", id).Find(&data).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/position?error=%s", err.Error()))
		return
	}

	if err := s.db.Where("id = ?", id).Delete(&data).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/position?error=%s", err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/list/position")
}
