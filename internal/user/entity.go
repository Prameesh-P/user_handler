package user

type User struct {
	ID uint   `json:"id" gorm:"primaryKey;unique"  `
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`
	LastName string `json:"last_name"  validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"email,required" `
	Age int `json:"age"`
	Phone    string `json:"phone"`
}
