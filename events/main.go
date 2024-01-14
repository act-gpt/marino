package events

import (
	"fmt"

	"github.com/act-gpt/marino/model"

	"github.com/kataras/go-events"
	"github.com/zeromicro/go-zero/core/logx"
)

var Emmiter events.EventEmmiter

func init() {
	Emmiter = events.New()

	Emmiter.On("user.created", func(payload ...interface{}) {
		user := payload[0].(model.User)
		go UserCreated(user)
		logx.Info(fmt.Sprintf("Register user %s, email %s", user.Username, user.Email))
	})

	Emmiter.On("user.login", func(payload ...interface{}) {
		user := payload[0].(model.User)
		go UserLogin(user)
		logx.Info(fmt.Sprintf("Login user %s, email %s", user.Username, user.Email))
	})

	Emmiter.On("org.created", func(payload ...interface{}) {
		/*
			org := payload[0].(controller.OrganizationPost)
			size := payload[1].(int)
			roles := payload[2].(int)
			logx.Info(fmt.Sprintf("Register org %s, size: %d, roles: %d", org.Name, size, roles))
		*/
	})

	Emmiter.On("bot.created", func(payload ...interface{}) {
		bot := payload[0].(model.Bot)
		go CreateBot(bot)
		logx.Info(fmt.Sprintf("Create bot %s ", bot.Id))
	})

	Emmiter.On("bot.deleted", func(payload ...interface{}) {
		bot := payload[0].(*model.Bot)
		go DeleteBot(bot)
		logx.Info(fmt.Sprintf("Delete bot %s ", bot.Id))
	})

	Emmiter.On("folder.created", func(payload ...interface{}) {
		folder := payload[0].(model.Folder)
		logx.Info(fmt.Sprintf("Create knowledge %s ", folder.Id))

	})
	Emmiter.On("folder.updated", func(payload ...interface{}) {
		folder := payload[0].(model.Folder)
		logx.Info(fmt.Sprintf("Update knowledge %s ", folder.Id))
	})

	Emmiter.On("folder.deleted", func(payload ...interface{}) {
		folder := payload[0].(*model.Folder)
		go DeleteFolder(folder)
		logx.Info(fmt.Sprintf("Delete knowledge %s ", folder.Id))
	})

	Emmiter.On("knowledge.created", func(payload ...interface{}) {
		knowledge := payload[0].(*model.Knowledge)
		go CreateKnowledge(knowledge)
		logx.Info(fmt.Sprintf("Create knowledge %s ", knowledge.Id))
	})

	Emmiter.On("knowledge.updated", func(payload ...interface{}) {
		knowledge := payload[0].(*model.Knowledge)
		go UpdateKnowledge(knowledge)
		logx.Info(fmt.Sprintf("Update knowledge %s ", knowledge.Id))
	})

	Emmiter.On("knowledge.deleted", func(payload ...interface{}) {
		knowledge := payload[0].(*model.Knowledge)
		go DeleteKnowledge(knowledge)
		logx.Info(fmt.Sprintf("Delete knowledge %s ", knowledge.Id))
	})

	Emmiter.On("query", func(payload ...interface{}) {
		msg := payload[0].(model.Message)
		go Query(msg)
		str := fmt.Sprintf("Convesation finished %s by %s, cost %.1f's, lft %.1f's, llm %.1f's",
			msg.Model,
			msg.Status,
			msg.CostTime,
			msg.LLMFirstTime,
			msg.LLMTime)
		logx.Info(str)
	})

	Emmiter.On("overquota", func(payload ...interface{}) {
		id := payload[0].(string)
		// TODO: send email notification
		logx.Info(fmt.Sprintf("Bot over quota %s ", id))
	})

	Emmiter.On("green", func(payload ...interface{}) {
		id := payload[0].(string)
		code := payload[1].(int)
		logx.Info(fmt.Sprintf("Bot not pass green %s, %d", id, code))
	})
}
