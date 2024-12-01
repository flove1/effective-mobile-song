// Song API.
// API for song library.
//
//	Schemes: http, https
//	BasePath: /
//	Version: 1.0.0
//	Host:
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package main

import (
	"effective-mobile/go/config"
	"effective-mobile/go/internal/api/http"
	"effective-mobile/go/internal/song"
	"effective-mobile/go/pkg/database"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//go:generate swagger generate spec -o ../swagger.json

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Error("failed to parse config: ", err)
		os.Exit(1)
	}

	if cfg.Mode == "development" {
		log.SetLevel(log.DebugLevel)
	}

	runMigrations(cfg)

	db, err := database.NewPostgresConnection(cfg.DB.ToDSN())
	if err != nil {
		log.Error("failed to connect to database: ", err)
		os.Exit(1)
	}

	defer db.Close()

	songRepo := song.NewSongRepository(cfg, db)
	songService := song.NewSongService(cfg, songRepo)
	songHandler := song.NewSongHandler(cfg, songService)

	server := http.NewServer(cfg, http.Handlers{
		SongHandler: songHandler,
	})
	server.Start()

	log.Info("server started on port ", cfg.HttpPort)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case s := <-interrupt:
		log.Debug("interrupt signal received: ", s.String())
	case err = <-server.Notify():
		log.Error("server notify: ", err.Error())
	}

	err = server.Shutdown()
	if err != nil {
		log.Error("server shutdown err: ", err)
	}

	log.Info("server exiting")
}

func runMigrations(cfg *config.Config) {
	m, err := migrate.New(
		`file://migrations`,
		cfg.DB.ToDSN())
	if err != nil {
		log.Error("failed to connect to database: ", err)
		os.Exit(1)
	}

	defer m.Close()
	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Error("failed to execute migrations: ", err)
		os.Exit(1)
	}

	log.Info("migrations executed successfully")
}
