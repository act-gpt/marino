package controller

import (
	"net/http"
	"os"
	"strconv"

	"github.com/act-gpt/marino/api"
	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/events"
	"github.com/act-gpt/marino/model"

	//"strings"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"gorm.io/datatypes"
)

type Ids []string

var MARKDOWN = []string{".markdown", ".mdown", ".mkdn", ".mkd", ".mdwn", ".md"}

func UploadKnowledges(c *gin.Context) {
	orgId := c.GetString("orgId")
	user := c.GetString("id")
	// for bot
	id := c.PostForm("bot_id")
	catgroy := c.PostForm("catgroy_id")
	path := c.PostForm("path")

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	files := form.File["files"]
	list := []model.Knowledge{}
	for _, file := range files {
		filename := filepath.Base(file.Filename)
		ext := filepath.Ext(filename)
		if slices.Contains(MARKDOWN, ext) {
			ext = ".md"
		}
		f := filepath.Join("/tmp", filename)
		c.SaveUploadedFile(file, f)
		res, _ := api.Client.Parse(f)
		knowledge := model.Knowledge{
			FolderId:   catgroy,
			BotId:      id,
			OrgId:      orgId,
			UserId:     user,
			Name:       filename,
			Content:    res.Data,
			Ext:        ext,
			Path:       path,
			Status:     2,
			UpdateUser: user,
			Sha:        common.ContentSha(res.Data),
			Tags:       datatypes.JSON("[]"),
		}
		err = knowledge.Insert()
		os.Remove(f)
		if err != nil {
			continue
		}
		events.Emmiter.Emit("knowledge.created", &knowledge)
		list = append(list, knowledge)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    0,
		"data":    list,
	})
}

func SearchKnowledges(c *gin.Context) {
	orgId := c.GetString("orgId")
	botId := c.Query("bot")
	q := c.Query("q")
	knowledges, err := model.SearchKnowledges(q, botId, orgId)
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
		"data":    knowledges,
	})
}

func GetKnowledgesByPath(c *gin.Context) {
	p, _ := strconv.Atoi(c.Query("p"))
	if p < 0 {
		p = 0
	}
	role := c.GetInt("role")
	orgId := c.GetString("orgId")
	path := c.Query("path")
	bot := c.Query("bot")

	var knowledges []*model.Knowledge
	var err error
	// root category
	if path == "0" {
		knowledges, err = model.GetKnowledgesByBotId(bot)
	} else {
		knowledges, err = model.GetKnowledgesByPath(bot, path)
	}
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
		"data": common.Filter(knowledges, func(m *model.Knowledge) bool {
			if role >= common.RoleOperatorUser {
				return true
			}
			return m.OrgId == orgId
		}),
	})
}

func GetKnowledge(c *gin.Context) {
	orgId := c.GetString("orgId")
	id := c.Param("id")
	role := c.GetInt("role")
	knowledge, err := model.GetKnowledgeById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// Pass OperatorUser
	if role < common.RoleOperatorUser {
		// forbidden
		if knowledge.OrgId != orgId {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    304,
				"message": "Forbidden",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    0,
		"data":    knowledge,
	})
}

func CreateKnowledge(c *gin.Context) {

	orgId := c.GetString("orgId")
	user := c.GetString("id")
	knowledge := model.Knowledge{}
	err := c.ShouldBindJSON(&knowledge)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	knowledge.OrgId = orgId
	knowledge.Status = 2
	knowledge.Sha = common.ContentSha(knowledge.Content)
	knowledge.UserId = user
	knowledge.UpdateUser = user
	err = knowledge.Insert()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	events.Emmiter.Emit("knowledge.created", &knowledge)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    0,
		"data":    knowledge,
	})
}

func UpdateKnowledge(c *gin.Context) {

	orgId := c.GetString("orgId")
	user := c.GetString("id")
	role := c.GetInt("role")
	id := c.Param("id")
	old, err := model.GetKnowledgeById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    404,
			"message": err.Error(),
		})
		return
	}

	// Pass OperatorUser
	if role < common.RoleOperatorUser {
		if old.OrgId != orgId {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    403,
				"message": "Forbidden",
			})
			return
		}
	}
	knowledge := model.Knowledge{}
	err = c.ShouldBindJSON(&knowledge)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	sha := common.ContentSha(knowledge.Content)
	if sha == old.Sha {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"code":    0,
			"data":    old,
		})
		return
	}
	// force
	knowledge.Id = id
	knowledge.OrgId = old.OrgId
	knowledge.Sha = sha
	knowledge.UserId = user
	knowledge.Status = 2
	knowledge.UpdateUser = user
	err = knowledge.Update()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	events.Emmiter.Emit("knowledge.updated", &knowledge)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    0,
		"data":    knowledge,
	})
}

func DeleteKnowledge(c *gin.Context) {
	orgId := c.GetString("orgId")
	id := c.Param("id")
	knowledge, err := model.GetKnowledgeById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    404,
			"message": err.Error(),
		})
		return
	}
	if knowledge.OrgId != orgId {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"code":    403,
			"message": "Forbidden",
		})
		return
	}
	knowledge.Delete()
	events.Emmiter.Emit("knowledge.deleted", knowledge)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    0,
		"data":    knowledge,
	})
}

func BatchDeleteKnowledge(c *gin.Context) {
	orgId := c.GetString("orgId")
	ids := Ids{}
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	list := []string{}
	for _, id := range ids {
		knowledge, err := model.GetKnowledgeById(id)
		if err != nil {
			continue
		}
		if knowledge.OrgId != orgId {
			continue
		}
		knowledge.Delete()
		events.Emmiter.Emit("knowledge.deleted", knowledge)
		list = append(list, id)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    0,
		"data":    list,
	})
}
