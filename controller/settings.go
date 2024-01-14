package controller

import (
	"net/http"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/model"

	"github.com/gin-gonic/gin"
)

func GetSetting(c *gin.Context) {

	role := c.GetInt("role")
	admin := c.GetString("id")
	id := c.Param("id")
	if role >= common.RoleOperatorUser {
		admin = ""
	}
	setting, err := model.GetSetting(id, admin)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    404,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    0,
		"data":    setting,
	})
}

func UpdateSetting(c *gin.Context) {
	role := c.GetInt("role")
	admin := c.GetString("id")
	id := c.Param("id")
	setting := model.BotSetting{}

	err := c.ShouldBindJSON(&setting)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	if role >= common.RoleOperatorUser {
		admin = ""
	}

	_, err = model.UpdateSetting(id, admin, setting)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    404,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    0,
		"data":    setting,
	})
}
