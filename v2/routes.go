package v2

import (
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func initializeRoutes() {
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

			if passwordCheck(loginVals.Username, loginVals.Password) {
				return loginVals.Username, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(userId interface{}, c *gin.Context) bool {
			return checkUserActivated(userId.(string))
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"code": code, "message": message})
		},
		TokenLookup: "header:Authorization", TokenHeadName: "Bearer", TimeFunc: time.Now,
	}

	// Handle the index route
	router.GET("/", renderIndexPage)
	router.GET("/search", renderSearchPage)

	adminRoutes := router.Group("/admin")
	{
		adminRoutes.GET("/govs/edit", showEditGovernments)
		adminRoutes.GET("/cabinet", showAdminPage)
	}

	userRoutes := router.Group("/u")
	{
		userRoutes.GET("/login", renderIndexPage)
		userRoutes.POST("/login", authMiddleware.LoginHandler)
		userRoutes.GET("/register", renderIndexPage)
		userRoutes.POST("/register", registrationHandler)
		userRoutes.GET("/cabinet", renderUserPage)
		userRoutes.GET("/activate/:hash", renderIndexPage)
		userRoutes.GET("/reset/:hash", renderIndexPage)

	}

	//get mtest page
	router.GET("/track/:mtest_id", renderMTESTPage)

	apiRoutes := router.Group("/api/v.1/")
	{
		//Get goverments names and ids
		apiRoutes.GET("/govs", getGovernmentsHandlers)
		apiRoutes.POST("/govs/save", authMiddleware.MiddlewareFunc(), postEditGovernments)

		//Get regions and edit
		apiRoutes.GET("/regions", getRegionsHandler)
		apiRoutes.POST("/regions/save", authMiddleware.MiddlewareFunc(), postEditRegions)
		apiRoutes.GET("/adm_actions", getAdministrativeActionsHandler)

		//Show and edit view
		apiRoutes.GET("/m/get/:mtest_id", getMTESTHandler)
		apiRoutes.POST("/m/update", authMiddleware.MiddlewareFunc(), updateMTESTHandler)

		//Creation of view
		apiRoutes.POST("/m/create", authMiddleware.MiddlewareFunc(), createMTESTHandler)

		//Delete handling
		apiRoutes.POST("/m/delete", authMiddleware.MiddlewareFunc(), deleteMTESTHandler)

		//executors and group calculations
		//Creation of view
		apiRoutes.POST("/m/excreate", authMiddleware.MiddlewareFunc(), createMTESTExecutorHandler)
		//Delete handling
		apiRoutes.POST("/m/exdelete", authMiddleware.MiddlewareFunc(), deleteExecutorHandler)

		//user
		apiRoutes.GET("/u/cabinet", authMiddleware.MiddlewareFunc(), userCabinetHandler)
		apiRoutes.POST("/u/edituser", authMiddleware.MiddlewareFunc(), editUserFieldHandler)
		apiRoutes.POST("/u/reset/", resetPasswordHandler)
		apiRoutes.GET("/u/reset/:hash", passwordCheckHandler)
		apiRoutes.POST("/u/reset/:hash", setNewPasswordHandler)
		apiRoutes.GET("/u/activate/:hash", activateAccountHandler)
	}

}
