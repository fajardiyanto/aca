package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aca/permit/models"
	"github.com/aca/permit/utils"
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

func (s *PermitService) PageSimper(c *gin.Context) {
	c.HTML(http.StatusOK, "generate_simper", gin.H{
		"title": "Simper",
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

	var dataUser models.Permit
	if err := s.db.Where("name = ?", c.PostForm("name")).First(&dataUser).Error; err == nil {
		fmt.Println(err.Error())
		c.Redirect(http.StatusFound, "/list/permit?error=you have already have permit")
		return
	}

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

	permitID := uuid.NewString()
	var isSimper bool
	if c.PostForm("is_simper") == "on" {
		isSimper = true

		dataVehicles := make([]models.Vehicle, 0)
		types := strings.Split(c.PostForm("type_vehicle"), ",")
		names := strings.Split(c.PostForm("name_vehicle"), ",")

		for i := 0; i < len(types) && i < len(names); i++ {
			dataVehicles = append(dataVehicles, models.Vehicle{
				Number: strconv.Itoa(i + 1),
				Type:   types[i],
				Name:   names[i],
			})
		}

		vehicle, err := json.Marshal(dataVehicles)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		dataSimper := models.Simper{
			PermitID:  permitID,
			Valid:     c.PostForm("simper_valid"),
			Type:      c.PostForm("simper_type"),
			Simpol:    c.PostForm("simpol"),
			NoSimpol:  c.PostForm("no_simpol"),
			BloodType: c.PostForm("blood_type"),
			Vehicle:   string(vehicle),
		}

		if err = s.db.Save(&dataSimper).Error; err != nil {
			c.Redirect(http.StatusFound, fmt.Sprintf("/create/permit?error=%s", err.Error()))
			return
		}
	}

	data := models.Permit{
		PermitID:    permitID,
		Name:        c.PostForm("name"),
		Region:      c.PostForm("region"),
		NIK:         c.PostForm("nik"),
		Company:     c.PostForm("company"),
		Departement: c.PostForm("departement"),
		Position:    c.PostForm("position"),
		Image:       pathFile,
		Valid:       c.PostForm("valid"),
		Type:        c.PostForm("type"),
		Violation:   c.PostForm("violation"),
		IsSimper:    isSimper,
		CreatedAt:   time.Now(),
	}

	if err = s.db.Save(&data).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/create/permit?error=%s", err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/list/permit")
}

func (s *PermitService) ListPermit(c *gin.Context) {
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

	data := make([]models.Permit, 0)
	if tokenModel.Role != "admin" {
		if err := s.db.Where("name = ?", c.Query("name")).Find(&data).Error; err != nil {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
	} else {
		if err := s.db.Find(&data).Error; err != nil {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, data)
}

func (s *PermitService) GetDetailPermit(c *gin.Context) {
	var data models.Permit
	if err := s.db.Where("id = ?", c.Query("id")).First(&data).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/permit?error=%s", err.Error()))
		return
	}

	var responseSimper models.SimperResponse
	if data.IsSimper {
		var dataSimper models.Simper
		if err := s.db.Where("permit_id = ?", data.PermitID).First(&dataSimper).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "data not found",
			})
			return
		}

		var vehicle []models.Vehicle
		json.Unmarshal([]byte(dataSimper.Vehicle), &vehicle)

		responseSimper = models.SimperResponse{
			ID:        dataSimper.ID,
			Valid:     dataSimper.Valid,
			PermitID:  dataSimper.PermitID,
			Type:      dataSimper.Type,
			Simpol:    dataSimper.Simpol,
			NoSimpol:  dataSimper.NoSimpol,
			BloodType: dataSimper.BloodType,
			Vehicle:   vehicle,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"permit": data,
		"simper": responseSimper,
	})
}

func (s *PermitService) GetDetailPermitByName(c *gin.Context) {
	var data models.Permit
	if err := s.db.Where("name = ?", c.Query("name")).First(&data).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "data not found",
		})
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
			c.Redirect(http.StatusFound, fmt.Sprintf("/update/permit?error=%s", err.Error()))
			return
		}

		imageFilePath = pathFile
	}

	var isSimper bool
	if c.PostForm("is_simper") == "on" {
		isSimper = true

		dataVehicles := make([]models.Vehicle, 0)
		types := strings.Split(c.PostForm("type_vehicle"), ",")
		names := strings.Split(c.PostForm("name_vehicle"), ",")

		for i := 0; i < len(types) && i < len(names); i++ {
			dataVehicles = append(dataVehicles, models.Vehicle{
				Number: strconv.Itoa(i + 1),
				Type:   types[i],
				Name:   names[i],
			})
		}

		vehicle, err := json.Marshal(dataVehicles)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		dataSimper := models.Simper{
			Valid:     c.PostForm("simper_valid"),
			Type:      c.PostForm("simper_type"),
			Simpol:    c.PostForm("simpol"),
			NoSimpol:  c.PostForm("no_simpol"),
			BloodType: c.PostForm("blood_type"),
			Vehicle:   string(vehicle),
		}

		if data.IsSimper {
			if err = s.db.Where("permit_id = ?", data.PermitID).Updates(&dataSimper).Error; err != nil {
				c.Redirect(http.StatusFound, fmt.Sprintf("/update/permit?error=%s", err.Error()))
				return
			}
		} else {
			dataSimper.PermitID = data.PermitID
			if err = s.db.Save(&dataSimper).Error; err != nil {
				c.Redirect(http.StatusFound, fmt.Sprintf("/update/permit?error=%s", err.Error()))
				return
			}
		}
	} else {
		var dataPermit models.Permit
		var dataSimper models.Simper
		if err = s.db.Model(&dataPermit).Where("id = ?", data.ID).Update("is_simper", false).Error; err != nil {
			c.Redirect(http.StatusFound, fmt.Sprintf("/update/permit?error=%s", err.Error()))
			return
		}
		if err := s.db.Where("permit_id = ?", data.PermitID).Delete(&dataSimper).Error; err != nil {
			c.Redirect(http.StatusFound, fmt.Sprintf("/update/permit?error=%s", err.Error()))
			return
		}
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
		Type:        c.PostForm("type"),
		Violation:   c.PostForm("violation"),
		IsSimper:    isSimper,
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
	var dataSimper models.Simper
	id := c.Query("id")
	if err := s.db.Where("id = ?", id).Find(&data).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/permit?error=%s", err.Error()))
		return
	}

	if err := s.db.Where("id = ?", id).Delete(&data).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/permit?error=%s", err.Error()))
		return
	}

	if err := s.db.Where("permit_id = ?", data.PermitID).Delete(&dataSimper).Error; err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/permit?error=%s", err.Error()))
		return
	}

	if err := os.Remove(data.Image); err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/list/permit?error=%s", err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/list/permit")
}
