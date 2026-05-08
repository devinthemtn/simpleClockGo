package main

import (
	_ "embed"
	"flag"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"reflect"
	"time"
	"unsafe"

	glfw "github.com/go-gl/glfw/v3.3/glfw"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/yaml.v3"
)

//go:embed assets/clockIcon.png
var clockIconBytes []byte

type ClockConfig struct {
	Timezone string `yaml:"timezone"`
	Label    string `yaml:"label"`
}

type Config struct {
	Clocks []ClockConfig `yaml:"clocks"`
}

func loadConfig() Config {
	dir, err := os.UserConfigDir()
	if err != nil {
		return Config{}
	}
	data, err := os.ReadFile(filepath.Join(dir, "simpleclock", "config.yaml"))
	if err != nil {
		return Config{}
	}
	var cfg Config
	yaml.Unmarshal(data, &cfg)
	return cfg
}

type secondaryClock struct {
	loc       *time.Location
	timeText  *canvas.Text
	label     *widget.Label
	dayOffset *widget.Label
}

// dragWidget is a transparent overlay that moves the window when dragged.
// It is used when the window has no title bar.
type dragWidget struct {
	widget.BaseWidget
	win      fyne.Window
	viewport *glfw.Window // cached lazily
	anchorX  float64      // cursor X within window at mouse-down
	anchorY  float64
}

func newDragWidget(win fyne.Window) *dragWidget {
	d := &dragWidget{win: win}
	d.ExtendBaseWidget(d)
	return d
}

func (d *dragWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(canvas.NewRectangle(color.Transparent))
}

func (d *dragWidget) ensureViewport() {
	if d.viewport != nil {
		return
	}
	rv := reflect.ValueOf(d.win).Elem()
	vp := rv.FieldByName("viewport")
	if vp.IsValid() {
		d.viewport = *(**glfw.Window)(unsafe.Pointer(vp.UnsafeAddr()))
	}
}

// MouseDown captures the in-window cursor position so Dragged can keep it fixed.
func (d *dragWidget) MouseDown(_ *desktop.MouseEvent) {
	d.ensureViewport()
	if d.viewport != nil {
		d.anchorX, d.anchorY = d.viewport.GetCursorPos()
	}
}

func (d *dragWidget) MouseUp(_ *desktop.MouseEvent) {}

// Dragged repositions the window so the anchor point stays under the cursor.
// Using GetCursorPos()+GetPos() rather than the Fyne delta avoids the
// feedback loop where SetPos shifts the window-relative cursor position,
// causing the next delta to partially cancel the move.
func (d *dragWidget) Dragged(_ *fyne.DragEvent) {
	if d.viewport == nil {
		return
	}
	curX, curY := d.viewport.GetCursorPos()
	winX, winY := d.viewport.GetPos()
	d.viewport.SetPos(
		winX+int(curX)-int(d.anchorX),
		winY+int(curY)-int(d.anchorY),
	)
}

func (d *dragWidget) DragEnd() {}

func main() {
	noTitlebar := flag.Bool("no-titlebar", false, "launch without window titlebar")
	flag.Parse()

	a := app.New()
	icon := fyne.NewStaticResource("clockIcon.png", clockIconBytes)
	a.SetIcon(icon)

	var w fyne.Window
	if *noTitlebar {
		if drv, ok := a.Driver().(desktop.Driver); ok {
			w = drv.CreateSplashWindow()
			w.SetPadded(true)
		} else {
			w = a.NewWindow("Clock")
		}
	} else {
		w = a.NewWindow("Clock")
	}
	w.SetIcon(icon)

	timeText := canvas.NewText("", theme.ForegroundColor())
	timeText.TextSize = 72
	timeText.TextStyle = fyne.TextStyle{Monospace: true}
	timeText.Alignment = fyne.TextAlignCenter

	dateLabel := widget.NewLabel("")
	dateLabel.Alignment = fyne.TextAlignCenter

	cfg := loadConfig()

	var clocks []secondaryClock
	for _, cc := range cfg.Clocks {
		if cc.Timezone == "" {
			continue
		}
		loc, err := time.LoadLocation(cc.Timezone)
		if err != nil {
			continue
		}
		tt := canvas.NewText("", theme.ForegroundColor())
		tt.TextSize = 36
		tt.TextStyle = fyne.TextStyle{Monospace: true}
		tt.Alignment = fyne.TextAlignCenter

		lbl := cc.Label
		if lbl == "" {
			lbl = cc.Timezone
		}
		l := widget.NewLabel(lbl)
		l.Alignment = fyne.TextAlignCenter

		dayOffset := widget.NewLabel("")
		clocks = append(clocks, secondaryClock{loc: loc, timeText: tt, label: l, dayOffset: dayOffset})
	}

	update := func() {
		now := time.Now()
		timeText.Text = now.Format("15:04:05")
		timeText.Color = theme.ForegroundColor()
		timeText.Refresh()
		dateLabel.SetText(now.Format("Monday, January 2, 2006"))

		localY, localM, localD := now.Date()
		localMidnight := time.Date(localY, localM, localD, 0, 0, 0, 0, now.Location())
		for _, c := range clocks {
			t2 := now.In(c.loc)
			c.timeText.Text = t2.Format("15:04")
			c.timeText.Color = theme.ForegroundColor()
			c.timeText.Refresh()

			secY, secM, secD := t2.Date()
			secMidnight := time.Date(secY, secM, secD, 0, 0, 0, 0, c.loc)
			diff := int(secMidnight.Sub(localMidnight).Hours() / 24)
			switch {
			case diff > 0:
				c.dayOffset.SetText(fmt.Sprintf("+%d", diff))
			case diff < 0:
				c.dayOffset.SetText(fmt.Sprintf("%d", diff))
			default:
				c.dayOffset.SetText("")
			}
		}
	}

	update()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			fyne.Do(update)
		}
	}()

	vboxItems := []fyne.CanvasObject{
		container.NewStack(timeText),
		widget.NewSeparator(),
		dateLabel,
	}
	for _, c := range clocks {
		vboxItems = append(vboxItems,
			widget.NewSeparator(),
			container.NewCenter(container.NewHBox(c.label, c.timeText, c.dayOffset)),
		)
	}
	content := container.NewPadded(container.NewVBox(vboxItems...))

	if *noTitlebar {
		dragHandle := newDragWidget(w)
		w.SetContent(container.NewStack(content, dragHandle))
	} else {
		w.SetContent(content)
	}
	w.ShowAndRun()
}
