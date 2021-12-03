package main

import (
	routerPkg "github.com/denlipov/omp-bot/internal/app/router"
	"github.com/denlipov/omp-bot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatal().Err(err).Msg("Failed init configuration")
	}
	cfg := config.GetConfig()

	if cfg.Project.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		log.Fatal().Err(err)
	}

	// Uncomment if you want debugging
	// bot.Debug = true

	log.Info().Msgf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.UpdateConfig{
		Timeout: cfg.Bot.UpdateTimeout,
	}

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal().Err(err)
	}

	routerHandler := routerPkg.NewRouter(bot)

	for update := range updates {
		routerHandler.HandleUpdate(update)
	}
}
