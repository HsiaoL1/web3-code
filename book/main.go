package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	// 启动服务器
	go func() {
		if err := app.Listen(":" + port); err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("Failed to start server")
		}
	}()

	// 优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit // 等待信号

	logrus.Info("Shutting down server...")

	// 设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error shutting down server")
	}

	logrus.Info("Server gracefully stopped")
}
