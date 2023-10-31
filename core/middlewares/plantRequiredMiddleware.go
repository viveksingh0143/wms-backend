package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"star-wms/app/admin/dto/user"
	"star-wms/app/admin/service"
	"star-wms/core/auth"
	"star-wms/core/common/responses"
)

func PlantRequiredMiddleware(service service.PlantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userForm user.Form
		value, userExists := c.Get(auth.AuthUserKey)
		if userExists {
			var ok bool
			userForm, ok = value.(user.Form)
			if !ok {
				err := errors.New("something went wrong")
				log.Error().Err(err).Msgf("at gin context, user set to %s is not found to be type of user.Form", auth.AuthUserKey)
				rest := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong", err)
				c.JSON(rest.Status, rest)
				c.Abort()
				return
			}
		} else {
			err := errors.New("login user not found in request")
			log.Error().Err(err).Msg("at gin context, no user found")
			rest := responses.NewErrorResponse(http.StatusUnauthorized, "login user not found in request", err)
			c.JSON(rest.Status, rest)
			c.Abort()
			return
		}

		if userForm.Plant == nil || userForm.Plant.ID == 0 {
			err := errors.New("login user is not associated with any plant")
			rest := responses.NewErrorResponse(http.StatusUnauthorized, "Login user is not associated with any plant", err)
			c.JSON(rest.Status, rest)
			c.Abort()
			return
		}

		plantFormPointer := userForm.Plant
		c.Set(auth.AuthPlantKey, *plantFormPointer)
		c.Next()
	}
}
