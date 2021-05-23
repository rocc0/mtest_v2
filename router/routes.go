package routes

import (
	"html/template"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Handlers
	*gin.Engine
	userValidator
}
type Handlers interface {
	UserCabinetHandler(c *gin.Context)
	RenderIndexPage(c *gin.Context)
	RenderSearchPage(c *gin.Context)
	RenderUserPage(c *gin.Context)
	RenderMTESTPage(c *gin.Context)
	GetMTESTHandler(c *gin.Context)
	CreateMTESTHandler(c *gin.Context)
	UpdateMTESTHandler(c *gin.Context)
	DeleteMTESTHandler(c *gin.Context)
	GetGovernmentsHandlers(c *gin.Context)
	GetRegionsHandler(c *gin.Context)
	CreateMTESTExecutorHandler(c *gin.Context)
	DeleteExecutorHandler(c *gin.Context)
	GetAdministrativeActionsHandler(c *gin.Context)
	ResetPasswordHandler(c *gin.Context)
	EditUserFieldHandler(c *gin.Context)
	PasswordCheckHandler(c *gin.Context)
	SetNewPasswordHandler(c *gin.Context)
	ActivateAccountHandler(c *gin.Context)
	PostEditRegions(c *gin.Context)
	PostEditGovernments(c *gin.Context)
	RegistrationHandler(c *gin.Context)
	ShowEditGovernments(c *gin.Context)
	ShowAdminPage(c *gin.Context)
}

type userValidator interface {
	PasswordCheck(email, password string) bool
	CheckUserActivation(email string) bool
}

func NewRouter(handlers Handlers, validator userValidator) Router {
	return Router{
		Handlers:      handlers,
		Engine:        gin.New(),
		userValidator: validator,
	}
}

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (r *Router) Init() error {
	// Set Gin to production mode
	gin.SetMode(gin.DebugMode)

	if err := r.initStatic(); err != nil {
		return err
	}

	// Use the setUserStatus middleware for every route to set a flag
	// indicating whether the request was from an authenticated user or not
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm: "test zone", Key: []byte("secret key"),
		Timeout: time.Hour, MaxRefresh: time.Hour,
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals Login
			if binderr := c.ShouldBind(&loginVals); binderr != nil {
				return "", jwt.ErrMissingLoginValues
			}

			if r.PasswordCheck(loginVals.Username, loginVals.Password) {
				return loginVals.Username, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(userId interface{}, c *gin.Context) bool {
			return r.CheckUserActivation(userId.(string))
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"code": code, "message": message})
		},
		TokenLookup: "header:Authorization", TokenHeadName: "Bearer", TimeFunc: time.Now,
	}

	// Handle the index route
	r.GET("/", r.RenderIndexPage)
	r.GET("/search", r.RenderSearchPage)

	adminRoutes := r.Group("/admin")
	{
		adminRoutes.GET("/govs/edit", r.ShowEditGovernments)
		adminRoutes.GET("/cabinet", r.ShowAdminPage)
	}

	userRoutes := r.Group("/u")
	{
		userRoutes.GET("/login", r.RenderIndexPage)
		userRoutes.POST("/login", authMiddleware.LoginHandler)
		userRoutes.GET("/register", r.RenderIndexPage)
		userRoutes.POST("/register", r.RegistrationHandler)
		userRoutes.GET("/cabinet", r.RenderUserPage)
		userRoutes.GET("/activate/:hash", r.RenderIndexPage)
		userRoutes.GET("/reset/:hash", r.RenderIndexPage)

	}

	//get mtest page
	r.GET("/track/:mtest_id", r.RenderMTESTPage)

	apiRoutes := r.Group("/api/v.1/")
	{
		//Get goverments names and ids
		apiRoutes.GET("/govs", r.GetGovernmentsHandlers)
		apiRoutes.POST("/govs/save", authMiddleware.MiddlewareFunc(), r.PostEditGovernments)

		//Get regions and edit
		apiRoutes.GET("/regions", r.GetRegionsHandler)
		apiRoutes.POST("/regions/save", authMiddleware.MiddlewareFunc(), r.PostEditRegions)
		apiRoutes.GET("/adm_actions", r.GetAdministrativeActionsHandler)

		//Show and edit view
		apiRoutes.GET("/m/get/:mtest_id", r.GetMTESTHandler)
		apiRoutes.POST("/m/update", authMiddleware.MiddlewareFunc(), r.UpdateMTESTHandler)

		//Creation of view
		apiRoutes.POST("/m/create", authMiddleware.MiddlewareFunc(), r.CreateMTESTHandler)

		//Delete handling
		apiRoutes.POST("/m/delete", authMiddleware.MiddlewareFunc(), r.DeleteMTESTHandler)

		//executors and group calculations
		//Creation of view
		apiRoutes.POST("/m/excreate", authMiddleware.MiddlewareFunc(), r.CreateMTESTExecutorHandler)
		//Delete handling
		apiRoutes.POST("/m/exdelete", authMiddleware.MiddlewareFunc(), r.DeleteExecutorHandler)

		//user
		apiRoutes.GET("/u/cabinet", authMiddleware.MiddlewareFunc(), r.UserCabinetHandler)
		apiRoutes.POST("/u/edituser", authMiddleware.MiddlewareFunc(), r.EditUserFieldHandler)
		apiRoutes.POST("/u/reset/", r.ResetPasswordHandler)
		apiRoutes.GET("/u/reset/:hash", r.PasswordCheckHandler)
		apiRoutes.POST("/u/reset/:hash", r.SetNewPasswordHandler)
		apiRoutes.GET("/u/activate/:hash", r.ActivateAccountHandler)
	}
	return nil
}

func (r *Router) initStatic() error {
	r.Static("static/", "assets/")
	// Set favicon path
	r.StaticFile("/favicon.ico", "static/favicon.ico")
	//Set templates path
	if tmpl, err := template.New("projectViews").Funcs(template.FuncMap{
		"toString": func(s []uint8) string {
			return string(s)
		},
	}).ParseGlob("templates/*"); err == nil {
		r.SetHTMLTemplate(tmpl)
	} else {
		return err
	}
	return nil
}
