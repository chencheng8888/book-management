package req

import (
	"book-management/internal/pkg/errcode"
	"book-management/pkg/logger"
	"github.com/gin-gonic/gin"
)

func ParseRequestBody(c *gin.Context, req any) error {
	err := c.ShouldBindBodyWithJSON(req)

	if err != nil {
		logger.LogPrinter.Errorf("Error parsing request: %v", err)
		return errcode.ParseRequestError
	}
	return nil
}

func ParseRequestQuery(c *gin.Context, req any) error {
	err := c.ShouldBindQuery(req)

	if err != nil {
		logger.LogPrinter.Errorf("Error parsing request: %v", err)
		return errcode.ParseRequestError
	}
	return nil
}
