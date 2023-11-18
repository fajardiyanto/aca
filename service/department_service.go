package service

import (
	"fmt"
	"net/http"

	"github.com/aca/permit/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DepartmentSerivce struct {
	db *gorm.DB
}

func NewDepartmentSerivce(db *gorm.DB) *DepartmentSerivce {
	return &DepartmentSerivce{db: db}
}

func (s *DepartmentSerivce) ListDepartmentHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "list_department", gin.H{
		"title": "List Department",
	})
}

func (s *DepartmentSerivce) CreateDepartment(c *gin.Context) {
	data := models.Department{
		Name: c.PostForm("name"),
	}

	if err := s.db.Save(&data).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/department?error=%s", err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/list/department")
}

func (s *DepartmentSerivce) ListDepartment(c *gin.Context) {
	data := make([]models.Department, 0)
	if err := s.db.Find(&data).Error; err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *DepartmentSerivce) DeleteDepartment(c *gin.Context) {
	var data models.Department
	id := c.Query("id")
	if err := s.db.Where("id = ?", id).Find(&data).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/department?error=%s", err.Error()))
		return
	}

	if err := s.db.Where("id = ?", id).Delete(&data).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/department?error=%s", err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/list/department")
}
