package main

import (
	"net/http"
	"time"

	"github.com/aca/permit/database"
	"github.com/aca/permit/service"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.MaxMultipartMemory = 1 // 8 MiB
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.HTMLRender = ginview.New(goview.Config{
		Root:         "views",
		Extension:    ".html",
		Master:       "layouts/master",
		Partials:     []string{},
		DisableCache: true,
	})

	r.Static("/assets", "./views/assets")
	r.Static("/images", "./images")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	db := database.ConnectDB()

	svcAuth := service.NewAuthService(db)
	svcPermit := service.NewPermitService(db, *svcAuth)
	svcPosition := service.NewPositionService(db)
	svcDepartment := service.NewDepartmentSerivce(db)

	r.GET("/login", svcAuth.LoginHTML)
	r.POST("/login", svcAuth.Login)
	r.POST("/register", svcAuth.Register)
	r.GET("/logout", svcAuth.Logout)
	r.GET("/me", svcAuth.Me)

	//r.Use(middleware.MainMiddleware)

	r.GET("/home", svcPermit.GetIndex)
	r.GET("/list/permit", svcPermit.ListPermitHTML)
	r.GET("/get/permit", svcPermit.ListPermit)
	r.GET("/create/permit", svcPermit.CreatePermitHTML)
	r.POST("/create/permit", svcPermit.CreatePermit)
	r.GET("/detail/permit", svcPermit.GetDetailPermit)
	r.GET("/delete/permit", svcPermit.DeletePermit)
	r.GET("/update/permit", svcPermit.UpdatePermitHTML)
	r.POST("/update/permit", svcPermit.UpdatePermit)
	r.GET("/backside/permit", svcPermit.BacksideHTML)
	r.GET("/detail/permit/by-name", svcPermit.GetDetailPermitByName)

	r.GET("/generate/permit", svcPermit.PagePermit)

	// position
	r.GET("/list/position", svcPosition.ListPositionHTML)
	r.GET("/get/position", svcPosition.ListPosition)
	r.POST("/create/position", svcPosition.CreatePosition)
	r.GET("/delete/position", svcPosition.DeletePosition)

	r.GET("/list/department", svcDepartment.ListDepartmentHTML)
	r.GET("/get/department", svcDepartment.ListDepartment)
	r.POST("/create/department", svcDepartment.CreateDepartment)
	r.GET("/delete/department", svcDepartment.DeleteDepartment)

	r.Run(":8088")
}
