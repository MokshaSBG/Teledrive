package telegram

/*
#cgo LDFLAGS: -ltdjson
#include <td/telegram/td_json_client.h>
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

type Client struct {
	rawClient unsafe.Pointer
}

func NewClient() *Client {
	return &Client{
		rawClient: C.td_json_client_create(),
	}
}

// Send sends a JSON string request to TDLib
func (c *Client) Send(jsonQuery string) {
	query := C.CString(jsonQuery)
	defer C.free(unsafe.Pointer(query))
	C.td_json_client_send(c.rawClient, query)
}

// Receive waits for an update from TDLib (timeout in seconds)
func (c *Client) Receive(timeout float64) string {
	result := C.td_json_client_receive(c.rawClient, C.double(timeout))
	if result == nil {
		return ""
	}
	return C.GoString(result)
}
