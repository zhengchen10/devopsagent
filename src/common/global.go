package common

import (
	"io"
	"net/http"
	"strings"
	"utils"
)
import "github.com/gorilla/mux"

type Global struct {
	log      *utils.Log
	config   *utils.Config
	router   *mux.Router
	random   *utils.Random
	http     *utils.HttpTools
	md5      *utils.MD5
	cmd      *utils.CmdTools
	strtool  *utils.StringTools
	ziptool  *utils.ZipTools
	plugins  map[string]AppPlugin
	handlers map[string]RequestHandler
}

func (g *Global) InitGlobal(r *mux.Router) {

	g.config = new(utils.Config)
	g.config.InitConfig("app.conf")

	logPath := g.config.GetProperty("logFilePath")
	logLevel := g.config.GetProperty("logLevel")
	logFileLevel := g.config.GetProperty("logFileLevel")
	g.log = new(utils.Log)
	g.log.InitLog(logPath+"app.log", logLevel, logFileLevel)

	g.random = new(utils.Random)
	g.random.InitRandom()
	g.http = new(utils.HttpTools)
	g.md5 = new(utils.MD5)
	g.cmd = new(utils.CmdTools)
	g.strtool = new(utils.StringTools)
	g.ziptool = new(utils.ZipTools)

	g.plugins = make(map[string]AppPlugin)
	g.handlers = make(map[string]RequestHandler)
	g.cmd.SetLogger(g.log)
	g.router = r
	g.log.InfoA("Global", "Init Global")
}

func (g *Global) RegisterPlugin(p AppPlugin) {
	g.log.InfoA("Global", "RegisterPlugin ["+p.GetName()+"]")
	g.plugins[p.GetName()] = p
}

func (g *Global) RegisterHandler(req string, h RequestHandler) {
	g.log.InfoA("Global", "Register Handler ["+h.GetName()+"] for URI [/"+req+"]")
	g.handlers["/"+req] = h
	g.router.HandleFunc("/"+req, g.handlerFunc)
}

func (g *Global) handlerFunc(w http.ResponseWriter, r *http.Request) {
	url := r.RequestURI
	pos := strings.Index(url, "?")
	if pos >= 0 {
		url = g.strtool.SubString(url, 0, pos)
	}
	handler := g.handlers[url]
	if handler == nil {
		io.WriteString(w, "{\"success\":false}")
		return
	}
	params := handler.GetRequestParams()
	reqParams := make(map[string]string)
	for _, p := range params {
		value := r.FormValue(p)
		reqParams[p] = value
	}

	result, err := handler.Execute(reqParams)
	if err != 0 {
		g.GetHttpTools().WriteError(w, err)
	} else {
		g.GetHttpTools().WriteData(w, result)
	}
}

func (g *Global) GetPlugin(name string) AppPlugin {
	return g.plugins[name]
}
func (g *Global) GetRouter() *mux.Router {
	return g.router
}

func (g *Global) GetLog() *utils.Log {
	return g.log
}

func (g *Global) GetConfig() *utils.Config {
	return g.config
}

func (g *Global) GetHttpTools() *utils.HttpTools {
	return g.http
}

func (g *Global) GetMD5() *utils.MD5 {
	return g.md5
}
func (g *Global) GetCmdTools() *utils.CmdTools {
	return g.cmd
}

func (g *Global) GetStringTools() *utils.StringTools {
	return g.strtool
}

func (g *Global) GetZipTools() *utils.ZipTools {
	return g.ziptool
}
