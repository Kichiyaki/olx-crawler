package cron

import (
	"fmt"
	"net"
	"net/http"
	"olx-crawler/colly/debug"
	"olx-crawler/config"
	"olx-crawler/models"
	"olx-crawler/notifications"
	"olx-crawler/observation"
	"olx-crawler/suggestion"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gocolly/colly/v2/storage"

	"github.com/goodsign/monday"

	"github.com/fsnotify/fsnotify"

	"github.com/gocolly/colly/v2/proxy"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/robfig/cron/v3"
)

type handler struct {
	notificationsManager notifications.Manager
	observationRepo      observation.Repository
	suggestionRepo       suggestion.Repository
	configManager        config.Manager
	collyStorage         storage.Storage
	collector            *colly.Collector
	logrus               *logrus.Entry
}

type Config struct {
	NotificationsManager notifications.Manager
	ObservationRepo      observation.Repository
	SuggestionRepo       suggestion.Repository
	ConfigManager        config.Manager
	CollyStorage         storage.Storage
}

func AttachHandlers(c *cron.Cron, cfg *Config) error {
	globalCfg, err := cfg.ConfigManager.Config()
	if err != nil {
		return err
	}
	collector, err := getCollector(cfg.CollyStorage, globalCfg)
	if err != nil {
		return err
	}
	h := &handler{
		cfg.NotificationsManager,
		cfg.ObservationRepo,
		cfg.SuggestionRepo,
		cfg.ConfigManager,
		cfg.CollyStorage,
		collector,
		logrus.WithField("package", "cron"),
	}
	cfg.ConfigManager.OnConfigChange(h.handleConfigChange)

	c.AddFunc("@every 1m", h.fetchSuggestions)

	return nil
}

func (h *handler) fetchSuggestions() {
	var mutex sync.Mutex
	pagination, err := h.observationRepo.Fetch(&models.ObservationFilter{
		Started: "true",
	})
	if err != nil {
		return
	}
	observations, _ := pagination.Items.([]*models.Observation)
	collector := h.collector.Clone()
	suggestions := make(map[string]*models.Suggestion)
	currentObservation := &models.Observation{}

	collector.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("url", r.URL.String())
	})

	collector.OnHTML("#textContent", func(e *colly.HTMLElement) {
		suggestion, ok := suggestions[e.Request.Ctx.Get("url")]
		if !ok {
			return
		}
		description := strings.Trim(e.Text, "")
		if isValid(currentObservation.Keywords, description, "description") {
			if err := h.suggestionRepo.Store(suggestion); err != nil {
				return
			}
			h.notificationsManager.Notify(fmt.Sprintf("%s: Hejho, powinieneś się tym zainteresować %s",
				currentObservation.Name,
				suggestion.URL))
			h.logrus.WithField("suggestion_id", suggestion.ID).Debug("new suggestion")
		}

		o := &models.Observation{
			Checked: []models.Checked{
				models.Checked{
					Value: suggestion.URL,
				},
			},
		}
		o.ID = currentObservation.ID
		if err := h.observationRepo.Update(o); err != nil {
			return
		}
		mutex.Lock()
		currentObservation.Checked = append(currentObservation.Checked, o.Checked[0])
		mutex.Unlock()
	})

	collector.OnHTML(".wrap", func(e *colly.HTMLElement) {
		s := parseHTMLElementToSuggestionStruct(e)
		if _, ok := suggestions[s.URL]; ok {
			return
		}
		for _, checked := range currentObservation.Checked {
			if checked.Value == s.URL {
				return
			}
		}
		s.ObservationID = currentObservation.ID

		if isValid(currentObservation.Keywords, s.Title, "title") {
			mutex.Lock()
			suggestions[s.URL] = s
			mutex.Unlock()
			collector.Visit(s.URL)
			return
		}

		o := &models.Observation{
			Checked: []models.Checked{
				models.Checked{
					Value: s.URL,
				},
			},
		}
		o.ID = currentObservation.ID
		if err := h.observationRepo.Update(o); err != nil {
			return
		}
		mutex.Lock()
		currentObservation.Checked = append(currentObservation.Checked, o.Checked[0])
		mutex.Unlock()
	})

	collector.OnHTML(`.wrapper:nth-child(3)`, func(e *colly.HTMLElement) {
		date := strings.TrimSpace(e.DOM.Find("#offers_table .wrap").Last().Find(`i[data-icon="clock"]`).Parent().Text())
		if isAfter(currentObservation.LastCheckAt, e.Request.Ctx.Get("url"), date) {
			if href, exists := e.DOM.Find(`a[data-cy="page-link-next"]`).Attr("href"); exists {
				e.Request.Visit(href)
			}
		}
	})

	for _, observation := range observations {
		h.logrus.WithField("observation", observation.Name).Info("Fetching suggestions...")
		currentObservation = observation
		collector.Visit(observation.URL)
		collector.Wait()
		currentObservation.LastCheckAt = time.Now()
		h.observationRepo.Update(currentObservation)
	}

}

