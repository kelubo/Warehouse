package middleware

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateRequest 参数验证中间件
func ValidateRequest(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用反射创建新的实例
		objType := reflect.TypeOf(obj)
		if objType.Kind() == reflect.Ptr {
			objType = objType.Elem()
		}
		newObj := reflect.New(objType).Interface()

		// 绑定请求体
		if err := c.ShouldBind(newObj); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "请求参数绑定失败",
				Error:   err.Error(),
			})
			c.Abort()
			return
		}

		// 验证参数
		if err := validate.Struct(newObj); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "参数验证失败",
				Error:   validationErrors.Error(),
			})
			c.Abort()
			return
		}

		// 将验证后的对象存入 context
		c.Set("validated", newObj)
		c.Next()
	}
}

// GetValidatedData 从 context 获取验证后的数据
func GetValidatedData(c *gin.Context) (interface{}, bool) {
	return c.Get("validated")
}
