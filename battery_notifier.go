package main

// #cgo pkg-config: libnotify
// #include <stdio.h>
// #include <errno.h>
// #include <libnotify/notify.h>
import "C"
import (
  "fmt"
  "github.com/distatus/battery"
)

func main() {
  batteries, err := battery.GetAll()
	if err != nil {
		fmt.Println("Could not get battery info!")
		return
	}
	for i, battery := range batteries {
		fmt.Printf("Bat%d: ", i)
		//fmt.Printf("state: %f, ", battery.State)
		//fmt.Printf("current capacity: %f mWh, ", battery.Current)
		//fmt.Printf("last full capacity: %f mWh, ", battery.Full)
		//fmt.Printf("design capacity: %f mWh, ", battery.Design)
		//fmt.Printf("charge rate: %f mW, ", battery.ChargeRate)
		//fmt.Printf("voltage: %f V, ", battery.Voltage)
		//fmt.Printf("design voltage: %f V\n", battery.DesignVoltage)
    //fmt.Printf("\n")
    percent := battery.Current / (battery.Full * 0.01)
		fmt.Printf("Charge percent: %.2f \n", percent)

    //var n C.NotifyNotification;
    cs := C.CString("test")
    ret := C.notify_init(cs)
		fmt.Printf("%t\n", ret)
    n := C.notify_notification_new (cs, nil, nil)
    C.notify_notification_show(n, nil)
		fmt.Printf("End\n")
	}
}

