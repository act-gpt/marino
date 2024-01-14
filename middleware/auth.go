package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/config/system"
	"github.com/act-gpt/marino/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func authHelper(c *gin.Context, minRole int) {
	session := sessions.Default(c)
	username := session.Get("username")
	role := session.Get("role")
	id := session.Get("id")
	status := session.Get("status")
	orgId := session.Get("orgId")
	/*
		if os.Getenv("DEMO") != "" {

		}
	*/
	if username == nil {
		// Check access token
		accessToken := c.Request.Header.Get("Authorization")
		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    401,
				"message": "无权进行此操作，未登录且未提供 access token",
			})
			c.Abort()
			return
		}
		user := model.ValidateUserAccessToken(accessToken)
		if user != nil && user.Username != "" {
			// Token is valid
			username = user.Username
			role = user.Role
			id = user.Id
			status = user.Status
			orgId = user.OrgId
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    401,
				"message": "无权进行此操作，access token 无效",
			})
			c.Abort()
			return
		}
	}
	if role.(int) < minRole {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"code":    403,
			"message": "无权进行此操作，权限不足",
		})
		c.Abort()
		return
	}

	if status.(int) == common.UserStatusDisabled {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"code":    6001,
			"message": "用户已被封禁",
		})
		c.Abort()
		return
	}
	c.Set("username", username)
	c.Set("role", role)
	c.Set("id", id)
	c.Set("orgId", orgId)
	c.Next()
}

func OpenAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}
		auth := strings.Replace(c.Request.Header.Get("Authorization"), "Bearer ", "", 1)
		// Authorization header check
		token, _ := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(system.Config.Secret), nil
		})
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    403,
			})
			return
		}

		if bot, ok := model.GetBotById(claims["bot"].(string)); ok == nil {
			c.Set("bot", bot)
			if org, ok := model.GetOrganizationById(bot.OrgId); ok == nil {
				c.Set("org", org)
			}
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    403,
			})
			return
		}
		c.Next()
	}
}

func AdminAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHelper(c, common.RoleAdminUser)
	}
}

func OwnerAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHelper(c, common.RoleOwnerUser)
	}
}

func DistributorAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHelper(c, common.RoleAdminUser)
	}
}

func OperatorAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHelper(c, common.RoleOperatorUser)
	}
}

func RootAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHelper(c, common.RoleRootUser)
	}
}

func TokenAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		access_token := c.Query("access_token")
		if access_token == "" {
			access_token = c.Request.Header.Get("Authorization")
		}
		access_token = strings.Replace(access_token, "Bearer ", "", 1)
		// Bot
		if strings.HasPrefix(access_token, "BA.") {
			id := c.Param("id")
			bot, err := model.BotAccessToken(access_token)
			if err != nil || (id != "" && id != bot.Id) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": gin.H{
						"success": false,
						"code":    403,
					},
				})
				return
			}
			c.Set("bot", bot)
			if org, ok := model.GetOrganizationById(bot.OrgId); ok == nil {
				c.Set("org", org)
			} else {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": gin.H{
						"success": false,
						"code":    403,
					},
				})
				return
			}
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"success": false,
				"code":    403,
			},
		})
		c.Abort()
	}
}
