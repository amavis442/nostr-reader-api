package http

import (
	"amavis442/nostr-reader/internal/db"
	domain "amavis442/nostr-reader/internal/domain"
	wrapper "amavis442/nostr-reader/internal/nostr"
	"fmt"
	"log/slog"
	"mime"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

type ServerConfig struct {
	Port int64 `json:"port"`
	// AllowedIPs lists client IPs (besides loopback) that may reach the API.
	// Empty means loopback-only.
	AllowedIPs []string `json:"allowed_ips"`
	// AllowedOrigins lists the CORS origins permitted to call the API. Empty
	// falls back to localhost origins.
	AllowedOrigins []string `json:"allowed_origins"`
}

type HttpServer struct {
	DevMode    bool
	Server     *ServerConfig
	Database   *db.Storage
	Nostr      *wrapper.Wrapper
	Router     *chi.Mux
	MaxNoteAge *domain.MaxNoteAge
}

func (s *HttpServer) Routes(c *Controller, port string) {
	s.Router = routes(c, port, s.Server)
}

// @title Nostr Reader API
// @version 1.0
// @description Nostr Reader Api.
// @termsOfService http://swaggerouter.io/terms/

// @contact.name API Support
// @contact.url http://www.swaggerouter.io/support
// @contact.email support@swaggerouter.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func (s *HttpServer) Start() {
	var c Controller
	c.Pubkey = s.Nostr.Cfg.PubKey
	c.Db = s.Database
	c.Nostr = s.Nostr
	c.MaxNoteAge = s.MaxNoteAge

	var port string = "8080"
	if s.Server.Port > 0 {
		port = fmt.Sprint(s.Server.Port)
	}

	s.Routes(&c, port)
	//router := routes(&c, port)

	// Windows may be missing this
	mime.AddExtensionType(".js", "application/javascript")

	slog.Info(fmt.Sprint("Server running: http://localhost:" + port))

	err := http.ListenAndServe(":"+port, s.Router)
	if err != nil {
		slog.Info(fmt.Sprintf("Could not start http server on this port: %s", port))
		slog.Error(err.Error())
		os.Exit(1)
	}
}
