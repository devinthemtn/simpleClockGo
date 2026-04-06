APP     := simpleclock
BIN     := bin
WINCC   := x86_64-w64-mingw32-gcc
MACCC   := o64-clang

.PHONY: all linux windows mac clean

all: linux

linux: | $(BIN)
	CGO_ENABLED=1 go build -o $(BIN)/$(APP) .

windows: | $(BIN)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=$(WINCC) go build -o $(BIN)/$(APP).exe .

mac: | $(BIN)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 CC=$(MACCC) go build -o $(BIN)/$(APP)-mac .

$(BIN):
	mkdir -p $(BIN)

clean:
	rm -rf $(BIN)
