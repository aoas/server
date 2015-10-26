package routers

import (
	"fmt"

	"github.com/aoas/server/config"
	"github.com/aoas/server/controllers"
	"github.com/aoas/server/middlewares"
	"github.com/aoas/server/models"
	"github.com/aoas/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/itsjamie/gin-cors"
)

func New(engine *xorm.Engine, c config.Config, logger utils.ILogger) (r *gin.Engine, err error) {

	r = gin.New()
	if c.CORS.Enable {
		options := cors.Config{}
		options.Origins = c.CORS.Origins
		options.Methods = c.CORS.Methods
		options.RequestHeaders = c.CORS.Headers
		r.Use(cors.Middleware(options))

	}
	// static
	r.Static("/files", c.File.UploadPath)

	public := r.Group("/api")
	loginMiddleware := middlewares.Auth(c.TokenSecret)

	base := controllers.Base{
		Config: c,
		Engine: engine,
		Logger: logger,
	}

	account := controllers.Account{base}
	public.POST("/login", account.Login)
	public.POST("/register", account.Register)
	public.GET("/tables", func(c *gin.Context) {
		x := models.Engine()
		for _, v := range x.Tables {
			fmt.Println(v.Name)
			for _, c := range v.Columns() {
				fmt.Println(c.Name, c.FieldName, c.MapType)
			}
			fmt.Println("\n")
		}

	})

	userRouter := r.Group("/api/users")
	userRouter.Use(loginMiddleware)

	user := controllers.User{base}
	userRouter.GET("/", user.List)
	userRouter.GET("/:id", user.Get)
	userRouter.POST("/:id/active", user.Active)
	userRouter.GET("/:id/roles", user.Roles)

	// 用token来获取用户, 为防止和用id获取用户信息路由冲突, 固放到跟目录下
	public.GET("/me", loginMiddleware, user.Me)

	// roles
	roleRouter := r.Group("/api/roles")
	roleRouter.Use(loginMiddleware)

	role := controllers.Role{base}
	roleRouter.GET("/", role.List)
	roleRouter.POST("/", role.Create)
	roleRouter.DELETE("/:id", role.Delete)
	roleRouter.POST("/:id/users", role.AddUser)
	roleRouter.GET("/:id/users", role.Users)
	roleRouter.DELETE("/:id/users", role.DeleteUsers)
	roleRouter.GET("/:id/permissions", role.Permissions)
	roleRouter.POST("/:id/permissions", role.AddPermissions)
	roleRouter.DELETE("/:id/permissions", role.DeletePermissions)

	// files
	fileRouter := r.Group("/api/files")
	fileRouter.Use(loginMiddleware)

	file := controllers.File{base}
	fileRouter.GET("/", file.Find)
	fileRouter.POST("/", file.Upload)
	fileRouter.POST("/:id/resize", file.ResizeImage)

	permissions := controllers.PermissionList()
	err = models.AddPermissionsByList(permissions)

	return
}
