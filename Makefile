all: build strip

build:
	go build -ldflags "-X main.version=${VERSION}" -o polybar_ab

strip:
	strip polybar_ab

clean:
	rm -rf polybar_ab
