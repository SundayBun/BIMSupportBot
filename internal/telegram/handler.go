package telegram

import (
	"BIMSupportBot/config"
	"BIMSupportBot/repository"
	"context"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Handlers interface {
	handleMessage(b *gotgbot.Bot, ctx *ext.Context) error
}

type messageHandler struct {
	cfg     *config.Config
	repo    repository.Repository
	context context.Context
}

func NewMessageHandler(cfg *config.Config, repository repository.Repository, context context.Context) Handlers {
	return &messageHandler{cfg: cfg, repo: repository, context: context}
}

func (msgh messageHandler) handleMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	//todo find response in db or return unknown error
	//todo if return unknown error-> add to db "uknown requests"
	income := ctx.EffectiveMessage.Text

	_, err := msgh.repo.GetById(msgh.context, income)

	if err != nil {
		return fmt.Errorf("failed to get answer: %w", err)
	}

	_, err = ctx.EffectiveMessage.Reply(b, ctx.EffectiveMessage.Text, nil)
	if err != nil {
		return fmt.Errorf("failed to echo message: %w", err)
	}
	return nil
}
