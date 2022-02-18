VERCMD  ?= git describe --long --tags 2> /dev/null
VERSION ?= $(shell $(VERCMD) || cat VERSION)
BINNAME ?= "polybar-ab"

PREFIX    ?= /usr/local
BINPREFIX ?= $(PREFIX)/bin

all: build strip install

build:
	go build -ldflags "-X main.version=$(VERSION)" -o $$(pwd)/$(BINNAME)

install:
	install -D -m 755 -o root -g root $(BINNAME) $(DESTDIR)$(BINPREFIX)/$(BINNAME)

uninstall:
	rm -rf "$(DESTDIR)$(BINPREFIX)/$(BINNAME)"

strip:
	strip $(BINNAME)

clean:
	rm -rf $(BINNAME)