func (h *handler) handleConfigChange(fsnotify.Event) {
	cfg, err := h.configManager.Config()
	if err != nil {
		return
	}
	c, err := getCollector(h.collyStorage, cfg)
	if err != nil {
		h.collector = c
	}
}

func getCollector(s storage.Storage, cfg *models.Config) (*colly.Collector, error) {
	collector := colly.NewCollector(
		colly.Debugger(&debug.LogrusDebugger{}),
		colly.Async(true),
		colly.AllowURLRevisit(),
	)
	collector.DisableCookies()
	extensions.RandomMobileUserAgent(collector)
	proxiesLength := len(cfg.Proxies)
	if cfg.Colly.Limit > 0 {
		limit := cfg.Colly.Limit
		if proxiesLength > 0 {
			limit *= proxiesLength
		}
		numCPU := runtime.NumCPU()
		if limit > numCPU*10 {
			limit = numCPU * 10
		}
		collector.Limit(&colly.LimitRule{
			DomainGlob:  "*",
			Parallelism: limit,
			RandomDelay: time.Duration(cfg.Colly.Delay) * time.Second,
		})
	}
	if s != nil {
		if err := collector.SetStorage(s); err != nil {
			return nil, err
		}
	}
	if proxiesLength > 0 {
		transport, err := getHTTPTransport(cfg.Proxies)
		if err != nil {
			return nil, err
		}
		collector.WithTransport(transport)
	}
	return collector, nil
}

func getHTTPTransport(proxies []string) (*http.Transport, error) {
	rp, err := proxy.RoundRobinProxySwitcher(proxies...)
	if err != nil {
		return nil, err
	}
	return &http.Transport{
		Proxy: rp,
		DialContext: (&net.Dialer{
			DualStack: true,
		}).DialContext,
		MaxIdleConns: 100,
		TLSNextProto: nil,
	}, nil
}

func parseHTMLElementToSuggestionStruct(e *colly.HTMLElement) *models.Suggestion {
	href, _ := e.DOM.Find("a strong").Parent().Attr("href")
	href = strings.TrimSpace(href)
	splitted := strings.Split(href, "#")
	href = splitted[0]
	price := e.DOM.Find(".price strong").Text()
	price = strings.TrimSpace(price)
	title := e.DOM.Find("a strong").Text()
	title = strings.TrimSpace(title)
	id, _ := e.DOM.Find("table").Attr("data-id")
	id = strings.TrimSpace(id)
	imgElement := e.DOM.Find("img")
	img := ""
	if imgElement != nil {
		img, _ = e.DOM.Find("img").Attr("src")
		img = strings.TrimSpace(img)
		img = strings.Replace(img, "261x203", "644x461", 1)
	}
	return &models.Suggestion{
		URL:   href,
		Title: title,
		OlxID: id,
		Image: img,
		Price: price,
	}
}

func isValid(keywords []models.Keyword, text, f string) bool {
	countExcluded := 0
	countOneOf := 0
	countTotalOneOf := 0

	for _, keyword := range keywords {
		if keyword.Type == "excluded" && keyword.For == f && strings.Contains(text, keyword.Value) {
			countExcluded++
			break
		} else if keyword.Type == "one_of" && keyword.For == f {
			countTotalOneOf++
			if strings.Contains(text, keyword.Value) {
				countOneOf++
			}
		}
	}
	return countExcluded == 0 && ((countTotalOneOf > 0 && countOneOf > 0) || countTotalOneOf == 0)
}

func isAfter(t time.Time, url, text string) bool {
	if strings.Contains(url, "olx.pl") {
		if strings.Contains(text, "wczoraj") && isTodayDate(t) {
			return false
		}
		if (strings.Contains(text, "dzisiaj") && !isTodayDate(t)) ||
			(strings.Contains(text, "wczoraj") && !isYestardayDate(t)) {
			return true
		}
		if strings.Contains(text, "dzisiaj") || strings.Contains(text, "wczoraj") {
			layout := "dzisiaj 15:04"
			if strings.Contains(text, "wczoraj") {
				layout = "wczoraj 15:04"
			}
			parsed, err := monday.ParseInLocation(layout, strings.ToLower(text), time.UTC, monday.LocalePlPL)
			if err != nil {
				return false
			}
			if parsed.Hour() > t.Hour() || (parsed.Hour() == t.Hour() && parsed.Minute() >= t.Minute()) {
				return true
			}
		} else {
			parsed, err := monday.ParseInLocation("2 Jan", strings.ToLower(text), time.UTC, monday.LocalePlPL)
			if err != nil {
				return false
			}
			if (parsed.Month() == t.Month() && parsed.Day() >= t.Day()) || parsed.Month() > t.Month() {
				return true
			}
		}
	}

	return false
}

func isTodayDate(t time.Time) bool {
	now := time.Now()
	return t.Month() == now.Month() && t.Day() == now.Day() && t.Year() == now.Year()
}

func isYestardayDate(t time.Time) bool {
	yesterday := time.Now().AddDate(0, 0, -1)
	return t.Month() == yesterday.Month() && t.Day() == yesterday.Day() && t.Year() == yesterday.Year()
}
