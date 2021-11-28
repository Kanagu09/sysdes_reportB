package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"todolist.go/db"
	"todolist.go/service"
)

const port = 8000

func main() {
	// initialize DB connection
	dsn := db.DefaultDSN(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	if err := db.Connect(dsn); err != nil {
		log.Fatal(err)
	}

	// initialize Gin engine
	engine := gin.Default()
	engine.LoadHTMLGlob("views/*.html")

	// routing
	engine.Static("/assets", "./assets")
	engine.GET("/", service.Home)
	engine.GET("/login", service.LoginPage)
	engine.GET("/register", service.RegisterPage)
	engine.GET("/change_name", service.ChangeNamePage)
	engine.GET("/change_pass", service.ChangePassPage)
	engine.GET("/delete_account", service.DeleteAccountPage)
	engine.GET("/logout", service.Logout)
	engine.GET("/list", service.TaskList)
	engine.GET("/task/:id", service.ShowTask) // ":id" is a parameter
	engine.GET("/edit/:id", service.ShowEdit)
	engine.POST("/login", service.Login)
	engine.POST("/register", service.Register)
	engine.POST("/change_name", service.ChangeName)
	engine.POST("/change_pass", service.ChangePass)
	engine.POST("/delete_account", service.DeleteAccount)
	engine.POST("/list", service.FilteredTaskList)
	engine.POST("/list/add", service.AddTask)
	engine.POST("/list/done/:id", service.DoneTask)
	engine.POST("/list/undone/:id", service.UndoneTask)
	engine.POST("/list/edit/:id", service.EditTask)
	engine.POST("/list/delete/:id", service.DeleteTask)

	// start server
	engine.Run(fmt.Sprintf(":%d", port))
}
