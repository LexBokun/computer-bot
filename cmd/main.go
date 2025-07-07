package main

import (
	"context"
	"github.com/LexBokun/ControlAgent/internal/delivery/telegram"
	"log/slog"
	"net/http"
	"os"

	"github.com/LexBokun/ControlAgent/config"
	setenable "github.com/LexBokun/ControlAgent/internal/application/service/command/set-enable"
	listdisplays "github.com/LexBokun/ControlAgent/internal/application/service/query/list-displays"
	privatedisplay "github.com/LexBokun/ControlAgent/internal/delivery/private/display"
	multimonitor "github.com/LexBokun/ControlAgent/internal/infrastructure/multi-monitor"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/utrack/clay/v3/server"
	"github.com/utrack/clay/v3/transport"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()

	// logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       config.GetInstance().Logger.Level,
		ReplaceAttr: nil,
	}))
	slog.SetDefault(logger)

	toolConfig := config.GetInstance().MonitorTool
	if _, err := os.Stat(toolConfig.Path); err != nil {
		slog.Error(err.Error())
		return
	}

	// TODO: fix this zalupa
	// repository
	tool := multimonitor.New(toolConfig.Path)
	repo := multimonitor.NewRepository(tool)

	// handlers
	listHandler := listdisplays.NewQueryHandler(repo)
	setHandler := setenable.NewCommandHandler(repo)

	// delivery
	// telegram bot
	bot := telegram.New(config.GetInstance().Telegram, listHandler, setHandler)
	err := bot.Start(ctx)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	privateDisplay := privatedisplay.NewImplementation(listHandler, setHandler)

	hmux := chi.NewMux()

	corsCfg := config.GetInstance().CORS
	hmux.Use(middleware.RealIP)
	hmux.Use(
		cors.Handler(
			cors.Options{
				AllowedOrigins:   corsCfg.AllowedOrigins,
				AllowedMethods:   corsCfg.AllowedMethods,
				AllowedHeaders:   corsCfg.AllowedHeaders,
				ExposedHeaders:   corsCfg.ExposedHeaders,
				AllowCredentials: corsCfg.AllowCredentials,
				MaxAge:           corsCfg.MaxAge,
			},
		),
	)

	hmux.HandleFunc(
		"/docs/*", func(w http.ResponseWriter, r *http.Request) {
			httpSwagger.Handler(httpSwagger.URL("swagger.json")).ServeHTTP(w, r)
		},
	)
	hmux.Get(
		"/docs", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
		},
	)
	hmux.Get(
		"/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/swagger.json", http.StatusMovedPermanently)
		},
	)

	var group errgroup.Group

	group.Go(
		func() error {
			httpServer := config.GetInstance().HTTPServer
			slog.Info("HTTP сервер запускается...", "port", httpServer.Port)
			return server.
				NewServer(
					httpServer.Port,
					server.WithHTTPMux(hmux),
				).Run(

				NewCompoundService(
					privateDisplay,
				),
			)
		})

	err = group.Wait()
	if err != nil {
		slog.Error(err.Error())
	}
}

type CompoundService struct {
	desc *transport.CompoundServiceDesc
}

func NewCompoundService(svcs ...transport.Service) *CompoundService {
	descs := make([]transport.ServiceDesc, 0, len(svcs))

	for _, svc := range svcs {
		descs = append(descs, svc.GetDescription())
	}

	return &CompoundService{
		desc: transport.NewCompoundServiceDesc(descs...),
	}
}

func (c *CompoundService) GetDescription() transport.ServiceDesc {
	return c.desc
}
