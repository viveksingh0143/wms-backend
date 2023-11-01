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
	adminRepository "star-wms/app/admin/repository"
	adminService "star-wms/app/admin/service"
	authHandlers "star-wms/app/auth/handlers"
	authServices "star-wms/app/auth/services"
	"star-wms/app/base/handlers/category"
	"star-wms/app/base/handlers/container"
	machine "star-wms/app/base/handlers/customer"
	customer "star-wms/app/base/handlers/machine"
	"star-wms/app/base/handlers/product"
	"star-wms/app/base/handlers/store"
	baseRepository "star-wms/app/base/repository"
	baseService "star-wms/app/base/service"
	"star-wms/configs"
	"star-wms/core/middlewares"
	"star-wms/core/validation"
	"star-wms/plugins/cache"
)

type AppContainer struct {
	DB                *gorm.DB
	PermissionRepo    adminRepository.PermissionRepository
	PermissionService adminService.PermissionService
	PermissionHandler *permission.Handler
	RoleRepo          adminRepository.RoleRepository
	RoleService       adminService.RoleService
	RoleHandler       *role.Handler
	PlantRepo         adminRepository.PlantRepository
	PlantService      adminService.PlantService
	PlantHandler      *plant.Handler
	UserRepo          adminRepository.UserRepository
	UserService       adminService.UserService
	UserHandler       *user.Handler
	CategoryRepo      baseRepository.CategoryRepository
	CategoryService   baseService.CategoryService
	CategoryHandler   *category.Handler
	ProductRepo       baseRepository.ProductRepository
	ProductService    baseService.ProductService
	ProductHandler    *product.Handler
	StoreRepo         baseRepository.StoreRepository
	StoreService      baseService.StoreService
	StoreHandler      *store.Handler
	ContainerRepo     baseRepository.ContainerRepository
	ContainerService  baseService.ContainerService
	ContainerHandler  *container.Handler
	MachineRepo       baseRepository.MachineRepository
	MachineService    baseService.MachineService
	MachineHandler    *machine.Handler
	CustomerRepo      baseRepository.CustomerRepository
	CustomerService   baseService.CustomerService
	CustomerHandler   *customer.Handler
	AuthService       authServices.AuthService
	AuthHandler       *authHandlers.Handler
	CacheManager      *cache.Manager
}

func NewAppContainer(db *gorm.DB, cacheManager *cache.Manager) *AppContainer {
	permissionRepo := adminRepository.NewPermissionGormRepository(db)
	permissionService := adminService.NewPermissionService(permissionRepo)
	permissionHandler := permission.NewPermissionHandler(permissionService)

	roleRepo := adminRepository.NewRoleGormRepository(db)
	roleService := adminService.NewRoleService(roleRepo)
	roleHandler := role.NewRoleHandler(roleService)

	plantRepo := adminRepository.NewPlantGormRepository(db)
	plantService := adminService.NewPlantService(plantRepo)
	plantHandler := plant.NewPlantHandler(plantService)

	userRepo := adminRepository.NewUserGormRepository(db)
	userService := adminService.NewUserService(userRepo, roleService, plantService)
	userHandler := user.NewUserHandler(userService)

	categoryRepo := baseRepository.NewCategoryGormRepository(db)
	categoryService := baseService.NewCategoryService(categoryRepo)
	categoryHandler := category.NewCategoryHandler(categoryService)

	productRepo := baseRepository.NewProductGormRepository(db)
	productService := baseService.NewProductService(productRepo, categoryService)
	productHandler := product.NewProductHandler(productService)

	storeRepo := baseRepository.NewStoreGormRepository(db)
	storeService := baseService.NewStoreService(storeRepo, categoryService, userService)
	storeHandler := store.NewStoreHandler(storeService)

	containerRepo := baseRepository.NewContainerGormRepository(db)
	containerService := baseService.NewContainerService(containerRepo, storeService, productService)
	containerHandler := container.NewContainerHandler(containerService)

	machineRepo := baseRepository.NewMachineGormRepository(db)
	machineService := baseService.NewMachineService(machineRepo)
	machineHandler := machine.NewMachineHandler(machineService)

	customerRepo := baseRepository.NewCustomerGormRepository(db)
	customerService := baseService.NewCustomerService(customerRepo)
	customerHandler := customer.NewCustomerHandler(customerService)

	authService := authServices.NewAuthService(userRepo, roleService, plantService, cacheManager)
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
		CategoryRepo:      categoryRepo,
		CategoryService:   categoryService,
		CategoryHandler:   categoryHandler,
		ProductRepo:       productRepo,
		ProductService:    productService,
		ProductHandler:    productHandler,
		StoreRepo:         storeRepo,
		StoreService:      storeService,
		StoreHandler:      storeHandler,
		ContainerRepo:     containerRepo,
		ContainerService:  containerService,
		ContainerHandler:  containerHandler,
		MachineRepo:       machineRepo,
		MachineService:    machineService,
		MachineHandler:    machineHandler,
		CustomerRepo:      customerRepo,
		CustomerService:   customerService,
		CustomerHandler:   customerHandler,
		AuthService:       authService,
		AuthHandler:       authHandler,
		CacheManager:      cacheManager,
	}
}

func (receiver *AppContainer) RunServer() {
	router := receiver.setupRouter()
	router.Use(middlewares.PaginationMiddleware())
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
