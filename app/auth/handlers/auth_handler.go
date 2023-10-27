package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"star-wms/app/admin/dto/ability"
	"star-wms/app/auth/dto"
	authServices "star-wms/app/auth/services"
	"star-wms/core/common/responses"
	"star-wms/core/validation"
)

type Handler struct {
	service authServices.AuthService
}

func NewAuthHandler(s authServices.AuthService) *Handler {
	return &Handler{
		service: s,
	}
}

// Login user
func (ph *Handler) Login(c *gin.Context) {
	var loginForm dto.LoginFormByEMail
	if err := c.ShouldBindJSON(&loginForm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, loginForm))
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(loginForm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, loginForm))
		return
	}

	userForm, err := ph.service.GetUserByEMail(loginForm.EMail, loginForm.Password)
	if err != nil {
		errResponse := responses.NewErrorResponse(http.StatusUnauthorized, "Invalid Username/Password", err)
		if len(errResponse.FieldErrors) > 0 {
			errResponse.Message = errResponse.FieldErrors[0].Message
		}
		c.JSON(http.StatusUnauthorized, errResponse)
		return
	}
	accessToken, err := ph.service.GenerateAccessToken(userForm, loginForm.RememberMe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "something failed to generate access token", err))
		return
	}
	refreshToken, err := ph.service.GenerateRefreshToken(userForm, loginForm.RememberMe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "something failed to generate refresh token", err))
		return
	}

	allAbilities := make([]*ability.Form, 0)
	for _, role := range userForm.Roles {
		allAbilities = append(allAbilities, role.Abilities...)
	}

	loginResponse := &dto.LoginTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Name:         userForm.Name,
		StaffID:      userForm.StaffID,
		Roles:        userForm.Roles,
		Abilities:    allAbilities,
		Plant:        userForm.Plant,
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "User logged-in successfully", loginResponse))
}

// LoginByEMail user
func (ph *Handler) LoginByEMail(c *gin.Context) {
	var loginForm dto.LoginFormByEMail
	if err := c.ShouldBindJSON(&loginForm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, loginForm))
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(loginForm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, loginForm))
		return
	}

	userForm, err := ph.service.GetUserByEMail(loginForm.EMail, loginForm.Password)
	if err != nil {
		errResponse := responses.NewErrorResponse(http.StatusUnauthorized, "Invalid Username/Password", err)
		if len(errResponse.FieldErrors) > 0 {
			errResponse.Message = errResponse.FieldErrors[0].Message
		}
		c.JSON(http.StatusUnauthorized, errResponse)
		return
	}
	accessToken, err := ph.service.GenerateAccessToken(userForm, loginForm.RememberMe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "something failed to generate access token", err))
		return
	}
	refreshToken, err := ph.service.GenerateRefreshToken(userForm, loginForm.RememberMe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "something failed to generate refresh token", err))
		return
	}

	allAbilities := make([]*ability.Form, 0)
	for _, role := range userForm.Roles {
		allAbilities = append(allAbilities, role.Abilities...)
	}

	loginResponse := &dto.LoginTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Name:         userForm.Name,
		StaffID:      userForm.StaffID,
		Roles:        userForm.Roles,
		Abilities:    allAbilities,
		Plant:        userForm.Plant,
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "User logged-in successfully", loginResponse))
}

// LoginByStaffID user
func (ph *Handler) LoginByStaffID(c *gin.Context) {
	var loginForm dto.LoginFormByStaffID
	if err := c.ShouldBindJSON(&loginForm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, loginForm))
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(loginForm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, loginForm))
		return
	}

	userForm, err := ph.service.GetUserByStaffID(loginForm.StaffID, loginForm.Password)
	if err != nil {
		errResponse := responses.NewErrorResponse(http.StatusUnauthorized, "Invalid Username/Password", err)
		if len(errResponse.FieldErrors) > 0 {
			errResponse.Message = errResponse.FieldErrors[0].Message
		}
		c.JSON(http.StatusUnauthorized, errResponse)
		return
	}
	accessToken, err := ph.service.GenerateAccessToken(userForm, loginForm.RememberMe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "something failed to generate access token", err))
		return
	}
	refreshToken, err := ph.service.GenerateRefreshToken(userForm, loginForm.RememberMe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "something failed to generate refresh token", err))
		return
	}

	allAbilities := make([]*ability.Form, 0)
	for _, role := range userForm.Roles {
		allAbilities = append(allAbilities, role.Abilities...)
	}

	loginResponse := &dto.LoginTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Name:         userForm.Name,
		StaffID:      userForm.StaffID,
		Roles:        userForm.Roles,
		Abilities:    allAbilities,
		Plant:        userForm.Plant,
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "User logged-in successfully", loginResponse))
}

