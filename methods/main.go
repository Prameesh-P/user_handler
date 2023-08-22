package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	method1Mutex sync.Mutex
	db           *gorm.DB
	err error
)

type User struct {
	ID   uint
	Name string
}

func method1(waitTime int, wg *sync.WaitGroup, names *[]string) {
	defer wg.Done()

	method1Mutex.Lock()
	defer method1Mutex.Unlock()

	fmt.Println("Executing Method 1...")
	time.Sleep(time.Duration(waitTime) * time.Second)

	var users []User
	db.Find(&users)
	for _, u := range users {
		*names = append(*names, u.Name)
	}
	fmt.Println("Method 1 Completed")
}

func method2(waitTime int, wg *sync.WaitGroup, names *[]string) {
	defer wg.Done()

	fmt.Println("Executing Method 2...")
	time.Sleep(time.Duration(waitTime) * time.Second)

	var users []User
	db.Find(&users)
	for _, u := range users {
		*names = append(*names, u.Name)
	}
	fmt.Println("Method 2 Completed")
}

func methodsHandler(c *gin.Context) {
	var requestData struct {
		Method   int `json:"method"`
		WaitTime int `json:"waitTime"`
	}

	// Parse the JSON request body
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var wg sync.WaitGroup
	var names []string

	if requestData.Method == 1 {
		wg.Add(1)
		go method1(requestData.WaitTime, &wg, &names)
	} else if requestData.Method == 2 {
		wg.Add(1)
		go method2(requestData.WaitTime, &wg, &names)
	}

	wg.Wait()

	// Return the list of user names in the response
	c.JSON(http.StatusOK, gin.H{"userNames": names})
}


func NewPostgresDB(dsn string) (*gorm.DB, error) {
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return db, nil
}
func main() {

	dsn := "user=username password=password dbname=database_name sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	r := gin.Default()
	r.POST("/methods", methodsHandler)

	r.Run(":7000")
}


