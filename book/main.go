package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

func main() {
	// Code
	if err := godotenv.Load(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error loading .env file")
		os.Exit(1)
	}

	var (
		// 读取环境变量
		// 这里可以使用 os.Getenv("ENV_VAR_NAME") 来获取环境变量
		// 例如：port := os.Getenv("PORT")
		port = os.Getenv("PORT")
		app  = fiber.New()
	)

	// Middleware
}
