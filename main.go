package main

import (
	"net/http"
	"time"

	"github.com/aca/permit/database"
	"github.com/aca/permit/middleware"
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

	r.GET("/login", svcAuth.LoginHTML)
	r.POST("/login", svcAuth.Login)
	r.POST("/register", svcAuth.Register)
	r.GET("/logout", svcAuth.Logout)

	r.Use(middleware.MainMiddleware)

	r.GET("/home", svcPermit.GetIndex)
	r.GET("/list/permit", svcPermit.ListPermitHTML)
	r.GET("/get/permit", svcPermit.ListPermit)
	r.GET("/create/permit", svcPermit.CreatePermitHTML)
	r.POST("/create/permit", svcPermit.CreatePermit)
	r.GET("/detail/permit", svcPermit.GetDetailPermit)
	r.GET("/delete/permit", svcPermit.DeletePermit)

	r.GET("/generate/permit", svcPermit.PagePermit)

	r.Run(":8088")
}