package http

import (
	"fmt"
	"net/http"
	_"github.com/Prameesh-P/user-handler/docs"
	"github.com/Prameesh-P/user-handler/grpcs"
	userpb "github.com/Prameesh-P/user-handler/pkg/pb"
	"github.com/gin-gonic/gin"
	ginSwagger"github.com/swaggo/gin-swagger"
	swaggerFiles"github.com/swaggo/files"
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
	router.GET("/docs/*any",ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/users", hs.CreateUserWithHTTP)
	router.GET("/users/:id", hs.GetUserByIDWithHTTP)
	router.PUT("/users",hs.UpdateUserWithHTTP)
	router.DELETE("/users/:id",hs.DeleteUserByIDWithHTTP)
	fmt.Println("HTTP server is running on :8081")
	if err := router.Run(":8081"); err != nil {
		fmt.Printf("Failed to start HTTP server: %v\n", err)
	}
}

type Body struct{
	ID uint   `json:"id" gorm:"primaryKey;unique"  `
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`
	LastName string `json:"last_name"  validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"email,required" `
	Age int `json:"age"`
	Phone    string `json:"phone"`
}

// @Summary Create User with HTTP
// @ID create-user
// @Description user-creation route.
// @Tags Users
// @Accept json
// @Produce json
// @Param user_details body Body true "User details"
// @Success 201
// @Failure 400 
// @Router /users [post]
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



// @Summary Update User with HTTP
// @ID update-user
// @Description user-creation route.
// @Tags Users
// @Accept json
// @Produce json
// @Param user_details body Body true "User details"
// @Success 200
// @Failure 400 
// @Router /users [put]
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



// @Summary Delete User with HTTP
// @ID delete-user
// @Description user-creation route.
// @Tags Users
// @Accept json
// @Produce json
// @Param id query string true "Email address of the user"
// @Success 200
// @Failure 400 
// @Router /users [delete]
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


// @Summary Get by id User with HTTP
// @ID getbyid-user
// @Description get-yb-id route.
// @Tags Users
// @Accept json
// @Produce json
// @Param id query string true "Email address of the user"
// @Success 201
// @Failure 400 
// @Router /users [get]
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

