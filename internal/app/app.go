package app

import (
	"context"
	"eMobile/docs"
	"eMobile/internal/config"
	"eMobile/internal/crud"
	"eMobile/internal/repo"
	"eMobile/internal/route"
	"eMobile/internal/service"
	"eMobile/pkg/logging"
	"eMobile/pkg/migrator"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	// init logger and config
	log := logging.GetLogger()
	conf := config.GetConfig(log)

	// init connection pool and repository
	db := crud.GetPool(conf, log)
	repositories := repo.NewRepository(db, log)

	// run migrations
	m, err := migrator.NewMigrator(migrator.Deps{
		Username: conf.Storage.Username,
		Password: conf.Storage.Password,
		Host:     conf.Storage.Host,
		Port:     conf.Storage.Port,
		Database: conf.Storage.Database,
		Source:   conf.Storage.Migration,
	})

	if err != nil {
		log.Fatal("Error initializing migrator: ", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("Error on run migration: ", err)
	}

	// init services
	services := service.NewService(&service.Deps{
		Repo:    repositories,
		Logger:  log,
		InfoURL: conf.Server.InfoServiceUrl,
	})

	// init router
	router := httprouter.New()

	if conf.Server.EnableSwag {
		router.Handler("GET", "/swagger/*any", httpSwagger.Handler(
			httpSwagger.URL("http://"+conf.Server.ExternalHost+":"+conf.Server.ExternalPort+"/swagger/doc.json"),
		))
		docs.SwaggerInfo.Host = conf.Server.ExternalHost + ":" + conf.Server.ExternalPort
		log.Info("Swagger enabled")
	}

	// init handler
	handler := route.NewHandler(route.Deps{
		Service: services,
		Logger:  log,
		Config:  conf,
	})
	handler.Init(router)

	// init and run server
	appMux := NewAppMux(&Deps{
		Router: router,
		Log:    log,
		Conf:   conf,
	})

	url := fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.Port)
	server := http.Server{
		Addr:    url,
		Handler: appMux,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		sig := <-sigChan
		log.Info("Shutdown signal:", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	log.Error(server.ListenAndServe())
}

// TODO add comments
// TODO validate after <= before
// TODO check if resp JSON is empty in getAudioInfo
