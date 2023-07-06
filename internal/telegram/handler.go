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
	handleCallbackQuery(b *gotgbot.Bot, ctx *ext.Context) error
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
	income := ctx.EffectiveMessage.Text
	entityList := msgh.repo.FullTextSearch(msgh.context, income)
	//startInlineQuery(b, ctx, entityList)
	startQuery(b, ctx, entityList)
	return nil
}

func (msgh messageHandler) handleCallbackQuery(b *gotgbot.Bot, ctx *ext.Context) error {
	cb := ctx.Update.CallbackQuery
	var entity, err = msgh.repo.GetById(msgh.context, cb.Data)

	if err != nil {
		return fmt.Errorf("failed to get entity by id: %w", err)
	}
	_, err = cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: entity.Description,
	})
	if err != nil {
		return fmt.Errorf("failed to answer to callback query: %w", err)
	}
	return nil
}

func startInlineQuery(b *gotgbot.Bot, ctx *ext.Context, entityList []repository.BimEntity) error {
	_, err := ctx.InlineQuery.Answer(b, createInlineQueryList(entityList), &gotgbot.AnswerInlineQueryOpts{
		IsPersonal: true,
	})

	//_, err := ctx.EffectiveMessage.Reply(b, ctx.EffectiveMessage.Text, nil)
	if err != nil {
		return fmt.Errorf("failed to answer message: %w", err)
	}
	return nil
}

func createInlineQueryList(entityList []repository.BimEntity) []gotgbot.InlineQueryResult {
	var queryArray []gotgbot.InlineQueryResult
	for _, i := range entityList {
		var inlineQR = gotgbot.InlineQueryResultArticle{
			Id: i.ID.String(),
			InputMessageContent: gotgbot.InputTextMessageContent{
				MessageText: i.Title,
			},
			Description: "Link to the source description",
		}
		queryArray = append(queryArray, inlineQR)
	}
	return queryArray
}

func startQuery(b *gotgbot.Bot, ctx *ext.Context, entityList []repository.BimEntity) error {
	_, err := ctx.EffectiveMessage.Reply(b, "Link to the source description", &gotgbot.SendMessageOpts{
		ParseMode: "html",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: createInlineKeyboard(entityList)},
	},
	)
	if err != nil {
		return fmt.Errorf("failed to send query message: %w", err)
	}
	return nil
}

func createInlineKeyboard(entityList []repository.BimEntity) [][]gotgbot.InlineKeyboardButton {
	var resultList [][]gotgbot.InlineKeyboardButton
	var inlineKeyboardButton []gotgbot.InlineKeyboardButton

	for _, i := range entityList {
		var inlineKB = gotgbot.InlineKeyboardButton{
			Text:         i.Title,
			CallbackData: i.ID.String(),
		}
		inlineKeyboardButton = append(inlineKeyboardButton, inlineKB)
	}
	resultList = append(resultList, inlineKeyboardButton)
	return resultList
}
