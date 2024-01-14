package controller

import (
	"net/http"
	"strconv"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/events"
	"github.com/act-gpt/marino/model"

	//"strings"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func GetAllOrganization(c *gin.Context) {
	p, _ := strconv.Atoi(c.Query("p"))
	if p < 0 {
		p = 0
	}
	orgs, err := model.GetAllOrganization(p*common.ItemsPerPage, common.ItemsPerPage)
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
		"data":    orgs,
	})
}

func GetOrganization(c *gin.Context) {
	orgId := c.Param("id")
	org, err := model.GetOrganizationById(orgId)
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
		"data":    org,
	})
}

func GetMyOrganization(c *gin.Context) {
	orgId := c.GetString("orgId")
	org, err := model.GetOrganizationById(orgId)
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
		"data":    org,
	})
}

type OrganizationPost struct {
	model.Organization
	Roles int `json:"roles"`
	Size  int `json:"size"`
	Role  int `json:"role"`
}

func CreateOrganization(c *gin.Context) {

	userId := c.GetString("id")
	org := OrganizationPost{}
	err := c.ShouldBindJSON(&org)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	user, err := model.GetUserById(userId, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	// Only one organization pre user current
	if user.OrgId != "" {
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"code":    7000,
				"message": "Over quota",
			})
			return
		}
	}
	org.Id = common.GetUUID()
	user.OrgId = org.Id
	if user.Role < common.RoleOwnerUser {
		user.Role = common.RoleOwnerUser
	}
	err = user.Update(false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	org.Owner = user.Id
	org.Admin = datatypes.JSON("[\"" + user.Id + "\"]")
	err = org.Insert()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	session := sessions.Default(c)
	session.Set("orgId", org.Id)
	session.Set("role", user.Role)
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    org,
	})
	events.Emmiter.Emit("org.created", org, org.Size, org.Roles)
}

func UpdateOrganization(c *gin.Context) {
	id := c.Param("id")
	userId := c.GetString("id")
	org := model.Organization{}
	err := c.ShouldBindJSON(&org)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	get, err := model.GetOrganizationById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	// Owner can change
	if get.Owner != userId {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Forbidden",
		})
		return
	}
	org.Id = id
	err = org.Update()
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
		"data":    org,
	})
}

func DeleteOrganization(c *gin.Context) {
	id := c.Param("id")
	userId := c.GetString("id")
	org, err := model.GetOrganizationById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	if org.Owner != userId {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Forbidden",
		})
		return
	}
	org.Delete()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    org,
	})
}
