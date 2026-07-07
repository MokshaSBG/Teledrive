package main

/*
#cgo LDFLAGS: -ltdjson
#include <td/telegram/td_json_client.h>
*/
import "C"

import (
	"fmt"
	"teledrive/internal/telegram"
)

/*func main() {
	// This is a test call to see if the library is linked correctly
	client := C.td_json_client_create()
	if client != nil {
		fmt.Println("✅ Success! Go is talking to TDLib.")
	} else {
		fmt.Println("❌ Bridge failed. Check library installation.")
	}
}*/

func main() {
	tgClient := telegram.NewClient()

	fmt.Println("🚀 Starting TG-Drive Update Loop...")

	for {
		update := tgClient.Receive(1.0) // Wait 1 second for updates
		if update == "" {
			continue
		}

		// For now, let's just print every raw JSON we get from Telegram
		fmt.Printf("📩 New Update: %s\n", update)

		// Next step: We will write a function to parse these JSONs
		// and handle "authorizationStateWaitPhoneNumber"
	}
}
