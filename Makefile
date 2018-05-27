VERCMD  ?= git describe --long --tags 2> /dev/null
VERSION ?= $(shell $(VERCMD) || cat VERSION)
BINNAME ?= "polybar-ab"

PREFIX    ?= /usr/local
BINPREFIX ?= $(PREFIX)/bin

all: build strip install

build:
	go build -ldflags "-X main.version=${VERSION}" -o ${BINNAME}

install:
	mkdir -p "$(DESTDIR)$(BINPREFIX)"
	cp -pf ${BINNAME} "$(DESTDIR)$(BINPREFIX)"
unistall:
	rm -rf "$(DESTDIR)$(BINPREFIX)/${BINNAME}"

strip:
	strip ${BINNAME}

clean:
	rm -rf ${BINNAME}
