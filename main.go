package main

import (
	"context"
	"fmt"
	"log"
	_collySqliteStorage "olx-crawler/colly/sqlite3"
	"olx-crawler/config"
	_configHTTPDelivery "olx-crawler/config/delivery/http"
	_cron "olx-crawler/cron"
	"olx-crawler/i18n"
	_middleware "olx-crawler/middleware"
	"olx-crawler/notifications"
	_observationHTTPDelivery "olx-crawler/observation/delivery/http"
	_observationRepository "olx-crawler/observation/repository"
	_observationUsecase "olx-crawler/observation/usecase"
	_suggestionHTTPDelivery "olx-crawler/suggestion/delivery/http"
	_suggestionRepository "olx-crawler/suggestion/repository"
	_suggestionUsecase "olx-crawler/suggestion/usecase"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
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

	//Logger
	if configManager.GetBool("debug") {
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

	//I18N
	i18n.LoadMessageFiles("i18n/locales")

	lang := configManager.GetString("lang")
	if lang == "" {
		var err error
		lang, err = systemLang.DetectLanguage()
		if err != nil {
			logrus.Fatal(err)
		}
	}
	i18n.SetLanguage(lang)

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
	observationRepo, err := _observationRepository.NewObservationRepository(db)
	if err != nil {
		logrus.Fatal(err)
	}
	suggestionRepo, err := _suggestionRepository.NewSuggestionRepository(db)
	if err != nil {
		logrus.Fatal(err)
	}

	// s := true
	// observationRepo.Store(&models.Observation{
	// 	Name: "Dysk SSD",
	// 	URL:  "https://www.olx.pl/elektronika/q-Dysk-SSD/?search%5Bfilter_float_price%3Afrom%5D=100&search%5Bfilter_float_price%3Ato%5D=500",
	// 	OneOf: []models.OneOf{
	// 		models.OneOf{
	// 			For:   "title",
	// 			Value: "128 GB",
	// 		},
	// 		models.OneOf{
	// 			For:   "title",
	// 			Value: "128 gb",
	// 		},
	// 		models.OneOf{
	// 			For:   "title",
	// 			Value: "128GB",
	// 		},
	// 		models.OneOf{
	// 			For:   "title",
	// 			Value: "128gb",
	// 		},
	// 		models.OneOf{
	// 			For:   "description",
	// 			Value: "m2",
	// 		},
	// 	},
	// 	Excluded: []models.Excluded{
	// 		models.Excluded{
	// 			For:   "title",
	// 			Value: "komputer",
	// 		},
	// 		models.Excluded{
	// 			For:   "title",
	// 			Value: "laptop",
	// 		},
	// 	},
	// 	Started: &s,
	// })

	//USECASES
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
	e.Use(middleware.Recover())
	e.Use(_middleware.Logger())
	e.Static("/", "/public")
	g := e.Group("/api")
	_observationHTTPDelivery.NewObservationHandler(g, observationUcase)
	_suggestionHTTPDelivery.NewSuggestionHandler(g, suggestionUcase)
	_configHTTPDelivery.NewConfigHandler(g, configManager)

	url := fmt.Sprintf(":%d", configManager.GetInt("port"))
	go func() {
		e.Start(url)
	}()
	logrus.Infof("Server is listening on port %s", e.Server.Addr)
	if err := openbrowser(fmt.Sprintf("http://localhost%s", url)); err != nil {
		logrus.Fatal(err)
	}

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	<-channel

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	e.Shutdown(ctx)
	logrus.Info("shutting down")
}

func openbrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("cmd", "/C", "start", url).Run()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}
