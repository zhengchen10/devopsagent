package server

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strings"
)
import "common"

type HttpAgent struct {
	me       string
	global   *common.Global
	router   *mux.Router
	handlers map[string]common.RequestHandler
}

func (s *HttpAgent) InitServer(g *common.Global) {
	s.router = mux.NewRouter()
	s.me = "HttpAgent"
	s.global = g
	s.handlers = make(map[string]common.RequestHandler)
}
func (s *HttpAgent) StartServer() {
	port := s.global.GetConfig().GetProperty("port")
	http.Handle("/", s.router)
	s.global.GetLog().Info("Start Http Agent at [" + port + "]")
	http.ListenAndServe(":"+port, nil)
}
func (s *HttpAgent) StopServer() {

}

func (s *HttpAgent) GetRouter() *mux.Router {
	return s.router
}

func (s *HttpAgent) RegisterHandler(req string, h common.RequestHandler) {
	s.global.GetLog().InfoA("Global", "Register Handler ["+h.GetName()+"] for URI [/"+req+"]")
	s.handlers["/"+req] = h
	s.router.HandleFunc("/"+req, s.handlerFunc)
}

func (s *HttpAgent) handlerFunc(w http.ResponseWriter, r *http.Request) {
	url := r.RequestURI
	pos := strings.Index(url, "?")
	if pos >= 0 {
		url = s.global.GetStringTools().SubString(url, 0, pos)
	}
	handler := s.handlers[url]
	if handler == nil {
		io.WriteString(w, "{\"success\":false}")
		return
	}
	params := handler.GetRequestParams()
	reqParams := make(map[string]interface{})
	for _, p := range params {
		value := r.FormValue(p)
		reqParams[p] = value
	}

	result, err := handler.Execute(reqParams)
	if err != 0 {
		s.global.GetHttpTools().WriteError(w, err)
	} else {
		s.global.GetHttpTools().WriteData(w, result)
	}
}

func (s *HttpAgent) Type() string {
	return "HTTP"
}
