package jwt

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		//var code int
		//var data interface{}
		//code = e.SUCCESS

		tokenString := c.GetHeader("Authorization")
		fmt.Println("tokenString:===== ", tokenString)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}
		// 使用密钥对 Token 进行签名
		secretKey := []byte("your-secret-key")
		token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil // 使用同一密钥来解析
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// 如果验证通过，将claims放入上下文中供后续处理函数使用
		claims, ok := token.Claims.(*jwt.StandardClaims)
		if ok && token.Valid {
			c.Set("user_id", claims.Subject)
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token validation failed"})
		c.Abort()

		//token := c.Query("token")
		//if token == "" {
		//	code = e.INVALID_PARAMS
		//} else {
		//	_, err := util.ParseToken(token)
		//	if err != nil {
		//		switch err.(*jwt.ValidationError).Errors {
		//		case jwt.ValidationErrorExpired:
		//			code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
		//		default:
		//			code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
		//		}
		//	}
		//}
		//if code != e.SUCCESS {
		//	c.JSON(http.StatusUnauthorized, gin.H{
		//		"code": code,
		//		"msg":  e.GetMsg(code),
		//		"data": data,
		//	})
		//	c.Abort()
		//	return
		//}
		//
		//c.Next()
	}
}
