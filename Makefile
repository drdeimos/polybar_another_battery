all: build strip

build:
	go build -o go_battery_notifier

strip:
	strip go_battery_notifier

clean:
	rm -rf go_battery_notifier
