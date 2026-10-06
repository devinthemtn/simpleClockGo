APP     := simpleclock
BIN     := bin
WINCC   := x86_64-w64-mingw32-gcc
MACCC   := o64-clang
PREFIX  ?= $(HOME)/.local

.PHONY: all linux windows mac clean install uninstall

all: linux

linux: | $(BIN)
	CGO_ENABLED=1 go build -o $(BIN)/$(APP) .

windows: | $(BIN)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=$(WINCC) go build -ldflags="-H windowsgui" -o $(BIN)/$(APP).exe .

mac: | $(BIN)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 CC=$(MACCC) go build -o $(BIN)/$(APP)-mac .

install: linux
	install -Dm755 $(BIN)/$(APP) $(PREFIX)/bin/$(APP)
	install -Dm644 assets/clockIcon.png $(PREFIX)/share/icons/hicolor/256x256/apps/$(APP).png
	mkdir -p $(PREFIX)/share/applications
	sed 's|^Exec=.*|Exec=$(PREFIX)/bin/$(APP)|' assets/$(APP).desktop > $(PREFIX)/share/applications/$(APP).desktop
	gtk-update-icon-cache -q $(PREFIX)/share/icons/hicolor 2>/dev/null || true
	update-desktop-database -q $(PREFIX)/share/applications 2>/dev/null || true

uninstall:
	rm -f $(PREFIX)/bin/$(APP) \
		$(PREFIX)/share/icons/hicolor/256x256/apps/$(APP).png \
		$(PREFIX)/share/applications/$(APP).desktop

$(BIN):
	mkdir -p $(BIN)

clean:
	rm -rf $(BIN)
