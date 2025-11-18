package web

import (
	database "Log_analyzer/pkg/dbmodels"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowFilterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Level":     []string{},
		"Component": []string{},
		"Host":      []string{},
		"RequestID": "",
		"Timestamp": "",
	})
}

func RunFilter(c *gin.Context) {

	levels := c.PostFormArray("level")
	components := c.PostFormArray("component")
	hosts := c.PostFormArray("host")

	requestID := c.PostForm("request_id")
	timestamp := c.PostForm("timestamp")

	entries, err := database.FilterLogs(DB, levels, components, hosts, requestID, timestamp)
	if err != nil {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Error":     err.Error(),
			"Level":     levels,
			"Component": components,
			"Host":      hosts,
			"RequestID": requestID,
			"Timestamp": timestamp,
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Entries":   entries,
		"Count":     len(entries),
		"Level":     levels,
		"Component": components,
		"Host":      hosts,
		"RequestID": requestID,
		"Timestamp": timestamp,
	})
}
