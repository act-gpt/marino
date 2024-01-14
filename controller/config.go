package controller

import (
	"net/http"

	"github.com/act-gpt/marino/model"

	"github.com/gin-gonic/gin"
)

func GetGetTemplate(c *gin.Context) {
	lang := c.Query("lang")
	data, err := model.GetTemplates(lang)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    data,
	})
}
