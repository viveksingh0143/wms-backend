package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"star-wms/app/admin/dto/user"
	authServices "star-wms/app/auth/services"
	"star-wms/core/auth"
	"star-wms/core/common/responses"
	"strings"
)

func AuthRequiredMiddleware(service authServices.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			err := errors.New("no access token given")
			rest := responses.NewErrorResponse(http.StatusUnauthorized, "Authorization header is missing", err)
			c.JSON(rest.Status, rest)
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userID, _, err := service.ParseAccessToken(tokenString)
		if err != nil {
			rest := responses.NewErrorResponse(http.StatusUnauthorized, "Login session expired", err)
			c.JSON(rest.Status, rest)
			c.Abort()
			return
		}
		var userForm user.Form
		userFormPointer, err := service.GetUserByID(userID)
		if err != nil {
			rest := responses.NewErrorResponse(http.StatusInternalServerError, err.Error(), err)
			c.JSON(rest.Status, rest)
			c.Abort()
			return
		}
		userForm = *userFormPointer
		c.Set(auth.AuthUserKey, userForm)
		c.Next()
	}
}
