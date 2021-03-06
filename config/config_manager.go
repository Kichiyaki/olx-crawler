package config

import (
	"encoding/json"
	"io/ioutil"
	"olx-crawler/errors"
	"olx-crawler/models"
	"olx-crawler/utils"
	"os"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

var (
	Version = "1.0.0"
)

type Manager interface {
	Init() error
	Config() (*models.Config, error)
	Save(*models.Config) error
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	OnConfigChange(run func(in fsnotify.Event)) int
}

type manager struct {
	*viper.Viper
	onConfigChange   map[int]func(in fsnotify.Event)
	onConfigChangeID int
	mutex            sync.Mutex
	logrus           *logrus.Entry
}

func NewManager() Manager {
	v := viper.New()
	m := &manager{
		v,
		make(map[int]func(in fsnotify.Event)),
		1,
		sync.Mutex{},
		logrus.WithField("package", "config"),
	}
	v.OnConfigChange(func(in fsnotify.Event) {
		for _, run := range m.onConfigChange {
			run(in)
		}
	})
	return m
}

func (m *manager) Init() error {
	m.SetConfigName("config.json")
	m.AddConfigPath(".")
	m.SetConfigType("json")
	m.WatchConfig()
	if err := m.ReadInConfig(); err != nil {
		m.logrus.Debugf("Cannot load config.json file: %s", err.Error())
		return errors.Wrap(errors.ErrCannotReadConfig, []error{err})
	}
	return nil
}

func (m *manager) Config() (*models.Config, error) {
	cfg := &models.Config{}
	if err := m.Unmarshal(cfg); err != nil {
		m.logrus.Debugf("Cannot unmarshal config into struct: %s", err.Error())
		return nil, errors.Wrap(errors.ErrCannotReadConfig, []error{err})
	}
	cfg.Version = Version
	return cfg, nil
}

func (m *manager) Save(cfg *models.Config) error {
	bytes, err := json.Marshal(cfg)
	if err != nil {
		m.logrus.Debug(err.Error())
		return errors.Wrap(errors.ErrCannotSaveConfig, []error{err})
	}
	if err := ioutil.WriteFile("config.json", bytes, os.ModePerm); err != nil {
		m.logrus.WithField("config", string(bytes)).Debug(err.Error())
		return errors.Wrap(errors.ErrCannotSaveConfig, []error{err})
	}
	return nil
}

func (m *manager) OnConfigChange(run func(in fsnotify.Event)) int {
	if cfg, err := m.Config(); err == nil {
		m.logrus.WithField("config", string(utils.MustMarshal(cfg))).Debug("Config change")
	}
	m.mutex.Lock()
	id := m.onConfigChangeID
	m.onConfigChange[m.onConfigChangeID] = run
	m.onConfigChangeID++
	m.mutex.Unlock()
	return id
}
