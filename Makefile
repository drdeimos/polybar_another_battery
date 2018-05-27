all: build strip

build:
	go build -ldflags "-X main.version=${VERSION}" -o go_battery_notifier

strip:
	strip go_battery_notifier

clean:
	rm -rf go_battery_notifier
