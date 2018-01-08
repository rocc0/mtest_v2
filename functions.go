package main

import (

	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)




var TemplateHelpers = template.FuncMap {
	"toString": func(s []uint8) string {
		return string(s)
	},
}

func check(e error) {
	if e != nil {
		log.Print(e.Error())
	}
}


func render(c *gin.Context, data gin.H, templateName string) {

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data["pl"])
	case "application/xml":
		// Respond with XML
		c.XML(http.StatusOK, data["pl"])
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}
}