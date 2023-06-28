package telegram

import (
	"BIMSupportBot/config"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"log"
	"net/http"
)

//type BotConfig struct {
//	cfg     *config.Config
//	bot     *gotgbot.Bot
//	updater *ext.Updater
//}

func InitBot(cfg *config.Config, msgHandler Handlers) {
	// Create bot
	bot, err := gotgbot.NewBot(cfg.TelegramApiConfig.Token, &gotgbot.BotOpts{
		Client: http.Client{},
		DefaultRequestOpts: &gotgbot.RequestOpts{
			Timeout: gotgbot.DefaultTimeout,
			APIURL:  gotgbot.DefaultAPIURL,
		},
	})

	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		Dispatcher: ext.NewDispatcher(&ext.DispatcherOpts{
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				log.Println("an error occurred while handling update:", err.Error())
				return ext.DispatcherActionNoop
			},
			MaxRoutines: ext.DefaultMaxRoutines,
		}),
	})

	// config handlers
	updater.Dispatcher.AddHandler(handlers.NewMessage(message.Text, msgHandler.handleMessage))
	//start webhook
	startWebHook(cfg, bot, updater)
	updater.Idle()
}

func startWebHook(cfg *config.Config, bot *gotgbot.Bot, updater *ext.Updater) {
	webhookOpts := ext.WebhookOpts{
		ListenAddr:  "localhost:8080",             // This example assumes you're in a dev environment running ngrok on 8080.
		SecretToken: cfg.TelegramApiConfig.Secret, // Setting a webhook secret here allows you to ensure the webhook is set by you (must be set here AND in SetWebhook!).
	}

	// We use the token as the urlPath for the webhook, as using a secret ensures that strangers aren't crafting fake updates.
	err := updater.StartWebhook(bot, cfg.TelegramApiConfig.Token, webhookOpts)
	if err != nil {
		panic("failed to start webhook: " + err.Error())
	}

	err = updater.SetAllBotWebhooks(cfg.TelegramApiConfig.Domain, &gotgbot.SetWebhookOpts{
		MaxConnections:     100,
		DropPendingUpdates: true,
		SecretToken:        webhookOpts.SecretToken,
	})
	if err != nil {
		panic("failed to set webhook: " + err.Error())
	}

}
