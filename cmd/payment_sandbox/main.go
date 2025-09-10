package main

import (
	"context"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/DavydAbbasov/trecker_bot/payment/events"
	"github.com/caarlos0/env/v10"
	"github.com/rs/zerolog"
)

type Config struct {
	Port string `env:"PAY_SANDBOX_PORT" envDefault:":8080"`
	//время ожидания перед завершением сервера (для graceful shutdown).
	ShutdownTimeout time.Duration `env:"PAY_SHUTDOWN_TIMEOUT" envDefault:"5s"`
}

func main() {

	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		logger.Fatal().Err(err).Msg("load config")
	}
	logger.Info().
		Str("port", cfg.Port).
		Dur("shutdown_timeout", cfg.ShutdownTimeout).
		Msg("sandbox config")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	bus := events.NewBus(1024)

	srv := startHTTP(logger, cfg, bus)

	go func() {
		logger.Info().Str("addr", cfg.Port).Msg("http listen")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("http listen")
		}
	}()

	logger.Info().Msg("sandbox started")
	<-ctx.Done()

	//Роль: тайм-лимит на мягкое завершение HTTP-сервера
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	//закрываем вход событий и завершаем воркеров (когда их добавишь)
	bus.Close()

	// srv.Shutdown(shutdownCtx)- мягко гасим
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("http shutdown")
	}
	defer logger.Info().Msg("sandbox stopped")

}

// функцию startHTTP — она создаёт и запускает сервер.
func startHTTP(logger zerolog.Logger, cfg Config, bus events.Bus) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	mux.HandleFunc("/payment/webhook", func(w http.ResponseWriter, r *http.Request) {
		// 1) Method

		if r.Method != http.MethodPost {
			http.Error(w, "metod not allowed", http.StatusMethodNotAllowed) //405
			return
		}
		// 2) type

		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || (mediaType != "application/json" && !strings.HasSuffix(mediaType, "+json")) {
			http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
			return
		}

		// 3) limit + reading
		defer r.Body.Close()
		lr := http.MaxBytesReader(w, r.Body, 1<<20)
		body, err := io.ReadAll(lr)
		if err != nil {
			http.Error(w, "bad reguest", http.StatusBadRequest)
			return
		}

		var m struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		}

		if err := json.Unmarshal(body, &m); err != nil || m.Type == "" || m.ID == "" {
			http.Error(w, "bad request", http.StatusBadRequest) //400
			return
		}

		ctxPub, cancel := context.WithTimeout(r.Context(), 100*time.Millisecond)
		defer cancel()

		evt := events.Event{
			Type:       m.Type,
			ID:         m.ID,
			Raw:        body,
			ReceivedAt: time.Now(),
		}

		if err := bus.Publish(ctxPub, evt); err != nil {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "queue busy", http.StatusServiceUnavailable) // 503
			return
		}

		logger.Info().
			Str("route", "/payment/webhook").
			Str("event_type", m.Type).
			Str("event_id", m.ID).
			Int("size", len(body)).
			Msg("webhook received")

		// 7) ответить провайдеру
		w.WriteHeader(http.StatusOK)

	})
	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: mux,
	}
	return srv

}
