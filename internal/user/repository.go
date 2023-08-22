package user

import (
	"context"
	"encoding/json"
	"strconv"

	//"errors"
	"fmt"
	"sync"
	"time"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *User) error
	GetByID(userID string) (*User, error)
	Update(user *User) (*User,error) 
	Delete(userId int)(string,error)
}

type userRepository struct {
	pgDB    *gorm.DB
    redisDB *redis.Client
	counter uint64 
	mu    sync.Mutex
}

func NewUserRepository(pgDB *gorm.DB, redisDB *redis.Client) UserRepository {
	return &userRepository{
		pgDB: pgDB,
		redisDB: redisDB,
	}
}

func (ur *userRepository) Create(user *User) error {
	ur.mu.Lock()
	defer ur.mu.Unlock()
	ur.counter++
	user.ID=uint(ur.counter)
	if err := ur.pgDB.Create(user).Error; err != nil {
        return fmt.Errorf("error inserting user data into database: %w", err)
    }
    // Cache user data in Redis
    userData, err := json.Marshal(user)
    if err != nil {
        return fmt.Errorf("error marshaling user data: %w", err)
    }
	u :=fmt.Sprintf("user:%d",user.ID)
    err = ur.redisDB.Set(context.Background(),u, userData, 10*time.Minute).Err()
    if err != nil {
        // Handle caching error
        fmt.Println("Error caching user data:", err)
    }

    return nil

}

func (ur *userRepository) GetByID(userID string) (*User, error) {
	ur.mu.Lock()
	defer ur.mu.Unlock()
	var user User
	ctx :=context.Background()
	userId ,err := strconv.Atoi(userID)
	if err!=nil{
		return nil,err
	}
	redisData, _ := ur.redisDB.Get(ctx, fmt.Sprintf("user:%d", userId)).Result()
	if redisData == "" {
		// Data not found in Redis cache, fetch from PostgreSQL
		fmt.Println("data from psql")
		result := ur.pgDB.First(&user,userId)
		if result.Error != nil {
			// fmt.Println("Psql error")
			fmt.Println(result.Error)
			if result.Error == gorm.ErrRecordNotFound {
				return nil, nil // User not found
			}
			return nil, fmt.Errorf("error fetching user data from PostgreSQL: %w", result.Error)
		}

		// Cache user data in Redis
		userData, err := json.Marshal(user)
		if err !=nil{
			return nil,err
		}
		ur.redisDB.Set(ctx, fmt.Sprintf("user:%d", userId), userData, 10*time.Minute)
	} else {
		// Use user data from Redis cache
		fmt.Println("data from redis")
		err := json.Unmarshal([]byte(redisData), &user)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling user data from Redis: %w", err)
		}
	}

	return &user, nil
}
func (ur *userRepository) Update(updateUser *User) (*User,error) {
	ur.mu.Lock()
	defer ur.mu.Unlock()
	var users User
	userID:=updateUser.ID
	result := ur.pgDB.First(&users,userID)
		if result.Error != nil {

			if result.Error == gorm.ErrRecordNotFound {
				return nil,errors.New("user not found")
			}
			return nil,fmt.Errorf("error fetching user data from PostgreSQL: %w", result.Error)
		}
		 newUser:=User{
			ID: userID,
			FirstName: updateUser.FirstName,
			LastName: updateUser.LastName,
			Email: updateUser.Email,
			Age: updateUser.Age,
			Phone: updateUser.Phone,
		 }
		if updateUser.FirstName ==""{
			newUser.FirstName=users.FirstName
		}
		if updateUser.LastName==""{
			newUser.LastName=users.LastName
		}
		if updateUser.Email==""{
			newUser.Email=users.Email
		}
		if updateUser.Phone==""{
			newUser.Phone=users.Phone
		}
		if updateUser.Age<=0{
			newUser.Age=users.Age
		}
		if err := ur.pgDB.Save(&newUser).Error; err != nil {
				return nil,err
		}
		userData, err := json.Marshal(newUser)
		if err != nil {
			return nil,fmt.Errorf("error marshaling user data: %w", err)
		}
		u :=fmt.Sprintf("user:%d",newUser.ID)
		err = ur.redisDB.Set(context.Background(),u, userData, 10*time.Minute).Err()
		if err != nil {
			fmt.Println("Error caching user data:", err)
		}
		
	return &newUser,nil
}

func (ur *userRepository) Delete(userId int)(string,error){
	msg := ""
	var user User
	if err := ur.pgDB.First(&user, "id = ?", userId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			msg ="User not found"
			fmt.Println("User not found")
		} else {
		
			fmt.Println("Error:", err)
		}
	}
	if err := ur.pgDB.Delete(&User{}, userId).Error; err != nil {
		return "",err
	}
	ctx := context.Background()
    u := fmt.Sprintf("user:%d",userId)

    // Delete the data associated with the ID from the cache
    if err:=ur.redisDB.Del(ctx,u).Err();err!=nil {
		return "",err
	}
	if msg == ""{
		msg = "User successfully deleted"
	}
	return msg,nil
}