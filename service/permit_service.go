package service

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aca/permit/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermitService struct {
	db   *gorm.DB
	auth AuthService
}

func NewPermitService(db *gorm.DB, auth AuthService) *PermitService {
	return &PermitService{db: db, auth: auth}
}

func (s *PermitService) GetIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{
		"title": "Permit",
	})
}

func (s *PermitService) CreatePermitHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "create_permit", gin.H{
		"title": "Create Permit",
	})
}

func (s *PermitService) ListPermitHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "list_permit", gin.H{
		"title": "List Permit",
	})
}

func (s *PermitService) PagePermit(c *gin.Context) {
	c.HTML(http.StatusOK, "generate_permit", gin.H{
		"title": "Permit",
	})
}

func (s *PermitService) UpdatePermitHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "update_permit", gin.H{
		"title": "Update Permit",
	})
}

func (s *PermitService) BacksideHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "backside", gin.H{
		"title": "Update Permit",
	})
}

func (s *PermitService) CreatePermit(c *gin.Context) {
	id := uuid.NewString()
	file, err := c.FormFile("file")
	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/create/permit?error=%s", err.Error()))
		return
	}

	fileExtension := filepath.Ext(file.Filename)
	pathFile := fmt.Sprintf("images/%s%s", id, fileExtension)
	if err = c.SaveUploadedFile(file, pathFile); err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/create/permit?error=%s", err.Error()))
		return
	}

	data := models.Permit{
		PermitID:    uuid.NewString(),
		Name:        c.PostForm("name"),
		Region:      c.PostForm("region"),
		NIK:         c.PostForm("nik"),
		Company:     c.PostForm("company"),
		Departement: c.PostForm("departement"),
		Position:    c.PostForm("position"),
		Image:       pathFile,
		Valid:       c.PostForm("valid"),
		CreatedAt:   time.Now(),
	}

	if err = s.db.Save(&data).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/create/permit?error=%s", err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/list/permit")
}

func (s *PermitService) ListPermit(c *gin.Context) {
	data := make([]models.Permit, 0)
	if err := s.db.Find(&data).Error; err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *PermitService) GetDetailPermit(c *gin.Context) {
	var data models.Permit
	if err := s.db.Where("id = ?", c.Query("id")).Find(&data).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/permit?error=%s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *PermitService) UpdatePermit(c *gin.Context) {
	var data models.Permit
	if err := s.db.Where("id = ?", c.Query("id")).First(&data).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/permit?error=%s", err.Error()))
		return
	}

	imageFilePath := data.Image
	file, err := c.FormFile("file")
	if err == nil {
		fileExtension := filepath.Ext(file.Filename)
		pathFile := fmt.Sprintf("images/%s%s", data.PermitID, fileExtension)
		if err = c.SaveUploadedFile(file, pathFile); err != nil {
			fmt.Println("222")
			c.Redirect(http.StatusFound, fmt.Sprintf("/update/permit?error=%s", err.Error()))
			return
		}

		imageFilePath = pathFile
	}

	permit := models.Permit{
		Name:        c.PostForm("name"),
		Region:      c.PostForm("region"),
		NIK:         c.PostForm("nik"),
		Company:     c.PostForm("company"),
		Departement: c.PostForm("departement"),
		Position:    c.PostForm("position"),
		Image:       imageFilePath,
		Valid:       c.PostForm("valid"),
		CreatedAt:   time.Now(),
	}

	if err = s.db.Where("id = ?", data.ID).Updates(&permit).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/update/permit?error=%s", err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/list/permit")
}

func (s *PermitService) DeletePermit(c *gin.Context) {
	var data models.Permit
	id := c.Query("id")
	if err := s.db.Where("id = ?", id).Find(&data).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/permit?error=%s", err.Error()))
		return
	}

	if err := s.db.Where("id = ?", id).Delete(&data).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/permit?error=%s", err.Error()))
		return
	}

	if err := os.Remove(data.Image); err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/permit?error=%s", err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/list/permit")
}
