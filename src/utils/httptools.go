package utils

import "fmt"
import "net/http"
import "encoding/json"

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
func (ht *HttpTools) WriteData(w http.ResponseWriter, data map[string]interface{}) {
	str, err := json.Marshal(data)

	if err != nil {
		//fmt.Println(err)
	}
	content := "{\"success\":true,\"data\":" + string(str) + "}"
	w.Write([]byte(content))
}

func (ht *HttpTools) WriteList(w http.ResponseWriter, data []map[string]interface{}) {
	str, err := json.Marshal(data)

	if err != nil {
		//fmt.Println(err)
	}
	content := "{\"success\":true,\"data\":" + string(str) + "}"
	w.Write([]byte(content))
}

func (ht *HttpTools) WriteObject(w http.ResponseWriter, data map[string]interface{}) {
	str, err := json.Marshal(data)

	if err != nil {
		//fmt.Println(err)
	}
	w.Write([]byte(string(str)))
}
