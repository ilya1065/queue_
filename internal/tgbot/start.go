package tgbot

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"queue/internal/infra"
	"queue/internal/server"
	"time"

	"golang.org/x/net/proxy"
	"gopkg.in/telebot.v4"
)

func newTelegramHTTPClient() (*http.Client, error) {
	proxyAddr := os.Getenv("TG_PROXY_ADDR") // 45.80.228.147:1080
	proxyUser := os.Getenv("TG_PROXY_USER") // proxyuser
	proxyPass := os.Getenv("TG_PROXY_PASS") // lusa

	if proxyAddr == "" {
		return &http.Client{Timeout: 30 * time.Second}, nil
	}

	dialer, err := proxy.SOCKS5("tcp", proxyAddr, &proxy.Auth{
		User:     proxyUser,
		Password: proxyPass,
	}, proxy.Direct)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		},
		ForceAttemptHTTP2:     true,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		IdleConnTimeout:       90 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   60 * time.Second,
	}, nil
}

// StartBot запуск Telegram-бота и регистрация всех роутов
func StartBot(srv *server.Server, inf *infra.Infra) {
	slog.Info("Старт бота")
	if os.Getenv("TG_KEY") == "" {
		log.Fatal("TG_KEY is not set")
	}

	httpClient, err := newTelegramHTTPClient()
	if err != nil {
		log.Fatal(err)
	}

	pref := telebot.Settings{
		Token:  os.Getenv("TG_KEY"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		Client: httpClient,
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	ctl := NewController(srv, inf)
	ctl.RegisterRoutes(b)
	slog.Info("Бот готов, запускаю polling")
	b.Start()
}
