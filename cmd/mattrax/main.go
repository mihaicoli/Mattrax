package main

import (
	"database/sql"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/authentication"
	"github.com/mattrax/Mattrax/internal/certificates"
	"github.com/mattrax/Mattrax/internal/db"
	"github.com/mattrax/Mattrax/internal/middleware"
	"github.com/mattrax/Mattrax/internal/settings"
	"github.com/mattrax/Mattrax/mdm"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var args mattrax.Arguments
	arg.MustParse(&args)
	// TODO: Verify arguments (eg. Domain is domain, cert paths exists, valid listen addr)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if args.Debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	dbconn, err := sql.Open("postgres", args.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initialising Postgres database connection")
	}
	defer dbconn.Close()

	if err := dbconn.Ping(); err != nil {
		log.Fatal().Err(err).Msg("Error communicating with Postgres database")
	}

	q := db.New(dbconn)
	defer q.Close()

	// TODO: Check DB is working by querying

	var srv = &mattrax.Server{
		Args:         args,
		GlobalRouter: mux.NewRouter(),
		DB:           q,
		Cache:        cache.New(5*time.Minute, 10*time.Minute),
	}
	if srv.Settings, err = settings.New(srv.DB); err != nil {
		log.Fatal().Err(err).Msg("Error starting settings service")
	}
	if srv.Cert, err = certificates.New(srv.DB); err != nil {
		log.Fatal().Err(err).Msg("Error starting certificates service")
	}
	if srv.Auth, err = authentication.New(srv.Cert, srv.Cache, srv.DB, args.Domain, args.Debug); err != nil {
		log.Fatal().Err(err).Msg("Error starting authentication service")
	}
	srv.GlobalRouter.Use(middleware.Logging())
	srv.GlobalRouter.Use(middleware.Headers())
	srv.Router = srv.GlobalRouter.Schemes("https").Host(args.Domain).Subrouter()
	mdm.Mount(srv)

	serve(args.Addr, args.Domain, args.TLSCert, args.TLSKey, nil, srv.GlobalRouter)
}
