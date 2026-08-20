//go:build tools

package tools

import (
	_ "github.com/gin-gonic/gin"
	_ "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/joho/godotenv"
	_ "google.golang.org/grpc"
	_ "google.golang.org/protobuf/proto"
	_ "gorm.io/gorm"
)
