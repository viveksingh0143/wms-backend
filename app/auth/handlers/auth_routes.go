package handlers

import (
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.RouterGroup, authHandler *Handler) {
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/login/email", authHandler.LoginByEMail)
	r.POST("/login/staffid", authHandler.LoginByStaffID)
	r.POST("/login/username", authHandler.LoginByUsername)
	r.POST("/refresh-token", authHandler.RefreshToken)
}
