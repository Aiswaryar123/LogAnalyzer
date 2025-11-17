package web

import (
	database "Log_analyzer/pkg/dbmodels"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ShowFilterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func RunFilter(c *gin.Context) {
	rawFilter := c.PostForm("filter")

	if strings.TrimSpace(rawFilter) == "" {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Error": "Filter cannot be empty",
		})
		return
	}

	parts := database.SplitUserFilter(rawFilter)

	entries, err := database.QueryDB(DB, parts)
	if err != nil {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Entries": entries,
	})
}
