package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/keisn1/lenslocked/models"
	"log"
	"os"
	"strconv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})

	err = es.ForgotPassword("jon@calhoun.io", "https://lenslocked.com/reset-pw?token=abc123")
	if err != nil {
		panic(err)
	}
	fmt.Println("Email sent")
}
