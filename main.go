package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_collySqliteStorage "olx-crawler/colly/sqlite3"
	"olx-crawler/config"
	_configHTTPDelivery "olx-crawler/config/delivery/http"
	_cron "olx-crawler/cron"
	"olx-crawler/i18n"
	_keywordHTTPDelivery "olx-crawler/keyword/delivery/http"
	_keywordRepository "olx-crawler/keyword/repository"
	_keywordUsecase "olx-crawler/keyword/usecase"
	"olx-crawler/menu"
	_middleware "olx-crawler/middleware"
	"olx-crawler/notifications"
	_observationHTTPDelivery "olx-crawler/observation/delivery/http"
	_observationRepository "olx-crawler/observation/repository"
	_observationUsecase "olx-crawler/observation/usecase"
	_suggestionHTTPDelivery "olx-crawler/suggestion/delivery/http"
	_suggestionRepository "olx-crawler/suggestion/repository"
	_suggestionUsecase "olx-crawler/suggestion/usecase"
	"olx-crawler/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/sirupsen/logrus"

	"github.com/robfig/cron/v3"

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
		logrus.Fatal(err)
	}
	cfg, err := configManager.Config()
	if err != nil {
		logrus.Fatal(err)
	}

	//Logger
	if cfg.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/olx_crawler.log",
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     1, //days
	})

	//I18N
	i18n.LoadMessageFiles("./i18n/locales")

	if cfg.Lang == "" {
		var err error
		cfg.Lang, err = systemLang.DetectLanguage()
		if err != nil {
			logrus.Fatal(err)
		}
		if err := configManager.Save(cfg); err != nil {
			logrus.Fatal(err)
		}
	}
	i18n.SetLanguage(cfg.Lang)

	//Notifications
	notificationsManager, err := notifications.NewManager(configManager)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		if err := notificationsManager.Close(); err != nil {
			logrus.Fatal(err)
		}
	}()

	//DB
	db, err := gorm.Open("sqlite3", "olx_crawler.db")
	defer func() {
		if err := db.Close(); err != nil {
			logrus.Fatal(err)
		}
	}()
	if err != nil {
		logrus.Fatal(err)
	}
	db.LogMode(false)
	db.Exec("PRAGMA foreign_keys = ON;")
	db.Exec("PRAGMA synchronous = OFF;")
	db.Exec("PRAGMA temp_store = MEMORY;")
	db.Exec("PRAGMA journal_mode = WAL;")

	//COLLY CACHE
	collyStorage := &_collySqliteStorage.Storage{
		Filename: "colly_cache.db",
	}
	defer func() {
		if err := collyStorage.Close(); err != nil {
			logrus.Fatal(err)
		}
	}()
	//REPOSITORIES
	keywordRepo, err := _keywordRepository.NewKeywordRepository(db)
	if err != nil {
		logrus.Fatal(err)
	}
	observationRepo, err := _observationRepository.NewObservationRepository(db)
	if err != nil {
		logrus.Fatal(err)
	}
	suggestionRepo, err := _suggestionRepository.NewSuggestionRepository(db)
	if err != nil {
		logrus.Fatal(err)
	}

	//USECASES
	keywordUcase := _keywordUsecase.NewKeywordUsecase(keywordRepo)
	observationUcase := _observationUsecase.NewObservationUsecase(observationRepo)
	suggestionUcase := _suggestionUsecase.NewSuggestionUsecase(suggestionRepo)

	//CRON
	c := cron.New(cron.WithChain(
		cron.SkipIfStillRunning(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags)))))
	err = _cron.AttachHandlers(c, &_cron.Config{
		NotificationsManager: notificationsManager,
		ObservationRepo:      observationRepo,
		SuggestionRepo:       suggestionRepo,
		ConfigManager:        configManager,
		CollyStorage:         collyStorage,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	c.Start()
	defer c.Stop()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	if os.Getenv("DEFAULT_HANDLER") != "true" {
		e.HTTPErrorHandler = customHTTPErrorHandler
	}
	e.Use(middleware.Recover())
	e.Use(_middleware.Logger())
	e.Static("/", "./public")
	g := e.Group("/api")
	_keywordHTTPDelivery.NewKeywordHandler(g, keywordUcase)
	_observationHTTPDelivery.NewObservationHandler(g, observationUcase)
	_suggestionHTTPDelivery.NewSuggestionHandler(g, suggestionUcase)
	_configHTTPDelivery.NewConfigHandler(g, configManager)

	url := fmt.Sprintf(":%d", cfg.Port)
	go func() {
		if err := e.Start(url); err != http.ErrServerClosed {
			logrus.Fatal(err)
		}
	}()
	logrus.Infof("Server is listening on port %s", url)
	serverURL := fmt.Sprintf("http://localhost%s", url)
	if err := utils.OpenBrowser(serverURL); err != nil {
		logrus.Fatal(err)
	}

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	if os.Getenv("DISABLE_MENU") != "true" {
		menu.New(serverURL, channel)
	}
	<-channel

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	e.Shutdown(ctx)
	logrus.Info("shutting down")
}

func customHTTPErrorHandler(err error, c echo.Context) {
	if _, ok := err.(*echo.HTTPError); ok {
		c.File("./public/index.html")
	}
}
