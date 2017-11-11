package main

/* Build flag for C libnotify library binding.
   Trimming binary. Reduce out binary size.
   Bind C libraries.
*/

// #cgo pkg-config: libnotify
// #cgo LDFLAGS: -s
// #include <stdio.h>
// #include <errno.h>
// #include <libnotify/notify.h>
import "C"
import (
  "os"
  "fmt"
  "time"
  "strconv"
  "flag"
  "github.com/distatus/battery"
)

var flagdebug bool
var flagsimple bool
var flagpolybar bool
var flagonce bool
var flagthr int

func main() {
  var state string

  flag_init()
  notify_init()

  if ; flagdebug {
    fmt.Printf("Debug: flagthr=%v\n", flagthr)
  }

  for {
    batteries, err := battery.GetAll()
	  if err != nil {
      fmt.Println("Could not get battery info!")
	    return
	  }
    for i, battery := range batteries {
      if ; flagdebug {
        fmt.Printf("Bat%d:\n", i)
        fmt.Printf("  state: %v %f\n", battery.State, battery.State)
      }

      switch ; battery.State {
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

      if ; percent < float64(flagthr) && battery.State != 3 {
        body := "Charge percent: " + strconv.FormatFloat(percent, 'f', 2, 32) + "\nState: " + state
        notify_send("Battery low!", body, 1)
      }

      if ; flagdebug {
        fmt.Printf("  Charge percent: %.2f \n", percent)
        fmt.Printf("  Sleep sec: %v \n", 10)
        fmt.Printf("  Time: %v \n", time.Now())
      }

      if ; flagsimple {
        fmt.Printf("%.2f\n", percent)
      }
      if ; flagpolybar {
        polybar_out(percent, battery.State)
      }
      if ; flagonce {
        os.Exit(0)
      }
      time.Sleep(10 * time.Second)
	  }
  }
}

func notify_init() {
  cs := C.CString("test")
  ret := C.notify_init(cs)
  if ; ret != 1 {
    fmt.Printf("Notification init failed. Returned: %v\n", ret)
  }
}

func flag_init() {
  // wordPtr := flag.String("word", "foo", "a string")
  // numbPtr := flag.Int("numb", 42, "an int")
  flag.BoolVar(&flagdebug, "debug", false, "Enable debug output to stdout")
  flag.BoolVar(&flagsimple, "simple", false, "Print battery level to stdout every check")
  flag.BoolVar(&flagpolybar, "polybar", false, "Print battery level in polybar format")
  flag.BoolVar(&flagonce, "once", false, "Check state and print once")
  flag.IntVar(&flagthr, "thr", 10, "Set threshould battery level for notificcations")

  flag.Parse()

  if ; flagdebug {
    fmt.Println("Debug:", flagdebug)
    fmt.Println("tail:", flag.Args())
  }
}

func notify_send(summary, body string, urg int) {
  csummary := C.CString(summary)
  cbody := C.CString(body)
  var curg C.NotifyUrgency

  switch ; urg {
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
  if ; ret != 1 {
    fmt.Printf("Notification show failed. Returned: %v\n", ret)
  }
}

func polybar_out(val float64, state battery.State) {
  col_empty := "FFFFFF"
  col_full := "00FF00"
  col_charging := "444444"
  //col_charging := "FFDF00"
  col_discharging := "ADDFAD"
  col_default := "DFDFDF"
  // case 1:"Empty"
  // case 2:"Full"
  // case 3:"Charging"
  // case 4:"Discharging"

  switch ; state {
    case 1:
      fmt.Printf("%%{F#%v}  %%{F#%v}%.2f%%\n", col_empty, col_default, val)
    case 2:
      fmt.Printf("%%{F#%v}  %%{F#%v}%.2f%%\n", col_full, col_default, val)
    case 3:
      fmt.Printf("%%{F#%v}  %%{F#%v}%.2f%%\n", col_charging, col_default, val)
    case 4:
      fmt.Printf("%%{F#%v}  %%{F#%v}%.2f%%\n", col_discharging, col_default, val)
    }
}
