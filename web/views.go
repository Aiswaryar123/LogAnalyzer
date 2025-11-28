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

// func RunFilter(c *gin.Context) {

// 	levels := c.PostFormArray("level")
// 	components := c.PostFormArray("component")
// 	hosts := c.PostFormArray("host")

// 	requestID := c.PostForm("request_id")
// 	timestamp := c.PostForm("timestamp")

// 	entries, err := database.FilterLogs(DB, levels, components, hosts, requestID, timestamp)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"entries": entries,
// 		"count":   len(entries),
// 	})
// }

// if err != nil {
// 	c.HTML(http.StatusOK, "index.html", gin.H{
// 		"Error":     err.Error(),
// 		"Level":     levels,
// 		"Component": components,
// 		"Host":      hosts,
// 		"RequestID": requestID,
// 		"Timestamp": timestamp,
// 	})
// 	return
// }

// 	c.HTML(http.StatusOK, "index.html", gin.H{
// 		"Entries":   entries,
// 		"Count":     len(entries),
// 		"Level":     levels,
// 		"Component": components,
// 		"Host":      hosts,
// 		"RequestID": requestID,
// 		"Timestamp": timestamp,
// 	})

// func Hello(c *gin.Context) {
// 	c.HTML(http.StatusOK, "hello.html", nil)
// }
// func ShowAllLogs(c *gin.Context) {
// 	entries, err := database.GetAllLogs(DB)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"entries": entries[0:10000],
// 	})

// }
func FilterPaginatedLogs(c *gin.Context) {
	//query params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "100"))

	offset := page * pageSize

	//json body
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