// LoginByUsername user
func (ph *Handler) LoginByUsername(c *gin.Context) {
	var loginForm dto.LoginFormByUsername
	if err := c.ShouldBindJSON(&loginForm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, loginForm))
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(loginForm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, loginForm))
		return
	}

	userForm, err := ph.service.GetUserByUsername(loginForm.Username, loginForm.Password)
	if err != nil {
		errResponse := responses.NewErrorResponse(http.StatusUnauthorized, "Invalid Username/Password", err)
		if len(errResponse.FieldErrors) > 0 {
			errResponse.Message = errResponse.FieldErrors[0].Message
		}
		c.JSON(http.StatusUnauthorized, errResponse)
		return
	}
	accessToken, err := ph.service.GenerateAccessToken(userForm, loginForm.RememberMe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "something failed to generate access token", err))
		return
	}
	refreshToken, err := ph.service.GenerateRefreshToken(userForm, loginForm.RememberMe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "something failed to generate refresh token", err))
		return
	}

	allAbilities := make([]*ability.Form, 0)
	for _, role := range userForm.Roles {
		allAbilities = append(allAbilities, role.Abilities...)
	}

	loginResponse := &dto.LoginTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Name:         userForm.Name,
		StaffID:      userForm.StaffID,
		Roles:        userForm.Roles,
		Abilities:    allAbilities,
		Plant:        userForm.Plant,
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "User logged-in successfully", loginResponse))
}

// RefreshToken user
func (ph *Handler) RefreshToken(c *gin.Context) {
	var tokenForm dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&tokenForm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, tokenForm))
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(tokenForm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, tokenForm))
		return
	}

	userID, claims, err := ph.service.ParseRefreshToken(tokenForm.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, responses.NewErrorResponse(http.StatusUnauthorized, "Invalid/Expired refresh token", err))
		return
	}
	userForm, err := ph.service.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, responses.NewErrorResponse(http.StatusUnauthorized, "No user found", err))
		return
	}
	accessToken, err := ph.service.GenerateAccessToken(userForm, claims.ExpireLong)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "something failed to generate access token", err))
		return
	}
	refreshToken, err := ph.service.GenerateRefreshToken(userForm, claims.ExpireLong)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "something failed to generate refresh token", err))
		return
	}

	allAbilities := make([]*ability.Form, 0)
	for _, role := range userForm.Roles {
		allAbilities = append(allAbilities, role.Abilities...)
	}

	loginResponse := &dto.LoginTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Name:         userForm.Name,
		StaffID:      userForm.StaffID,
		Roles:        userForm.Roles,
		Abilities:    allAbilities,
		Plant:        userForm.Plant,
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "User logged-in successfully", loginResponse))
}

// Register user
func (ph *Handler) Register(c *gin.Context) {
	//var loginForm user.Form
	//if err := c.ShouldBindJSON(&loginForm); err != nil {
	//	c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, loginForm))
	//	return
	//}
	//
	//validate := validation.GetValidator()
	//if loginForm.Plant != nil {
	//	if err := validate.Var(loginForm.Plant.ID, "required,gt=0"); err != nil {
	//		iErr := responses.NewInputError("plant.id", "invalid ID for Plant", loginForm.Plant.ID)
	//		c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, "Invalid Plant ID format", iErr))
	//	}
	//}
	//if err := validate.Struct(loginForm); err != nil {
	//	c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, loginForm))
	//	return
	//}
	//
	//err := ph.service.CreateUser(&loginForm)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err))
	//	return
	//}
	c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "Functionality is not designed & developed", nil))
}
