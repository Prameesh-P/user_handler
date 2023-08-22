package user

import (
	"errors"
	"strconv"
)

type UserUsecase interface {
	CreateUser(user *User) (*User, error)
	GetUserByID(userID string) (*User, error)
	UpdateUser(user *User) (*User, error)
	DeleteUser(userId string)(string,error)
}

type userUsecase struct {
	userRepository UserRepository
}

func NewUserUsecase(userRepository UserRepository) UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
	}
}

func (uc *userUsecase) CreateUser(user *User) (*User, error) {
	if user.FirstName == "" || user.Email == "" {
		return nil, errors.New("fistname and email are required")
	}

	if err := uc.userRepository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *userUsecase) GetUserByID(userID string) (*User, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}

	user, err := uc.userRepository.GetByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (uc *userUsecase) UpdateUser(user *User) (*User, error) {
	if user.ID<=0{
		return nil, errors.New("userID is required")
	}
	var users *User
	var err error
	if users,err = uc.userRepository.Update(user); err != nil {
		return nil, err
	}
	return users, nil
}
func  (uc *userUsecase)DeleteUser(userId string)(string,error){
	var msg string
	var err error
	if userId == ""{
		return "",errors.New("User Id is required")
	}
	userid ,err:=strconv.Atoi(userId)
	if err !=nil{
		return "",err
	}
	if msg,err = uc.userRepository.Delete(userid);err!=nil{
		return "",err
	}
	return msg,nil
}