package web

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupRouter(db *gorm.DB) *gin.Engine {
	DB = db
	r := gin.Default()
	r.SetFuncMap(map[string]interface{}{
		"contains": func(arr []string, val string) bool {
			for _, v := range arr {
				if v == val {
					return true
				}
			}
			return false
		},
	})
	r.Use(CORSMiddleware())
	r.POST("/filter", FilterPaginatedLogs)
	// Load templates
	//r.LoadHTMLGlob("pkg/web/templates/*")
	// Routes
	//r.GET("/", ShowAllLogs)
	// r.POST("/filter", RunFilter)
	// r.GET("/hello", Hello)
	// r.POST("/search", RunFilter1)
	return r
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")
		//disable caching
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
