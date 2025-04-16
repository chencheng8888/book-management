package req

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type SimpleRequest struct {
	Donate bool `json:"donate" binding:"required"`
}

func TestParseRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/test", func(c *gin.Context) {
		var req SimpleRequest
		if err := ParseRequestBody(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	// 测试用例
	tests := []struct {
		name         string
		requestBody  string
		expectedCode int
	}{
		{
			name:         "donate=false",
			requestBody:  `{"donate": false}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "No donate field is passed",
			requestBody:  `{}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer([]byte(tt.requestBody)))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			//check code
			assert.Equal(t, tt.expectedCode, w.Code, fmt.Sprintf("status code should be %d", tt.expectedCode))
		})
	}
}
