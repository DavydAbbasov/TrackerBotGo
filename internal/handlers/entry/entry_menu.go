package entry

import (
	"context"
	"time"

	"github.com/DavydAbbasov/trecker_bot/interfaces"
	helper "github.com/DavydAbbasov/trecker_bot/internal/lib/postgresql"
	"github.com/DavydAbbasov/trecker_bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/rs/zerolog/log"
)

type EntryModule struct {
	bot        interfaces.BotAPI
	fsm        interfaces.FSMManager
	repository interfaces.UserRepository
	validator  *helper.Validator
}

func New(bot interfaces.BotAPI, fsm interfaces.FSMManager, repository interfaces.UserRepository, validator *helper.Validator) *EntryModule {
	return &EntryModule{
		bot:        bot,
		fsm:        fsm,
		repository: repository,
		validator:  validator,
	}
}
func (e *EntryModule) HandleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case model.CommandStart:
		if err := e.handleStart(message); err != nil {
			return
		}
		e.HandleLanguageStart(message)
	default:
		e.handleUnknown(message)
	}
}

func (e *EntryModule) handleStart(m *tgbotapi.Message) error {
	if m.From == nil { // защита от каналов
		return nil
	}
	// Save user into DB
	user := model.User{
		TgUserID: m.From.ID,
		// указатель оставляем только для InsertUser (он может принять NULL)
		UserName: helper.PtrIfNotEmpty(m.From.UserName),
	}

	log.Info().Msgf("HandleStart triggered for userID=%d chatID=%d", user.TgUserID, m.Chat.ID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.repository.InsertUser(ctx, user.TgUserID, user.UserName); err != nil {
		log.Error().Err(err).Int64("tg_user_id", user.TgUserID).Msg("ensure user failed")
		return err
	}

	if _, err := e.repository.UpdateUsername(ctx, user.TgUserID, m.From.UserName); err != nil {
		log.Warn().Err(err).Msg("update username skipped")

	}

	if code, err := e.validator.ValidateLanguage(m.From.LanguageCode); err == nil {
		if _, err := e.repository.UpdateLanguage(ctx, m.From.ID, code); err != nil {
			log.Warn().Err(err).Str("code", code).Msg("update language failed")

		}
	}

	if _, err := e.bot.Send(tgbotapi.NewMessage(m.Chat.ID, model.WELCOME)); err != nil {
		log.Error().Err(err).Msg("failed to send welcome")
		return err
	}
	return nil
}

func (e *EntryModule) handleUnknown(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, model.TextUnknownCmd)
	_, err := e.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("unknown command error")
	}
}
