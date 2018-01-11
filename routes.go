package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
)

func initializeRoutes() {

	// Use the setUserStatus middleware for every route to set a flag
	// indicating whether the request was from an authenticated user or not

	authMiddleware := &jwt.GinJWTMiddleware {
		Realm:      "test zone",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(eMail string, password string, c *gin.Context) (string, bool) {
			if (loginCheck(eMail, password) == true) {
				return eMail, true
			}

			return eMail, false
		},
		Authorizator: func(userId string, c *gin.Context) bool {
			return true
			//return authCheck(userId)
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},

		TokenLookup: "header:Authorization",

		TokenHeadName: "Bearer",

		TimeFunc: time.Now,
	}

	// Handle the index route
	router.GET("/", showIndexPage) //+

	adminRoutes :=router.Group("/admin")
	{
		adminRoutes.GET("/govs/edit", showEditGovernments)
		adminRoutes.GET("/cabinet", showAdminPage)
	}


	userRoutes := router.Group("/u")
	{
		userRoutes.GET("/login", showIndexPage)//+

		userRoutes.POST("/login", authMiddleware.LoginHandler)

		userRoutes.GET("/register", showIndexPage) //+

		userRoutes.POST("/register", registerHandler) //+

		userRoutes.GET("/cabinet", showUserPage)

	}


	//+ get mtest page
	router.GET("/track/:mtest_id", showMtestPage)//+


	apiRoutes := router.Group("/api/v.1/")

	{
		//Get goverments names and ids
		apiRoutes.GET("/govs", getGovernments) //+
		apiRoutes.POST("/govs/save", authMiddleware.MiddlewareFunc(), postEditGovernments) //+

		//Get regions and edit

		apiRoutes.GET("/regions", getRegions) //+
		apiRoutes.POST("/regions/save", authMiddleware.MiddlewareFunc(), postEditRegions) //+

		apiRoutes.GET("/adm_actions", getAdmActions)

		//Show and edit view
		apiRoutes.GET("/m/get/:mtest_id", getReadMtest) //+

		apiRoutes.POST("/m/update", authMiddleware.MiddlewareFunc(), postUpdateMtest) //+

		//Creation of view
		apiRoutes.POST("/m/create", authMiddleware.MiddlewareFunc(), postCreateMtest) //+

		//Delete handling
		apiRoutes.POST("/m/delete", authMiddleware.MiddlewareFunc(), postDeleteMtest) //+


		//user
		apiRoutes.GET("/u/cabinet", authMiddleware.MiddlewareFunc(), cabinetHandler)
		apiRoutes.POST("/u/edituser",authMiddleware.MiddlewareFunc(), editUserField) //+
	}

}
