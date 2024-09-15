package bootstrap

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/clicklord/lms/config"
)

func LoadMainWindow(cfg *config.DmsConfig) fyne.Window {
	a := app.New()
	w := a.NewWindow("Local Media Server")

	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu("LMS",
			fyne.NewMenuItem("Show", func() {
				w.Show()
			}),
			fyne.NewMenuItem("Quit", func() {
				a.Quit()
			}))
		desk.SetSystemTrayMenu(m)
	}

	mainLabel := widget.NewLabel("Shared folder path: " + cfg.Path)
	w.SetContent(container.NewVBox(
		mainLabel,
		widget.NewButton("Hide to system tray", func() {
			w.Hide()
		}),
		widget.NewButton("Quit", func() {
			a.Quit()
		}),
	))
	w.Resize(fyne.NewSize(300, 100))
	w.SetCloseIntercept(func() {
		w.Hide()
	})

	return w
}
