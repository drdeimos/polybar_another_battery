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
  "github.com/distatus/battery"
)

func main() {
  var state string
  var thr float64 = 10.0
  notify_init()

  // FIXME threshould hardcoded now
  threshould := os.Args[1:]
  fmt.Printf("Arg %v\n", threshould)

  for {
    batteries, err := battery.GetAll()
	  if err != nil {
      fmt.Println("Could not get battery info!")
	    return
	  }
    for i, battery := range batteries {
      fmt.Printf("Bat%d: ", i)
      fmt.Printf("state: %f\n", battery.State)
      switch ; battery.State {
      case 3:
        state = "Charging"
      case 4:
        state = "Discharging"
      default:
        state = "Unknown"
      }

      percent := battery.Current / (battery.Full * 0.01)
      if ; percent < thr && battery.State != 3 {
        body := "Charge percent: " + strconv.FormatFloat(percent, 'f', 2, 32) + "\nState: " + state
        notify_send("Battery low!", body, 1)
      }
      fmt.Printf("Charge percent: %.2f \n", percent)
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
