package tray

import (
	"fmt"

	"github.com/farmsimcompanymanager/companion/internal/config"
	"github.com/farmsimcompanymanager/companion/internal/server"
	"github.com/getlantern/systray"
)

func Run(cfg *config.Config, srv *server.Server) {
	systray.Run(func() { onReady(cfg, srv) }, onExit)
}

func onReady(cfg *config.Config, srv *server.Server) {
	systray.SetTitle("FarmSim Manager")
	systray.SetTooltip("FarmSim Company Manager — running")

	mStatus := systray.AddMenuItem(fmt.Sprintf("Running on port %d", cfg.Port), "")
	mStatus.Disable()

	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "Stop the companion app")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {}
