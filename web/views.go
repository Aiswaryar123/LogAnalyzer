package web

import (
	database "Log_analyzer/pkg/dbmodels"
	"net/http"
	"strconv"

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

func FilterPaginatedLogs(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "100"))

	offset := page * pageSize

	var body struct {
		Levels     []string `json:"levels"`
		Components []string `json:"components"`
		Hosts      []string `json:"hosts"`
		RequestID  string   `json:"requestId"`
		StartTime  string   `json:"startTime"`
		EndTime    string   `json:"endTime"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filtered, err := database.FilterLogs(
		DB,
		body.Levels,
		body.Components,
		body.Hosts,
		body.RequestID,
		body.StartTime,
		body.EndTime,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	total := len(filtered)

	start := offset
	end := offset + pageSize

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	pageEntries := filtered[start:end]

	c.JSON(http.StatusOK, gin.H{
		"entries": pageEntries,
		"total":   total,
	})
}
