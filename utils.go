package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func GetEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	env := &Env{
		URI:  os.Getenv("URI"),
		DB:   os.Getenv("DB"),
		Addr: os.Getenv("ADDR"),
	}

	return env
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func MatchPasswords(toCheck string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(toCheck))
	return err == nil
}

func RanHash() string {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var slug string

	for i := 0; i < 8; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
		if err != nil {
			fmt.Println("Error generating random index:", err)
			return ""
		}
		slug += string(characters[randomIndex.Int64()])
	}

	return slug
}

func ConvertTitleToHash(title string) string {
	title = strings.ToLower(title)
	title = strings.Replace(title, " ", "-", -1)
	reg, err := regexp.Compile("[^a-zA-Z0-9-]+")

	if err != nil {
		log.Default().Println("Error converting title to hash:", err)
	}

	title = reg.ReplaceAllString(title, "")

	return title
}
