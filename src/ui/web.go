package ui

import (
	"dk4/utils"
	"encoding/json"
	"fmt"

	"github.com/jchv/go-webview2"
)

type Response struct {
  Code int    `json:"code"`
	Error string `json:"error"`
	Data interface{} `json:"data"`
}

func dk4Action (action string, payload string) string {
	handle, ok := jsFuncMaps[action]
	if !ok {
		return "{\"error\":\"unknown action " + action + ", \"code\": 400}"
	}

	data := ""
	err := utils.Try(func() {
		data = handle(payload)
	})
	if err != nil {
		t.Reset()
		leadSeaman = nil
		currentOrg = nil
		fmt.Printf("run action %s error %s", action, err.Error())
		return "{\"error\":\"" + err.Error() + "\", \"code\": 500}"
	}

	return data
}
func injectJSSDK(w webview2.WebView) {
	w.Bind("dk4Action", dk4Action)
}

func makeErrorResponse(msg string, code int) string {
	resp := Response{Code: code, Error: msg}
	result, _ := json.Marshal(resp)
	return string(result)
}
func makeResponse(code int) Response {
	resp := Response{Code: code, Error: ""}
	return resp
}

func formatResponse(r *Response) string {
	result, err := json.Marshal(*r)
	if err != nil {
		return "{\"error\":\"" + err.Error() + ", \"code\": 400}"
	}
	return string(result)
}