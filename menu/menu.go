package menu

import (
	"io/ioutil"
	"olx-crawler/utils"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/getlantern/systray"
)

type menu struct {
	url string
	ch  chan<- os.Signal
}

func New(url string, ch chan<- os.Signal) {
	m := &menu{
		url,
		ch,
	}
	go systray.Run(m.onReady, m.onExit)
}

func (m *menu) onReady() {
	icon, err := ioutil.ReadFile("icon.ico")
	if err != nil {
		logrus.Fatalf("Cannot read icon.ico file: %s", err.Error())
		return
	}
	systray.SetTemplateIcon(icon, icon)
	systray.SetTitle("Olx-Crawler")
	systray.SetTooltip("Pretty awesome")
	mOpenBrowser := systray.AddMenuItem("New OLX-Crawler window", "Open new OLX-Crawler window in the browser")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		for {
			select {
			case <-mOpenBrowser.ClickedCh:
				utils.OpenBrowser(m.url)
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func (m *menu) onExit() {
	m.ch <- os.Interrupt
}
