package http

import (
	"fmt"
	"net/http"

	"github.com/Prameesh-P/user-handler/grpcs"
	userpb "github.com/Prameesh-P/user-handler/pkg/pb"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	UserClient userpb.UserServiceClient
}

func NewHTTPServer(userClient userpb.UserServiceClient) *HTTPServer {
	return &HTTPServer{
		UserClient: userClient,
	}
}
func (hs *HTTPServer) StartHTTPServer() {
	router := gin.Default()

	router.POST("/users", hs.CreateUserWithHTTP)
	router.GET("/users/:id", hs.GetUserByIDWithHTTP)
	router.PUT("/users",hs.UpdateUserWithHTTP)
	router.DELETE("/users/:id",hs.DeleteUserByIDWithHTTP)
	fmt.Println("HTTP server is running on :8081")
	if err := router.Run(":8081"); err != nil {
		fmt.Printf("Failed to start HTTP server: %v\n", err)
	}
}

func (hs *HTTPServer) CreateUserWithHTTP(c *gin.Context) {
	var user userpb.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	resp,err:=grpcs.CreateUserWithGRPC(grpcs.CreateGRPCConnection(),&user)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}
func (hs *HTTPServer)UpdateUserWithHTTP(c *gin.Context){
	var user  userpb.User
	
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userid := user.ID

	if userid <=0{
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"Enter proper id",
		})
		return
	}
	resp,err:=grpcs.UpdateUserWithHTTP(grpcs.CreateGRPCConnection(),&user)
	if err !=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"message":err.Error(),
		})
		return 
	}
	c.JSON(http.StatusOK,gin.H{
		"success":resp,
	})
	

}
func (hs *HTTPServer)DeleteUserByIDWithHTTP(c *gin.Context){
	userID := c.Param("id")

	msg,err:=grpcs.DeleteUserWithHTTP(grpcs.CreateGRPCConnection(),userID)
	if err != nil{

		c.JSON(http.StatusInternalServerError,gin.H{
			"message":"error from server",
		})
		return 
	}
	if msg.Msg== "User not found"{
		c.JSON(http.StatusBadRequest,gin.H{
			"Message":"User not existing",
		})
		return 
	}
	c.JSON(http.StatusOK,gin.H{
		"succuss":msg.Msg,
	})
}

func (hs *HTTPServer) GetUserByIDWithHTTP(c *gin.Context) {
	userID := c.Param("id")
	resp, err := grpcs.GetUserByIDWithGRPC(grpcs.CreateGRPCConnection(),userID)
	 if err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	if resp.Email == "" {
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"can't find user with this id",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":resp,
	})
}

