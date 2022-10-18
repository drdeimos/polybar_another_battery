package main

/* Build flag for C libnotify library binding.
   Trimming binary. Reduce out binary size.
   Bind C libraries.
*/

// #cgo pkg-config: libnotify
// #include <stdio.h>
// #include <errno.h>
// #include <libnotify/notify.h>
import "C"
import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/distatus/battery"
)

var batdetected bool
var flagdebug bool
var flagfont int
var flagonce bool
var flagpolybar bool
var flagsimple bool
var flagthr int
var flagversion bool

var version string

func main() {
	var state string

	flag_init()
	if flagversion {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}
	notify_init()

	if flagdebug {
		fmt.Printf("Debug: flagdebug=%v\n", flagdebug)
		fmt.Printf("       flagfont=%v\n", flagfont)
		fmt.Printf("       flagonce=%v\n", flagonce)
		fmt.Printf("       flagpolybar=%v\n", flagpolybar)
		fmt.Printf("       flagsimple=%v\n", flagsimple)
		fmt.Printf("       flagthr=%v\n", flagthr)
		fmt.Printf("       flagversion=%v\n", flagversion)
	}

	for {
		waitBat()
		batteries, err := battery.GetAll()
		if err != nil {
			if flagdebug {
				fmt.Println("Could not get battery info!")
				fmt.Printf("%+v\n", err)
				//return
			}
		}
		for i, battery := range batteries {
			if flagdebug {
				fmt.Printf("%s:\n", battery)
				fmt.Printf("Bat%d:\n", i)
				fmt.Printf("  state: %v %f\n", battery.State, battery.State)
			}

			switch battery.State {
			case 0:
				state = "Not charging"
			case 1:
				state = "Empty"
			case 2:
				state = "Full"
			case 3:
				state = "Charging"
			case 4:
				state = "Discharging"
			default:
				state = "Unknown"
			}

			percent := battery.Current / (battery.Full * 0.01)
			if percent > 100.0 {
				percent = 100.0
			} else if battery.Full == 0 { // Workaround. Sometime sysfs don't know full charge level. Dunno why
				percent = 100
			}

			if percent < float64(flagthr) && battery.State != 3 {
				body := "Charge percent: " + strconv.FormatFloat(percent, 'f', 2, 32) + "\nState: " + state
				notify_send("Battery low!", body, 1)
			}

			if flagdebug {
				fmt.Printf("  Charge percent: %.2f \n", percent)
				fmt.Printf("  Sleep sec: %v \n", 10)
				fmt.Printf("  Time: %v \n", time.Now())
			}

			if flagsimple {
				fmt.Printf("%.2f\n", percent)
			}
			if flagpolybar {
				polybar_out(percent, battery.State)
			}
			if flagonce {
				os.Exit(0)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func notify_init() {
	cs := C.CString("test")
	ret := C.notify_init(cs)
	if ret != 1 {
		fmt.Printf("Notification init failed. Returned: %v\n", ret)
	}
}

func flag_init() {
	flag.BoolVar(&flagdebug, "debug", false, "Enable debug output to stdout")
	flag.BoolVar(&flagsimple, "simple", false, "Print battery level to stdout every check")
	flag.BoolVar(&flagpolybar, "polybar", false, "Print battery level in polybar format")
	flag.BoolVar(&flagonce, "once", false, "Check state and print once")
	flag.IntVar(&flagthr, "thr", 10, "Set threshould battery level for notifications")
	flag.IntVar(&flagfont, "font", 1, "Set font numbler for polybar output")
	flag.BoolVar(&flagversion, "version", false, "Print version info and exit")

	flag.Parse()

	if flagdebug {
		fmt.Println("Debug:", flagdebug)
		fmt.Println("tail:", flag.Args())
	}
}

func notify_send(summary, body string, urg int) {
	csummary := C.CString(summary)
	cbody := C.CString(body)
	var curg C.NotifyUrgency

	switch urg {
	case 1:
		curg = C.NOTIFY_URGENCY_CRITICAL
	case 2:
		curg = C.NOTIFY_URGENCY_NORMAL
	case 3:
		curg = C.NOTIFY_URGENCY_LOW
	}
	n := C.notify_notification_new(csummary, cbody, nil)
	C.notify_notification_set_urgency(n, curg)
	ret := C.notify_notification_show(n, nil)
	if ret != 1 {
		fmt.Printf("Notification show failed. Returned: %v\n", ret)
	}
}

func polybar_out(val float64, state battery.State) {
	if flagdebug {
		fmt.Printf("Debug polybar: val=%v, state=%v\n", val, state)
	}

	bat_icons := []string{
		"\xef\x95\xb9",
		"\xef\x95\xba",
		"\xef\x95\xbb",
		"\xef\x95\xbc",
		"\xef\x95\xbd",
		"\xef\x95\xbe",
		"\xef\x95\xbf",
		"\xef\x96\x80",
		"\xef\x96\x81",
		"\xef\x95\xb8"}
	color_default := "DFDFDF"
	color := get_color(val)

	switch state {
	// Not charging
	case 0:
		level := val / 10
		fmt.Printf("%%{T%d}%%{F#%v} %s %%{F#%v}%%{T-}%.2f%%\n", flagfont, "00DDFF", bat_icons[int(level)], color_default, val)
		if flagdebug {
			fmt.Printf("Polybar discharge pict: %v\n", int(level))
		}
	// Empty
	case 1:
		fmt.Printf("%%{T%d}%%{F#%v} %v %%{F#%v}%%{T-}%.2f%%\n", flagfont, color, bat_icons[0], color_default, val)
	// Full
	case 2:
		fmt.Printf("%%{T%d}%%{F#%v} %v %%{F#%v}%%{T-}%.2f%%\n", flagfont, color, bat_icons[9], color_default, val)
	// Unknown, Charging
	case 3:
		for i := 0; i < 10; i++ {
			fmt.Printf("%%{T%d}%%{F#%v} %s %%{F#%v}%%{T-}%.2f%%\n", flagfont, color, bat_icons[i], color_default, val)
			time.Sleep(100 * time.Millisecond)
		}
	// Discharging
	case 4:
		level := val / 10
		fmt.Printf("%%{T%d}%%{F#%v} %s %%{F#%v}%%{T-}%.2f%%\n", flagfont, color, bat_icons[int(level)], color_default, val)
		if flagdebug {
			fmt.Printf("Polybar discharge pict: %v\n", int(level))
		}
	}
}

func get_color(val float64) string {
	var color string

	switch {
	case val <= 5.0:
		color = "FF0000"
	case val <= 10.0:
		color = "FF1A00"
	case val <= 15.0:
		color = "FF3500"
	case val <= 20.0:
		color = "FF5000"
	case val <= 25.0:
		color = "FF6B00"
	case val <= 30.0:
		color = "FF8600"
	case val <= 35.0:
		color = "FFA100"
	case val <= 40.0:
		color = "FFBB00"
	case val <= 45.0:
		color = "FFD600"
	case val <= 50.0:
		color = "FFF100"
	case val <= 55.0:
		color = "F1FF00"
	case val <= 60.0:
		color = "D6FF00"
	case val <= 65.0:
		color = "BBFF00"
	case val <= 70.0:
		color = "A1FF00"
	case val <= 75.0:
		color = "86FF00"
	case val <= 80.0:
		color = "6BFF00"
	case val <= 85.0:
		color = "50FF00"
	case val <= 90.0:
		color = "35FF00"
	case val <= 95.0:
		color = "1AFF00"
	case val <= 100.0:
		color = "00FF00"
	}

	if flagdebug {
		fmt.Printf("Selected color: %v", color)
	}

	return color
}

func waitBat() {
	batdetected = false
	for batdetected != true {
		_, err := os.Stat("/sys/class/power_supply/BAT0")
		if os.IsNotExist(err) {
			if flagdebug {
				fmt.Println("Could not find battery!")
			}
			if flagpolybar {
				polybar_out(0, 4)
			}
			if flagonce {
				os.Exit(0)
			}
			time.Sleep(1 * time.Second)
		} else {
			batdetected = true
		}
	}
}
