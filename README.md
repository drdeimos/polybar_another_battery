# About

Simple battery charge level watcher with notifications (libnotify)

# Requirements

- (Build) go1.9.2 (But it should work on earlier versions)
- (Build) libnotify-dev
- (Run) Font for battery indicator - 3270Medium NF
- (Run) libnotify4

# Packages

- [PPA for Ubuntu](https://launchpad.net/~drdeimosnn/+archive/ubuntu/survive-on-wm)

# Build manually

```
go get -u github.com/distatus/battery/cmd/battery
make build
```

# Usage

Run with key `-h` for get actual help
```
$ ./polybar-ab -h
Usage of ./polybar-ab:
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
exec = polybar-ab -polybar -thr 10
tail = true
```

# TODO
- [ ] ETA battery life when discharging
- [ ] Battery health level (based on full/design capacity)
