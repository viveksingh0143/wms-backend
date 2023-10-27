package core

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"star-wms/app/admin/handlers/permission"
	"star-wms/app/admin/handlers/plant"
	"star-wms/app/admin/handlers/role"
	"star-wms/app/admin/handlers/user"
	"star-wms/app/admin/repository"
	"star-wms/app/admin/service"
	authHandlers "star-wms/app/auth/handlers"
	authServices "star-wms/app/auth/services"
	"star-wms/configs"
	"star-wms/core/validation"
)

type AppContainer struct {
	DB                *gorm.DB
	PermissionRepo    repository.PermissionRepository
	PermissionService service.PermissionService
	PermissionHandler *permission.Handler
	RoleRepo          repository.RoleRepository
	RoleService       service.RoleService
	RoleHandler       *role.Handler
	PlantRepo         repository.PlantRepository
	PlantService      service.PlantService
	PlantHandler      *plant.Handler
	UserRepo          repository.UserRepository
	UserService       service.UserService
	UserHandler       *user.Handler
	AuthHandler       *authHandlers.Handler
}

func NewAppContainer(db *gorm.DB) *AppContainer {
	permissionRepo := repository.NewPermissionGormRepository(db)
	permissionService := service.NewPermissionService(permissionRepo)
	permissionHandler := permission.NewPermissionHandler(permissionService)

	roleRepo := repository.NewRoleGormRepository(db)
	roleService := service.NewRoleService(roleRepo)
	roleHandler := role.NewRoleHandler(roleService)

	plantRepo := repository.NewPlantGormRepository(db)
	plantService := service.NewPlantService(plantRepo)
	plantHandler := plant.NewPlantHandler(plantService)

	userRepo := repository.NewUserGormRepository(db)
	userService := service.NewUserService(userRepo, roleService, plantService)
	userHandler := user.NewUserHandler(userService)

	authService := authServices.NewAuthService(userRepo, roleService, plantService)
	authHandler := authHandlers.NewAuthHandler(authService)

	return &AppContainer{
		DB:                db,
		PermissionRepo:    permissionRepo,
		PermissionService: permissionService,
		PermissionHandler: permissionHandler,
		RoleRepo:          roleRepo,
		RoleService:       roleService,
		RoleHandler:       roleHandler,
		PlantRepo:         plantRepo,
		PlantService:      plantService,
		PlantHandler:      plantHandler,
		UserRepo:          userRepo,
		UserService:       userService,
		UserHandler:       userHandler,
		AuthHandler:       authHandler,
	}
}

func (receiver *AppContainer) RunServer() {
	router := receiver.setupRouter()
	SetupRoutes(router, receiver)

	log.Info().Msgf("Server started on %s:%d", configs.ServerCfg.Address, configs.ServerCfg.Port)
	err := router.Run(fmt.Sprintf("%s:%d", configs.ServerCfg.Address, configs.ServerCfg.Port))
	if err != nil {
		log.Fatal().Msgf("Failed to start the server: %v", err)
	}
}

func (receiver *AppContainer) setupRouter() *gin.Engine {
	// gin.DisableConsoleColor()
	if configs.AppCfg.Debug {
		log.Info().Msg("Debug mode active for GIN")
		gin.SetMode(gin.DebugMode)
	} else {
		log.Info().Msg("Release mode active for GIN")
		gin.SetMode(gin.ReleaseMode)
	}
	binding.Validator = new(validation.DefaultValidator)
	router := gin.Default()

	// Logging middleware for debugging
	router.Use(func(c *gin.Context) {
		log.Info().Msgf("Incoming request: origin=%s, method=%s, url=%s",
			c.Request.Header.Get("Origin"), c.Request.Method, c.Request.URL)
		c.Next()
	})

	log.Info().Msg("CORS activated for local only (http://localhost:3000) or (http://localhost:5173)")
	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	router.Use(cors.New(config))

	return router
}
