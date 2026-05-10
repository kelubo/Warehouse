package middleware

import (
	"net/http"

	"warehouse-management/config"
	"warehouse-management/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorResponse 统一错误响应
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ErrorHandler 统一错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			statusCode := http.StatusInternalServerError

			switch err.Type {
			case gin.ErrorTypeBind:
				statusCode = http.StatusBadRequest
			case gin.ErrorTypeRender:
				statusCode = http.StatusInternalServerError
			default:
				statusCode = http.StatusInternalServerError
			}

			cfg := config.LoadConfig()
			response := ErrorResponse{
				Status:  statusCode,
				Message: "请求处理失败",
			}

			// 开发环境显示详细错误信息
			if cfg.IsDevelopment() {
				response.Error = err.Error()
			}

			// 记录错误日志
			utils.Error("Request error",
				zap.Int("status", statusCode),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("error", err.Error()),
			)

			c.JSON(statusCode, response)
		}
	}
}
