package main

import (
	"context"
	"fmt"
	"log"
	"olx-crawler/config"
	_configHTTPDelivery "olx-crawler/config/delivery/http"
	"olx-crawler/i18n"
	"olx-crawler/notifications"
	_observationHTTPDelivery "olx-crawler/observation/delivery/http"
	_observationRepository "olx-crawler/observation/repository"
	_observationUsecase "olx-crawler/observation/usecase"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"

	systemLang "github.com/cloudfoundry-attic/jibber_jabber"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	//Config
	configManager := config.NewManager()
	if err := configManager.Init(); err != nil {
		log.Fatal(err)
	}

	notificationsManager, err := notifications.NewNotificationsManager(configManager)
	if err != nil {
		log.Fatal(err)
	}
	defer notificationsManager.Close()

	//I18N
	i18n.LoadMessageFiles("i18n/locales")

	lang := configManager.GetString("lang")
	if lang == "" {
		var err error
		lang, err = systemLang.DetectLanguage()
		if err != nil {
			log.Fatal(err)
		}
	}
	i18n.SetLanguage(lang)

	//DB
	db, err := gorm.Open("sqlite3", "olx_crawler.db")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	//REPOSITORIES
	observationRepo, err := _observationRepository.NewObservationRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	//USECASES
	observationUcase := _observationUsecase.NewObservationUsecase(observationRepo)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	g := e.Group("/api")
	_observationHTTPDelivery.NewObservationHandler(g, observationUcase)
	_configHTTPDelivery.NewConfigHandler(g, configManager)

	port := fmt.Sprintf(":%d", configManager.GetInt("port"))
	go func() {
		e.Start(port)
	}()
	log.Printf("Server is listening on port %s", port)

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	<-channel

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	e.Shutdown(ctx)
	log.Print("shutting down")
}
