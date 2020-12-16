package logger

import (
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/kuops/go-example-app/server/pkg/log"
	"time"
)

func GinLogger() gin.HandlerFunc {
	logger := log.Logger
	return ginzap.Ginzap(logger,time.RFC3339,true)
}

func GinRecovery() gin.HandlerFunc {
	logger := log.Logger
	return ginzap.RecoveryWithZap(logger,true)
}