package utils

import "fmt"
import "net/http"

type HttpTools struct {
}

func (ht *HttpTools) CheckHeader(w http.ResponseWriter, r *http.Request, secretString string) bool {
	s := ht.GetHeader(w, r, "Secret")
	if s != secretString {
		w.WriteHeader(403)
		ht.WriteError(w, -1)
		return false
	}
	return true
}

func (ht *HttpTools) GetHeader(w http.ResponseWriter, r *http.Request, headerName string) string {
	var headerValue = ""
	if len(r.Header) > 0 {
		for k, v := range r.Header {
			if k == headerName {
				headerValue = v[0]
			}
		}
	}
	return headerValue
}

func (ht *HttpTools) WriteError(w http.ResponseWriter, errorCode int) {
	content := fmt.Sprintf("{\"success\":false,\"errorCode\":%d}", errorCode)
	w.Write([]byte(content))
}

func (ht *HttpTools) WriteSuccess(w http.ResponseWriter) {
	w.Write([]byte("{\"success\":true}"))
}
