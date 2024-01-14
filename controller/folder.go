package controller

import (
	"net/http"

	"github.com/act-gpt/marino/events"
	"github.com/act-gpt/marino/model"

	"github.com/gin-gonic/gin"
)

func GetFolderByBot(c *gin.Context) {
	orgId := c.GetString("orgId")
	id := c.Param("id")
	folders, err := model.GetFoldersByBotId(id, orgId)
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
		"data":    folders,
	})
}

func GetFolder(c *gin.Context) {
	orgId := c.GetString("orgId")
	id := c.Param("id")
	folder, err := model.GetFolderById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	// forbidden
	if folder.OrgId != orgId {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Forbidden",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    folder,
	})
}

func CreateFolder(c *gin.Context) {
	//id := c.GetString("id")
	orgId := c.GetString("orgId")
	folder := model.Folder{}
	err := c.ShouldBindJSON(&folder)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	folder.OrgId = orgId

	err = folder.Insert()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	events.Emmiter.Emit("folder.created", folder)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    folder,
	})
}

func UpdateFolder(c *gin.Context) {
	orgId := c.GetString("orgId")
	id := c.Param("id")
	old, err := model.GetFolderById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	if old.OrgId != orgId {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Forbidden",
		})
		return
	}
	folder := model.Folder{}
	err = c.ShouldBindJSON(&folder)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	// force
	folder.Id = id
	folder.OrgId = orgId
	err = folder.Update()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	events.Emmiter.Emit("folder.updated", folder)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    folder,
	})
}

func DeleteFolder(c *gin.Context) {
	orgId := c.GetString("orgId")
	id := c.Param("id")
	folder, err := model.GetFolderById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	if folder.OrgId != orgId {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Forbidden",
		})
		return
	}
	folder.Delete()
	events.Emmiter.Emit("folder.deleted", folder)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    folder,
	})
}
