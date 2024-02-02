package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Todo struct {
	Id        int    `gorm:"primaryKey" json:"id"  `
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var db *gorm.DB

func main() {
	var err error
	e := echo.New()

	dsn := "user=postgres password=password dbname=testdb host=localhost port=5433 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
		panic(err)
	}

	db.AutoMigrate(&Todo{})
	e.Use(middleware.CORS())

	e.POST("/api/data", InsertOne)
	e.GET("/api/data", GetAll)
	e.DELETE("/api/data/:id", DeleteOne)
	e.PUT("/api/data/:id", UpdateOne)

	e.Start(":8080")

}
func InsertOne(c echo.Context) error {

	var todos Todo
	err := c.Bind(&todos)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	fmt.Println(todos.Title)
	db.Create(&todos)

	return c.JSON(http.StatusCreated, todos)

}

func DeleteOne(c echo.Context) error {
	id := c.Param("id")

	var todos Todo

	err := db.Delete(&todos, id).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	reponse := map[string]string{
		"Delete": "success",
	}
	return c.JSON(http.StatusOK, reponse)

}
func GetAll(c echo.Context) error {
	var todos []Todo

	err := db.Find(&todos).Error
	if err != nil {
		return nil
	}
	fmt.Println(todos)

	return c.JSON(http.StatusOK, todos)

}
func UpdateOne(c echo.Context) error {
	id := c.Param("id")

	var todo Todo

	err := c.Bind(&todo)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	result := db.Model(&Todo{}).Where("id = ?", id).Updates(todo)
	if result.Error != nil {
		log.Println(result.Error)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating todo"})
	}

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "No changes made"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Update successful"})
}
