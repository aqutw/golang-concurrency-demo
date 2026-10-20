// Exercise the real pprof HTTP handler in memory; no port or server to clean up.
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/pprof"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	request := httptest.NewRequest(http.MethodGet, "/debug/pprof/goroutine?debug=1", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "goroutine profile:") {
		panic("missing goroutine profile")
	}
	fmt.Println("pprof HTTP status:", response.Code)
	fmt.Println("goroutine profile captured")
}
