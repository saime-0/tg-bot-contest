package tg

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/jmoiron/sqlx"

	tgClient "github.com/Saime-0/tg-bot-contest/internal/tg/client"
	tgUpdatesController "github.com/Saime-0/tg-bot-contest/internal/tg/updatesController"
)

func Run(ctx context.Context, token string, db *sqlx.DB) (err error) {
	client := &tgClient.Client{}
	if client.Bot, err = tgClient.NewBot(token); err != nil {
		return fmt.Errorf("tgClient.NewBot: %w", err)
	}

	updatesController := &tgUpdatesController.Controller{
		DB:  db,
		Bot: client.Bot,
	}

	client.Dispatcher = tgClient.NewDispatcher()
	if err = updatesController.AddHandlers(client.Dispatcher); err != nil {
		return fmt.Errorf("updatesController.AddHandlers: %w", err)
	}

	client.Updater = ext.NewUpdater(client.Dispatcher, nil)

	if err = tgClient.StartPolling(client.Updater, client.Bot); err != nil {
		return fmt.Errorf("tgClient.StartPolling: %w", err)
	}

	updaterCtx, cancel := context.WithCancel(context.Background())
	go func() {
		client.Updater.Idle()
		cancel()
	}()

	select {
	case <-ctx.Done():
		cancel()
		if err = client.Updater.Stop(); err != nil {
			slog.Warn("client.Updater.Stop: " + err.Error())
		}
		return nil
	case <-updaterCtx.Done():
		return nil
	}
}
