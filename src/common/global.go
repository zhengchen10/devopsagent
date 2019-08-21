package common

import (
	"utils"
)

type Global struct {
	log    *utils.Log
	config *utils.Config
	//router   *mux.Router
	random   *utils.Random
	http     *utils.HttpTools
	md5      *utils.MD5
	cmd      *utils.CmdTools
	strtool  *utils.StringTools
	ziptool  *utils.ZipTools
	filetool *utils.File
	plugins  map[string]AppPlugin
	agent    AppServer
}

func (g *Global) InitGlobal() {

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
	g.filetool = new(utils.File)

	g.plugins = make(map[string]AppPlugin)

	g.cmd.SetLogger(g.log)
	//g.router = r
	g.log.InfoA("Global", "Init Global")
}

func (g *Global) InitAgent(appServer AppServer) {
	g.agent = appServer
	g.agent.InitServer(g)
}

func (g *Global) StartAgent() {
	g.agent.(AppServer).StartServer()
}

func (g *Global) RegisterPlugin(p AppPlugin) {
	g.log.InfoA("Global", "RegisterPlugin ["+p.GetName()+"]")
	g.plugins[p.GetName()] = p
}

func (g *Global) RegisterHandler(req string, h RequestHandler) {
	g.agent.(AppServer).RegisterHandler(req, h)
	//g.handlers["/"+req] = h
	//g.router.HandleFunc("/"+req, g.handlerFunc)
}

func (g *Global) GetPlugin(name string) AppPlugin {
	return g.plugins[name]
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

func (g *Global) GetFileTools() *utils.File {
	return g.filetool
}

func (g *Global) GetAppServer() AppServer {
	return g.agent
}
