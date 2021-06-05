package routes

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Handlers
	*gin.Engine
	userValidator
}
type Handlers interface {
	renderer
	userHandlers
	mtestHandlers
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
	gin.SetMode(gin.ReleaseMode)

	if err := r.initStatic(); err != nil {
		return err
	}

	authMiddleware, err := r.newMiddleware()
	if err != nil {
		return err
	}
	// Handle the index route
	r.GET("/", r.RenderIndexPage)
	r.GET("/search", r.RenderSearchPage)

	adminRoutes := r.Group("/a")
	{
		adminRoutes.GET("/cabinet", r.RenderAdminPage)
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
		apiRoutes.GET("/admin", authMiddleware.MiddlewareFunc(), r.ValidateAdminRights)
		//Get governments names and ids
		apiRoutes.GET("/governments", r.GetGovernmentsHandlers)
		apiRoutes.POST("/governments", authMiddleware.MiddlewareFunc(), r.AddGovernmentHandler)
		apiRoutes.PUT("/governments", authMiddleware.MiddlewareFunc(), r.EditGovernmentNameHandler)
		apiRoutes.DELETE("/governments", authMiddleware.MiddlewareFunc(), r.RemoveGovernmentHandler)

		//Get regions and edit
		apiRoutes.GET("/regions", r.GetRegionsHandler)
		apiRoutes.PUT("/regions", authMiddleware.MiddlewareFunc(), r.PostEditRegions)

		apiRoutes.GET("/actions", r.GetAdministrativeActionsHandler)
		apiRoutes.POST("/actions", authMiddleware.MiddlewareFunc(), r.AddAdministrativeActionHandler)
		apiRoutes.PUT("/actions", authMiddleware.MiddlewareFunc(), r.EditAdministrativeActionsHandler)
		apiRoutes.DELETE("/actions", authMiddleware.MiddlewareFunc(), r.DeleteAdministrativeActionsHandler)

		//users
		apiRoutes.GET("/users", r.GetUsersHandler)
		//Show and edit view
		apiRoutes.GET("/m/get/:mtest_id", r.GetMTESTHandler)
		apiRoutes.POST("/m/update", authMiddleware.MiddlewareFunc(), r.UpdateMTESTHandler)

		//Creation of view
		apiRoutes.POST("/m/create", authMiddleware.MiddlewareFunc(), r.CreateMTESTHandler)

		//Delete handling
		apiRoutes.POST("/m/delete", authMiddleware.MiddlewareFunc(), r.DeleteMTESTHandler)

		//Regulatory acts uploading and deletion
		apiRoutes.POST("/m/regact/list", r.ActsListHandler)
		apiRoutes.GET("/m/regact", r.ActGetHandler)
		apiRoutes.POST("/m/regact", authMiddleware.MiddlewareFunc(), r.ActUploadHandler)
		apiRoutes.DELETE("/m/regact", authMiddleware.MiddlewareFunc(), r.ActDeleteHandler)

		//executors and group calculations
		//Creation of view
		apiRoutes.POST("/m/executor", authMiddleware.MiddlewareFunc(), r.CreateMTESTExecutorHandler)
		//Delete handling
		apiRoutes.DELETE("/m/executor", authMiddleware.MiddlewareFunc(), r.DeleteExecutorHandler)

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

func (r Router) newMiddleware() (*jwt.GinJWTMiddleware, error) {
	// Use the setUserStatus middleware for every route to set a flag
	// indicating whether the request was from an authenticated user or not
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm: "test zone", Key: []byte("secret key"),
		Timeout: time.Hour, MaxRefresh: time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(string); ok {
				return jwt.MapClaims{"id": v}
			}
			log.Errorf("data not in string format %v", data)
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			if r.PasswordCheck(loginVals.Username, loginVals.Password) {
				return loginVals.Username, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(userId interface{}, c *gin.Context) bool {
			id, ok := userId.(string)
			if !ok {
				return false
			}
			if !r.CheckUserActivation(id) {
				return false
			}

			u := c.Request.URL.String()

			if strings.Contains(u, "admin") || isWriteActions(u, c.Request.Method) {
				return r.CheckUserAdmin(id)
			}

			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"code": code, "message": message})
		},
		TokenLookup: "header:Authorization", TokenHeadName: "Bearer", TimeFunc: time.Now,
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			if cl, ok := claims["id"]; ok {
				if c, ok := cl.(string); ok {
					return c
				}
			}
			log.Errorf("IdentityHandler: id is not string %v", claims)
			return nil
		},
	}
	if err := authMiddleware.MiddlewareInit(); err != nil {
		return nil, err
	}

	return authMiddleware, nil
}

func isWriteActions(u string, m string) bool {
	return (strings.Contains(u, "actions") ||
		strings.Contains(u, "governments") ||
		strings.Contains(u, "regions")) &&
		(m == http.MethodPut ||
			m == http.MethodPost ||
			m == http.MethodDelete)
}
