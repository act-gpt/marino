package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/config/system"
	"github.com/act-gpt/marino/model"

	"github.com/gin-gonic/gin"
)

func GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"start_time":  common.StartTime,
			"system_name": system.Config.SystemName,
			"init":        system.Config.Initialled,
		},
	})
}

func GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    system.Config,
	})
}

type EmainRequest struct {
	Email string `json:"email"`
}

func SendEmailVerification(c *gin.Context) {
	req := EmainRequest{}
	err := c.ShouldBindJSON(&req)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	email := req.Email
	if err := common.Validate.Var(email, "required,email"); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}
	if model.IsEmailAlreadyTaken(email) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "邮箱地址已被占用",
		})
		return
	}
	code := common.GenerateVerificationCode(6)
	common.RegisterVerificationCodeWithKey(email, code, common.EmailVerificationPurpose)
	subject := fmt.Sprintf("%s 邮箱验证", system.Config.SystemName)
	content := fmt.Sprintf("<p>您好，你正在进行 %s 邮箱验证。</p>"+
		"<p>您的验证码为: <strong>%s</strong></p>"+
		"<p>验证码 %d 分钟内有效，如果不是本人操作，请忽略。</p>", system.Config.SystemName, code, common.VerificationValidMinutes)
	err = common.SendEmailByResend(subject, email, content)
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
	})
}

func SendPasswordResetEmail(c *gin.Context) {
	req := EmainRequest{}
	err := c.ShouldBindJSON(&req)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	email := req.Email
	if err := common.Validate.Var(email, "required,email"); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}
	if !model.IsEmailAlreadyTaken(email) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "该邮箱地址未注册",
		})
		return
	}
	code := common.GenerateVerificationString(16)
	common.RegisterVerificationCodeWithKey(email, code, common.PasswordResetPurpose)
	print(c.GetHeader("host"))
	link := fmt.Sprintf("%s/user/reset?email=%s&token=%s", c.GetHeader("host"), email, code)
	subject := fmt.Sprintf("%s密码重置", system.Config.SystemName)
	content := fmt.Sprintf("<p>您好，你正在进行%s密码重置。</p>"+
		"<p>点击<a href='%s'>此处</a>进行密码重置。</p>"+
		"<p>重置链接 %d 分钟内有效，如果不是本人操作，请忽略。</p>", system.Config.SystemName, link, common.VerificationValidMinutes)
	err = common.SendEmailByResend(subject, email, content)
	fmt.Println("error", err.Error())
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
	})
}

type PasswordResetRequest struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func ResetPassword(c *gin.Context) {
	var req PasswordResetRequest
	json.NewDecoder(c.Request.Body).Decode(&req)
	if req.Email == "" || req.Token == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}
	if !common.VerifyCodeWithKey(req.Email, req.Token, common.PasswordResetPurpose, common.VerificationValidMinutes) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "重置链接非法或已过期",
		})
		return
	}
	password := common.GenerateVerificationString(12)
	err := model.ResetUserPasswordByEmail(req.Email, password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	common.DeleteKey(req.Email, common.PasswordResetPurpose)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    password,
	})
}
