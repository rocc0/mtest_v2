package main


import (
	"html/template"

	"github.com/gin-gonic/gin"

)

var router *gin.Engine


func main(){
	// Set Gin to production mode
	gin.SetMode(gin.DebugMode)
	// Set the router as the default one provided by Gin
	router = gin.Default()
	// Set static routes
	router.Static( "static/", "assets/",)
	// Set favicon path
	router.StaticFile("/favicon.ico", "static/favicon.ico")
	//Set templates path
	if tmpl, err := template.New("projectViews").Funcs(TemplateHelpers).ParseGlob("templates/*"); err == nil {
		router.SetHTMLTemplate(tmpl)
	} else {
		panic(err)
	}
	//Create admin user
	userInit()
	// Initialize the routes
	initializeRoutes()
	//Search indexing
	elasticIndex()
	// Start serving the application
	router.Run(":8888")
}
