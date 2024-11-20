package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"javacode-test/internal/app/processor"
	"javacode-test/internal/config/dbconfig"
	"javacode-test/internal/config/getconfig"
	"javacode-test/internal/config/httpconfig"
	"javacode-test/internal/config/postconfig"
	"javacode-test/internal/http-server/handlers/apply"
	"javacode-test/internal/http-server/handlers/get"
	"javacode-test/internal/storage/psql"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

const (
	postURL       = "/api/v1/wallet"
	walletIdParam = "walletId"
)

var (
	errInvalidLogger  = errors.New("failed to initialize zap logger")
	errInvalidStorage = errors.New("failed to initialize storage")
	getURL            = fmt.Sprintf("/api/v1/wallets/{%s}", walletIdParam)
)

func main() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(errInvalidLogger)
	}
	defer log.Sync()

	dbCfg, err := dbconfig.New()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		dbCfg.DBUser,
		dbCfg.DBPass,
		dbCfg.DBHost,
		dbCfg.DBPort,
		dbCfg.DBName,
	)
	storage, err := psql.New(dbURL)
	if err != nil {
		log.Fatal(errInvalidStorage.Error(), zap.Error(err))
		os.Exit(2)
	}

	proc := processor.New(storage)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	getCfg, err := getconfig.New()
	router.Get(getURL, get.New(getCfg, log, proc, walletIdParam))

	postCfg, err := postconfig.New()
	router.Post(postURL, apply.New(postCfg, log, proc))

	httpCfg, err := httpconfig.New()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	log.Info(fmt.Sprintf("starting server at %s:%d", httpCfg.HTTPHost, httpCfg.HTTPPort))
	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%d", httpCfg.HTTPHost, httpCfg.HTTPPort),
		Handler:      router,
		ReadTimeout:  httpCfg.HTTPReadTimeout,
		WriteTimeout: httpCfg.HTTPWriteTimeout,
		IdleTimeout:  httpCfg.HTTPIdleTimeout,
	}
	log.Info(fmt.Sprintf("running server at %s:%d", httpCfg.HTTPHost, httpCfg.HTTPPort))
	if err := srv.ListenAndServe(); err != nil {
		log.Error(err.Error())
	}
}
