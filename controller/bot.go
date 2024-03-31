package controller

import (
	"encoding/json"
	"strings"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/config"
	"github.com/act-gpt/marino/config/system"
	"github.com/act-gpt/marino/events"

	"net/http"
	"strconv"

	"github.com/act-gpt/marino/model"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type Template struct {
	Id string `json:"id"`
}

// Paginations
type Pagination struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Size  int `json:"size"`
}

func GetBotByOrgId(c *gin.Context) {
	orgId := c.GetString("orgId")
	p, _ := strconv.Atoi(c.Query("p"))
	if p < 0 {
		p = 0
	}
	bots, err := model.GetAllBotByOrgId(orgId, p*common.ItemsPerPage, common.ItemsPerPage)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    403,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": bots,
	})
}

func GetBostByOwner(c *gin.Context) {
	id := c.GetString("id")

	bots, err := model.GetAllBotByOwner(id)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    403,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bots,
	})
}

func GetBostByAdmin(c *gin.Context) {
	id := c.GetString("id")

	bots, err := model.GetAllBotByAdmin(id)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    403,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bots,
	})
}

func GetBot(c *gin.Context) {
	id := c.Param("id")
	_bot, ok := c.Get("bot")
	if ok {
		bot := _bot.(*model.Bot)
		bot.AccessToken = ""
		bot.Admin = nil
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"code":    0,
			"data":    bot,
		})
		return
	}
	bot, err := model.GetBotById(id)
	if err != nil || bot.Id == "" {
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
		"data":    bot,
	})
}

func CreateBot(c *gin.Context) {

	id := c.GetString("id")
	orgId := c.GetString("orgId")

	template := Template{}
	err := c.ShouldBindJSON(&template)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    "500",
			"message": err.Error(),
		})
		return
	}
	bot := model.Bot{}
	// 配置 config
	conf, err := model.GetTemplate(template.Id)
	if err != nil || conf.Id == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    "404",
			"message": err.Error(),
		})
		return
	}

	// 组织
	org, err := model.GetOrganizationById(orgId)
	if err != nil || org.Id == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    "404",
			"message": err.Error(),
		})
		return
	}

	bot.Id = common.GetUUID()
	bot.OrgId = orgId
	bot.Owner = id
	bot.Name = conf.Setting["name"].(string)
	// 管理员
	bot.Admin = datatypes.JSON("[\"" + id + "\"]")

	// 设置
	bot.Setting = conf.Setting
	bot.Setting["corpus"] = orgId + ":" + bot.Id

	// TODO: setting
	model := config.GetAvailableModel()

	bot.Setting["model"] = model.Name
	err = bot.Insert()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    "500",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bot,
	})
	events.Emmiter.Emit("bot.created", bot)
}

func UpdateBot(c *gin.Context) {
	conf := system.Config
	orgId := c.GetString("orgId")
	admin := c.GetString("id")
	id := c.Param("id")
	// 无数据
	old, err := model.GetBotByAdminAndId(id, admin)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    "500",
			"message": err.Error(),
		})
		return
	}
	// 非当前组织
	if old.OrgId != orgId {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"code":    "403",
			"message": "Forbidden",
		})
		return
	}
	bot := model.Bot{}
	err = c.ShouldBindJSON(&bot)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// force
	bot.Id = id
	bot.OrgId = orgId
	//bot.Config = old.Config

	var setting model.BotSetting
	s, _ := json.Marshal(bot.Setting)
	err = json.Unmarshal(s, &setting)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	link := setting.Link
	model := setting.Model
	if link != "" {
		model = link
	}
	item := config.MODELS[model]

	if item == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Model is not exists",
		})
		return
	}

	co := item.(*config.MODEL)
	owner := co.Owner
	disabled := co.Disabled

	if disabled {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Model is disable",
		})
		return
	}

	if strings.Contains(model, "act-gpt") && conf.ActGpt.AccessKey == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "ACT GPT is not configured",
		})
		return
	} else if owner == "baidu" && (conf.Baidu.ClientId == "" || conf.Baidu.ClientSecret == "") {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Baidu is not configured",
		})
		return
	}

	err = bot.Update()
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
		"data":    bot,
	})
}

func DeleteBot(c *gin.Context) {
	orgId := c.GetString("orgId")
	admin := c.GetString("id")
	id := c.Param("id")
	bot, err := model.GetBotByAdminAndId(id, admin)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	if bot.OrgId != orgId {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"code":    403,
			"message": "Forbidden",
		})
		return
	}
	bot.Delete()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    0,
		"data":    bot,
	})
	events.Emmiter.Emit("bot.deleted", bot)
}

func GetMessagesByBot(c *gin.Context) {
	orgId := c.GetString("orgId")
	admin := c.GetString("id")
	id := c.Param("id")
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	size, _ := strconv.Atoi(c.Query("size"))
	if size <= 0 {
		size = common.ItemsPerPage
	}
	bot, err := model.GetBotByAdminAndId(id, admin)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    "404",
			"message": err.Error(),
		})
		return
	}
	if bot.OrgId != orgId {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"code":    "403",
			"message": "Forbidden",
		})
		return
	}
	data, total, err := model.GetMessagesByBot(id, (page-1)*size, size, true)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    "404",
			"message": err.Error(),
		})
		return
	}
	if total > 500 {
		total = 500
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"meta": &Pagination{
			Total: int(total),
			Page:  page,
			Size:  size,
		},
	})
}
