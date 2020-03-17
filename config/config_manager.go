package config

import (
	"encoding/json"
	"io/ioutil"
	"olx-crawler/errors"
	"olx-crawler/models"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
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
}

func NewManager() Manager {
	v := viper.New()
	m := &manager{
		v,
		make(map[int]func(in fsnotify.Event)),
		1,
		sync.Mutex{},
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
		return errors.Wrap(errors.ErrCannotReadConfig, []error{err})
	}
	return nil
}

func (m *manager) Config() (*models.Config, error) {
	cfg := &models.Config{}
	if err := m.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(errors.ErrCannotReadConfig, []error{err})
	}
	return cfg, nil
}

func (m *manager) Save(cfg *models.Config) error {
	jsonString, err := json.Marshal(cfg)
	if err != nil {
		return errors.Wrap(errors.ErrCannotSaveConfig, []error{err})
	}
	if err := ioutil.WriteFile("config.json", jsonString, os.ModePerm); err != nil {
		return errors.Wrap(errors.ErrCannotSaveConfig, []error{err})
	}
	return nil
}

func (m *manager) OnConfigChange(run func(in fsnotify.Event)) int {
	m.mutex.Lock()
	id := m.onConfigChangeID
	m.onConfigChange[m.onConfigChangeID] = run
	m.onConfigChangeID++
	m.mutex.Unlock()
	return id
}
