package bootstrap

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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

	configWindow := configWindow(cfg, a)

	w.SetContent(container.NewVBox(
		widget.NewButton("Configure", func() {
			configWindow.Show()
		}),
		widget.NewButton("Hide to system tray", func() {
			w.Hide()
		}),
		widget.NewButton("Quit", func() {
			a.Quit()
		}),
	))

	w.SetCloseIntercept(func() {
		w.Hide()
	})

	return w
}

func configWindow(cfg *config.DmsConfig, a fyne.App) fyne.Window {
	w := a.NewWindow("Config")

	pathInput := widget.NewEntry()
	pathInput.SetPlaceHolder("Enter or browse for a file...")
	pathInput.SetText(cfg.Path)

	friendlyNameInput := widget.NewEntry()
	friendlyNameInput.SetPlaceHolder("Enter browse name")
	friendlyNameInput.SetText(cfg.FriendlyName)

	ifNameInput := widget.NewEntry()
	ifNameInput.SetPlaceHolder("Enter broadcast interface name")
	ifNameInput.SetText(cfg.IfName)

	openFolderDialog := dialog.NewFolderOpen(
		func(reader fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			if reader != nil {
				// Update the text input with the selected file path
				pathInput.SetText(reader.Path())
			}
		},
		w,
	)

	browseButton := widget.NewButton("...", func() {
		openFolderDialog.Show()
	})

	saveButton := widget.NewButton("Save", func() {
		cfg.Path = pathInput.Text
		cfg.FriendlyName = friendlyNameInput.Text
		cfg.IfName = ifNameInput.Text

		err := cfg.Save()
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		dialog.ShowInformation("Success", "Configuration saved successfully!", w)
	})

	clearCacheButton := widget.NewButton("Clear local cache", func() {
		err := os.Remove(cfg.FFprobeCachePath)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		dialog.ShowInformation("Success", "Cache successfully cleared!", w)
	})

	pathContainer := container.NewGridWithColumns(
		3, widget.NewLabel("Path:"),
		pathInput,
		container.NewGridWithColumns(6, browseButton),
	)
	browseNameContainer := container.NewGridWithColumns(
		3, widget.NewLabel("Browse name:"),
		friendlyNameInput,
	)
	ifNameContainer := container.NewGridWithColumns(
		3,
		widget.NewLabel("Interface name:"),
		ifNameInput,
	)

	content := container.NewVBox(
		pathContainer,
		browseNameContainer,
		ifNameContainer,
		clearCacheButton,
		saveButton,
		widget.NewButton("Back", func() {
			w.Hide()
		}),
	)

	w.SetContent(content)

	return w
}
