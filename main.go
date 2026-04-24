package main

import (
	_ "embed"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//go:embed assets/clockIcon.png
var clockIconBytes []byte

func main() {
	a := app.New()
	icon := fyne.NewStaticResource("clockIcon.png", clockIconBytes)
	a.SetIcon(icon)

	w := a.NewWindow("Clock")
	w.SetIcon(icon)

	timeText := canvas.NewText("", theme.ForegroundColor())
	timeText.TextSize = 72
	timeText.TextStyle = fyne.TextStyle{Monospace: true}
	timeText.Alignment = fyne.TextAlignCenter

	dateLabel := widget.NewLabel("")
	dateLabel.Alignment = fyne.TextAlignCenter

	update := func() {
		now := time.Now()
		timeText.Text = now.Format("15:04:05")
		timeText.Color = theme.ForegroundColor()
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

	// container.NewStack forces timeText to fill the full width so
	// TextAlignCenter works correctly instead of centering a narrow box
	content := container.NewPadded(
		container.NewVBox(
			container.NewStack(timeText),
			widget.NewSeparator(),
			dateLabel,
		),
	)

	w.SetContent(content)
	w.ShowAndRun()
}
