package notification

import (
	"context"
	"fmt"
	"strconv"

	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
	"github.com/carlosgab83/matrix/go/internal/tank/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramNotifier struct {
	Config domain.Config
	Logger logging.Logger
	Bot    *tgbotapi.BotAPI
}

func NewTelegramNotifier(cfg domain.Config, logger logging.Logger) (*TelegramNotifier, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotAPIToken)
	if err != nil {
		return nil, fmt.Errorf("error creating telegram bot: %v", err)
	}

	bot.Debug = false // Enable debug mode for more detailed logs
	logger.Info("Telegram bot authorized", "account", bot.Self.UserName)

	return &TelegramNotifier{
		Config: cfg,
		Logger: logger,
		Bot:    bot,
	}, nil
}

func (tn *TelegramNotifier) Notify(ctx context.Context, chatID string, payload string) error {
	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return fmt.Errorf("error converting CHAT_ID: %v", err)
	}

	msg := tgbotapi.NewMessage(chatIDInt, payload)
	_, err = tn.Bot.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending message \"%s\" to %d: %v", payload, chatIDInt, err)
	}

	tn.Logger.Info("Message sent to", "chatID", chatID)
	return nil
}

func (tn *TelegramNotifier) Register(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := tn.Bot.GetUpdatesChan(u)
	tn.Logger.Info("Bot listenning messages...")
	for {
		select {
		case update, ok := <-updates:
			if !ok {
				return
			}

			if update.Message == nil {
				continue
			}

			chatID := update.Message.Chat.ID
			userName := update.Message.From.UserName
			tn.Logger.Info("Received message", "user", userName, "chatID", chatID, "text", update.Message.Text)
			msg := tgbotapi.NewMessage(chatID, "Welcome!.")
			tn.Bot.Send(msg)
		case <-ctx.Done():
			return
		}
	}
}

func (tn *TelegramNotifier) Close() error {
	//TODO
	return nil
}
