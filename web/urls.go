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

	r.LoadHTMLGlob("web/templates/*")

	r.GET("/", ShowFilterPage)
	r.POST("/filter", RunFilter)

	return r
}
