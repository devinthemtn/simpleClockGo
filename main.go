package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Clock")
	w.Resize(fyne.NewSize(320, 150))
	w.SetFixedSize(true)

	timeText := canvas.NewText("", theme.ForegroundColor())
	timeText.TextSize = 64
	timeText.TextStyle = fyne.TextStyle{Monospace: true}
	timeText.Alignment = fyne.TextAlignCenter

	dateLabel := widget.NewLabel("")
	dateLabel.Alignment = fyne.TextAlignCenter

	update := func() {
		now := time.Now()
		timeText.Text = now.Format("15:04:05")
		timeText.Refresh()
		dateLabel.SetText(now.Format("Monday, January 2, 2006"))
	}

	update()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			fyne.Do(update)
		}
	}()

	w.SetContent(container.NewVBox(
		container.NewPadded(timeText),
		dateLabel,
	))
	w.ShowAndRun()
}
