package main

import (
	"fmt"
	"time"
	"net/http"
	//"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string
}



var (

    db *gorm.DB
    err error
)

func NewPostgresDB(dsn string) (*gorm.DB, error) {
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return db, nil
}

func init(){
	
}
func method1Wait(waitTime int) {
	time.Sleep(time.Duration(waitTime) * time.Second)
}

func method1(c *gin.Context) {
	var requestData struct {
		Method    int `json:"method"`
		WaitTime int `json:"waitTime"`
	}

	err := c.BindJSON(&requestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	fmt.Println("Method 1 Request Received")

	var users []User
	db.Find(&users)

	fmt.Println("Number of users:", len(users))
	fmt.Println("Names of users:")
	for _, user := range users {
		fmt.Println(user.Name)
	}

	method1Wait(requestData.WaitTime)

	userNames := make([]string, len(users))
	for i, user := range users {
		userNames[i] = user.Name
	}

	c.JSON(http.StatusOK, userNames)
}

func method2Wait(waitTime int) {
	time.Sleep(time.Duration(waitTime) * time.Second)
}

func method2(c *gin.Context) {
	var requestData struct {
		Method    int `json:"method"`
		WaitTime int `json:"waitTime"`
	}

	err := c.BindJSON(&requestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	fmt.Println("Method 2 Request Received")

	var users []User
	db.Find(&users)

	fmt.Println("Number of users:", len(users))
	fmt.Println("Names of users:")
	for _, user := range users {
		fmt.Println(user.Name)
	}

	method2Wait(requestData.WaitTime)

	userNames := make([]string, len(users))
	for i, user := range users {
		userNames[i] = user.Name
	}

	c.JSON(http.StatusOK, userNames)
}

func main() {
	dsn :="host=localhost user=postgres password=pramee-12345 dbname=user_handler port=5432 sslmode=disable"
	NewPostgresDB(dsn)

	router := gin.Default()

	router.POST("/methods1", method1)
	router.POST("/methods2", method2)

	fmt.Println("Server started on :9000")
	router.Run(":9000")
}