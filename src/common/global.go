package common

import "utils"
import "github.com/gorilla/mux"

type Global struct {
	log     *utils.Log
	config  *utils.Config
	router  *mux.Router
	random  *utils.Random
	http    *utils.HttpTools
	md5     *utils.MD5
	cmd     *utils.CmdTools
	strtool *utils.StringTools
	plugins map[string]AppPlugin
}

func (g *Global) InitGlobal(r *mux.Router) {

	g.config = new(utils.Config)
	g.config.InitConfig("app.conf")

	logPath := g.config.GetProperty("logPath")
	g.log = new(utils.Log)
	g.log.InitLog(logPath + "app.log")

	g.random = new(utils.Random)
	g.random.InitRandom()
	g.http = new(utils.HttpTools)
	g.md5 = new(utils.MD5)
	g.cmd = new(utils.CmdTools)
	g.strtool = new(utils.StringTools)
	g.plugins = make(map[string]AppPlugin)

	g.cmd.SetLogger(g.log)
	g.router = r
	g.log.InfoA("Global", "Init Global")
}

func (g *Global) RegisterPlugin(p AppPlugin) {
	g.log.InfoA("Global", "RegisterPlugin ["+p.GetName()+"]")
	g.plugins[p.GetName()] = p
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
