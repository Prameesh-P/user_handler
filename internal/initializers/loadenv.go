package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariable() {

	if err := godotenv.Load();err!=nil{

		log.Fatal(err)
	}
}