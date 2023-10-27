package services

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"star-wms/app/admin/dto/role"
	"star-wms/app/admin/dto/user"
	"star-wms/app/admin/models"
	"star-wms/app/admin/repository"
	"star-wms/app/admin/service"
	"star-wms/app/auth/dto"
	"star-wms/configs"
	"star-wms/core/common/responses"
	"star-wms/core/types"
	"star-wms/core/utils"
	"strconv"
	"time"
)

type AuthService interface {
	GetUserByID(id uint) (*user.Form, error)
	GetUserByUsername(username string, password string) (*user.Form, error)
	GetUserByEMail(email string, password string) (*user.Form, error)
	GetUserByStaffID(staffID string, password string) (*user.Form, error)
	ToForm(userModel *models.User) *user.Form
	GenerateAccessToken(userForm *user.Form, expireLong bool) (string, error)
	GenerateRefreshToken(userForm *user.Form, expireLong bool) (string, error)
	ParseRefreshToken(refreshToken string) (uint, *dto.CustomRefreshClaims, error)
}

type DefaultAuthService struct {
	repo         repository.UserRepository
	roleService  service.RoleService
	plantService service.PlantService
}

func NewAuthService(repo repository.UserRepository, roleService service.RoleService, plantService service.PlantService) AuthService {
	return &DefaultAuthService{repo: repo, roleService: roleService, plantService: plantService}
}

func (s *DefaultAuthService) GetUserByID(id uint) (*user.Form, error) {
	userData, err := s.repo.GetLoginInfoByField("id", id)
	if err != nil {
		return nil, responses.NewInputError("id", "no user found", nil)
	}
	if userData.Status != types.StatusActive {
		return nil, responses.NewInputError("id", "account is not active", nil)
	}
	return s.ToForm(userData), nil
}

func (s *DefaultAuthService) GetUserByUsername(username string, password string) (*user.Form, error) {
	userData, err := s.repo.GetLoginInfoByField("username", username)
	if err != nil {
		return nil, responses.NewInputError("username", "invalid credentials", nil)
	}
	if err := utils.CompareHashAndPassword(password, userData.Password); err != nil {
		return nil, responses.NewInputError("username", "invalid credentials", nil)
	}
	if userData.Status != types.StatusActive {
		return nil, responses.NewInputError("username", "account is not active", nil)
	}
	return s.ToForm(userData), nil
}

func (s *DefaultAuthService) GetUserByEMail(email string, password string) (*user.Form, error) {
	userData, err := s.repo.GetLoginInfoByField("email", email)
	if err != nil {
		return nil, responses.NewInputError("username", "invalid credentials", nil)
	}
	if err := utils.CompareHashAndPassword(password, userData.Password); err != nil {
		return nil, responses.NewInputError("username", "invalid credentials", nil)
	}
	if userData.Status != types.StatusActive {
		return nil, responses.NewInputError("username", "account is not active", nil)
	}
	return s.ToForm(userData), nil
}

func (s *DefaultAuthService) GetUserByStaffID(staffID string, password string) (*user.Form, error) {
	userData, err := s.repo.GetLoginInfoByField("staff_id", staffID)
	if err != nil {
		return nil, responses.NewInputError("username", "invalid credentials", nil)
	}
	if err := utils.CompareHashAndPassword(password, userData.Password); err != nil {
		return nil, responses.NewInputError("username", "invalid credentials", nil)
	}
	if userData.Status != types.StatusActive {
		return nil, responses.NewInputError("username", "account is not active", nil)
	}
	return s.ToForm(userData), nil
}

func (s *DefaultAuthService) ToForm(userModel *models.User) *user.Form {
	userForm := &user.Form{
		ID:       userModel.ID,
		Name:     userModel.Name,
		StaffID:  userModel.StaffID,
		Username: userModel.Username,
		EMail:    userModel.EMail,
		Password: "",
		Status:   userModel.Status,
	}
	if userModel.Roles != nil {
		roles := make([]*role.Form, 0)
		if len(userModel.Roles) > 0 {
			roles = s.roleService.ToFormSlice(userModel.Roles)
		}
		userForm.Roles = roles
	}
	if userModel.Plant != nil {
		userForm.Plant = s.plantService.ToForm(userModel.Plant)
	}
	return userForm
}

func (s *DefaultAuthService) GenerateAccessToken(userForm *user.Form, expireLong bool) (string, error) {
	var expirationTime time.Time
	if expireLong {
		expirationTime = time.Now().Add(time.Second * time.Duration(configs.AuthCfg.ExpiryLongDuration))
	} else {
		expirationTime = time.Now().Add(time.Second * time.Duration(configs.AuthCfg.ExpiryDuration))
	}
	expiresAt := jwt.NewNumericDate(expirationTime)
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatUint(uint64(userForm.ID), 10),
		ExpiresAt: expiresAt,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(configs.AuthCfg.SecretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s *DefaultAuthService) GenerateRefreshToken(userForm *user.Form, expireLong bool) (string, error) {
	var expirationTime time.Time
	if expireLong {
		expirationTime = time.Now().Add(time.Second * time.Duration(configs.AuthCfg.RefreshExpiryLongDuration))
	} else {
		expirationTime = time.Now().Add(time.Second * time.Duration(configs.AuthCfg.RefreshExpiryDuration))
	}
	expiresAt := jwt.NewNumericDate(expirationTime)

	claims := dto.CustomRefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(userForm.ID), 10),
			ExpiresAt: expiresAt,
		},
		ExpireLong: expireLong,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(configs.AuthCfg.SecretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s *DefaultAuthService) ParseRefreshToken(refreshToken string) (uint, *dto.CustomRefreshClaims, error) {
	claims := &dto.CustomRefreshClaims{}
	tkn, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.AuthCfg.SecretKey), nil
	})
	if err != nil || !tkn.Valid {
		log.Error().Err(err).Msg("Invalid/Expired refresh token")
		return 0, nil, errors.New("given token is Invalid or Expired")
	}
	uint64Val, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Refresh token has invalid user information")
		return 0, nil, errors.New("given refresh token has invalid user information ")
	}

	return uint(uint64Val), claims, nil
}
