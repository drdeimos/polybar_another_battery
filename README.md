# About

Simple battery charge level watcher with libnotify usage

# Requirements

- (Build) Tested on go1.9.2 (But it should work on earlier versions)
- (Build) For notifications: libnotify-dev
- (Usage) As polybar battery indicator - [Siji iconic bitmap font](https://github.com/stark/siji)

# Build

```
go get -u github.com/distatus/battery/cmd/battery
make build
```

# Usage

Run with key `-h` for get actual help
```
$ ./polybar_ab -h
Usage of ./polybar_ab:
  -debug
      Enable debug output to stdout
  -once
      Check state and print once
  -polybar
      Print battery level in polybar format
  -simple
      Print battery level to stdout every check
  -thr int
      Set threshould battery level for notificcations (default 10)
  -version
      Print version info and exit
```

## Polybar

Built in [polybar](https://github.com/jaagr/polybar) support.
Add flag `-polybar` for get stdout output in polybar format:
![Charging](/screenshots/charging.gif?raw=true "Charging")

### Polybar module example
```
[module/custom-battery]
type = custom/script
exec = ../polybar_ab -polybar
tail = true
```

# TODO
- [ ] ETA battery life when discharging
- [ ] Battery health level (based on full/design capacity)
- [ ] Improve my English
