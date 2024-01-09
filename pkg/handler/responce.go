package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newResponceError(c *gin.Context, log *logrus.Logger, statusCode int, message string) {
	log.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
