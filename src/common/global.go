package common

import "utils"
import "github.com/gorilla/mux"

type Global struct {
	log     *utils.Log
	config  *utils.Config
	router  *mux.Router
	random  *utils.Random
	plugins map[string]AppPlugin
}

func (g *Global) InitGlobal(r *mux.Router) {

	g.random = new(utils.Random)
	g.random.InitRandom()
	g.plugins = make(map[string]AppPlugin)
	g.config = new(utils.Config)
	g.config.InitConfig("app.conf")

	g.log = new(utils.Log)
	g.log.InitLog("app.log")

	g.router = r
	g.log.InfoA("Global", "Init Global")
}

func (g *Global) RegisterPlugin(p AppPlugin) {
	g.log.InfoA("Global", "RegisterPlugin ["+p.GetName()+"]")
	g.plugins[p.GetName()] = p
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
