package cron

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"olx-crawler/models"
	_notificationsMocks "olx-crawler/notifications/mocks"
	"olx-crawler/observation/mocks"
	_suggestionMocks "olx-crawler/suggestion/mocks"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gocolly/colly/v2"
)

func TestFetchSuggestions(t *testing.T) {
	t.Run("Should correctly navigate to next page", func(t *testing.T) {
		mux := http.NewServeMux()
		ts := httptest.NewServer(mux)
		defer ts.Close()
		now := time.Now()
		called := false

		mux.HandleFunc("/pagination", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(fmt.Sprintf(`
			<div class="wrapper"></div>
			<div class="wrapper"></div>
			<div class="wrapper">
			<div id="offers_table">
				<div class="wrap">
				<p class="price">
					<strong>Cena</strong>
				</p>
				<a href="other_page">
					<strong>Title</strong>
				</a>
				<table data-id="123123">
					<div>
					<i data-icon="clock"></i>
					dzisiaj %s
					</div>
				</table>
				</div>
			</div>
			<a href="%s/page_2" data-cy="page-link-next">next</a>
			</div>
				`, now.Format("15:04"), ts.URL)))
		})

		mux.HandleFunc("/page_2", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			called = true
		})

		o := &models.Observation{
			Name: "Random name",
			URL:  ts.URL + "/pagination",
			Checked: []models.Checked{
				models.Checked{
					Value: "other_page",
				},
			},
			LastCheckAt: time.Now(),
		}

		observationMockRepo := &mocks.Repository{
			FetchFunc: func(*models.ObservationFilter) (models.PaginatedResponse, error) {
				return models.PaginatedResponse{
					Items: []*models.Observation{
						o,
					},
					Total: 1,
				}, nil
			},
			UpdateFunc: func(*models.Observation) error {
				return nil
			},
		}

		h := &handler{
			observationRepo: observationMockRepo,
			collector:       newCollector(),
			logrus:          logrus.WithField("package", "cron_test"),
		}

		h.fetchSuggestions()
		if !called {
			t.Error("didnt change page")
		}
	})

	t.Run("should not visit ad page", func(t *testing.T) {
		mux := http.NewServeMux()
		ts := httptest.NewServer(mux)
		defer ts.Close()
		now := time.Now()
		called := false
		o := &models.Observation{
			Name:        "Random name",
			URL:         ts.URL + "/notstore",
			LastCheckAt: now.AddDate(0, 0, 1),
		}
		observationMockRepo := &mocks.Repository{
			FetchFunc: func(*models.ObservationFilter) (models.PaginatedResponse, error) {
				return models.PaginatedResponse{
					Items: []*models.Observation{
						o,
					},
					Total: 1,
				}, nil
			},
			UpdateFunc: func(*models.Observation) error {
				return nil
			},
		}
		h := &handler{
			observationRepo: observationMockRepo,
			collector:       newCollector(),
			logrus:          logrus.WithField("package", "cron_test"),
		}
		mux.HandleFunc("/notstore", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(fmt.Sprintf(`
			<div class="wrapper"></div>
			<div class="wrapper"></div>
			<div class="wrapper">
			<div id="offers_table">
				<div class="wrap">
				<p class="price">
					<strong>125 zł</strong>
				</p>
				<a href="%s/ad">
					<strong>Ad Title</strong>
				</a>
				<table data-id="123123">
					<div>
					<i data-icon="clock"></i>
					dzisiaj %s
					</div>
				</table>
				</div>
			</div>
			</div>
				`, ts.URL, now.Format("15:04"))))
		})

		mux.HandleFunc("/ad", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			called = true
		})

		t.Run("if created_at < last_check_at", func(t *testing.T) {
			h.fetchSuggestions()
			if called {
				t.Error("visited")
			}
		})

		t.Run("if title is invalid", func(t *testing.T) {
			o.LastCheckAt = now.AddDate(0, 0, -1)
			o.Keywords = append(o.Keywords, models.Keyword{
				Value: "title",
				Type:  "excluded",
				For:   "title",
			})
			o.Checked = []models.Checked{}
			called = false
			h.fetchSuggestions()
			if called {
				t.Error("visited")
			}
		})
	})

	t.Run("should not store suggestion", func(t *testing.T) {
		mux := http.NewServeMux()
		ts := httptest.NewServer(mux)
		defer ts.Close()
		now := time.Now()
		storeCalled := false
		adPageVisited := false
		o := &models.Observation{
			Name:        "Random name",
			URL:         ts.URL + "/url",
			LastCheckAt: now,
			Keywords: []models.Keyword{
				models.Keyword{
					Type:  "excluded",
					For:   "description",
					Value: "opis",
				},
			},
		}
		observationMockRepo := &mocks.Repository{
			FetchFunc: func(*models.ObservationFilter) (models.PaginatedResponse, error) {
				return models.PaginatedResponse{
					Items: []*models.Observation{
						o,
					},
					Total: 1,
				}, nil
			},
			UpdateFunc: func(*models.Observation) error {
				return nil
			},
		}
		suggestionMockRepo := &_suggestionMocks.Repository{
			StoreFunc: func(*models.Suggestion) error {
				storeCalled = true
				return nil
			},
		}

		h := &handler{
			observationRepo: observationMockRepo,
			suggestionRepo:  suggestionMockRepo,
			collector:       newCollector(),
			logrus:          logrus.WithField("package", "cron_test"),
		}
		mux.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(fmt.Sprintf(`
			<div class="wrapper"></div>
			<div class="wrapper"></div>
			<div class="wrapper">
			<div id="offers_table">
				<div class="wrap">
				<p class="price">
					<strong>125 zł</strong>
				</p>
				<a href="%s/ad">
					<strong>Ad Title</strong>
				</a>
				<table data-id="123123">
					<div>
					<i data-icon="clock"></i>
					dzisiaj %s
					</div>
				</table>
				</div>
			</div>
			</div>
				`, ts.URL, now.Format("15:04"))))
		})

		mux.HandleFunc("/ad", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`
			<div id="textContent">
				Jakiś tam opis
			</div>
				`))
			adPageVisited = true
		})

		t.Run("if description is invalid", func(t *testing.T) {
			h.fetchSuggestions()
			if !adPageVisited {
				t.Error("didnt visit ad page")
			}
			if storeCalled {
				t.Error("stored")
			}
		})
	})

	t.Run("should successfully store suggestion", func(t *testing.T) {
		mux := http.NewServeMux()
		ts := httptest.NewServer(mux)
		defer ts.Close()
		now := time.Now()
		storeCalled := false
		adPageVisited := false
		notifyCalled := false
		o := &models.Observation{
			Name:        "Random name",
			URL:         ts.URL + "/url",
			LastCheckAt: now,
		}
		observationMockRepo := &mocks.Repository{
			FetchFunc: func(*models.ObservationFilter) (models.PaginatedResponse, error) {
				return models.PaginatedResponse{
					Items: []*models.Observation{
						o,
					},
					Total: 1,
				}, nil
			},
			UpdateFunc: func(*models.Observation) error {
				return nil
			},
		}
		suggestionMockRepo := &_suggestionMocks.Repository{
			StoreFunc: func(*models.Suggestion) error {
				storeCalled = true
				return nil
			},
		}
		notificationsMockManager := &_notificationsMocks.Manager{
			NotifyFunc: func(string) error {
				notifyCalled = true
				return nil
			},
		}
		h := &handler{
			observationRepo:      observationMockRepo,
			suggestionRepo:       suggestionMockRepo,
			notificationsManager: notificationsMockManager,
			collector:            newCollector(),
			logrus:               logrus.WithField("package", "cron_test"),
		}
		mux.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(fmt.Sprintf(`
			<div class="wrapper"></div>
			<div class="wrapper"></div>
			<div class="wrapper">
			<div id="offers_table">
				<div class="wrap">
				<p class="price">
					<strong>125 zł</strong>
				</p>
				<a href="%s/ad">
					<strong>Ad Title</strong>
				</a>
				<table data-id="123123">
					<div>
					<i data-icon="clock"></i>
					dzisiaj %s
					</div>
				</table>
				</div>
			</div>
			</div>
				`, ts.URL, now.Format("15:04"))))
		})

		mux.HandleFunc("/ad", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`
			<div id="textContent">
				Jakiś tam opis
			</div>
				`))
			adPageVisited = true
		})

		h.fetchSuggestions()
		if !adPageVisited {
			t.Error("didnt visit ad page")
		}
		if !storeCalled {
			t.Error("didnt save suggestion")
		}
		if !notifyCalled {
			t.Error("didnt call notificationsManager.Notify")
		}
	})
}

func newCollector() *colly.Collector {
	return colly.NewCollector(
		colly.Async(true),
		colly.AllowURLRevisit(),
	)
}
