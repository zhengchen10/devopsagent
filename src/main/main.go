// devopsagent project main.go
package main

// import "fmt"
//import "net/http"
//import "github.com/gorilla/mux"
import (
	"common"
	"server"
)
import "plugins"

func main() {
	global := new(common.Global)
	//r := mux.NewRouter()
	global.InitGlobal()

	agent := global.GetConfig().GetProperty("agent")
	if agent == "TCP" {
		global.InitAgent(new(server.TcpAgent))
	} else {
		global.InitAgent(new(server.HttpAgent))
	}
	loadPlugins(global)

	global.StartAgent()
}

func loadPlugins(g *common.Global) {
	healthCheck := new(plugins.HealthCheck)
	healthCheck.InitPlugin(g)

	uploadfile := new(plugins.UploadFile)
	uploadfile.InitPlugin(g)

	downloadfile := new(plugins.DownloadFile)
	downloadfile.InitPlugin(g)

	uptimeinfo := new(plugins.UptimeInfo)
	uptimeinfo.InitPlugin(g)

	pidinfo := new(plugins.GetPidInfo)
	pidinfo.InitPlugin(g)

	getthreads := new(plugins.GetProcessThreads)
	getthreads.InitPlugin(g)

	jstat := new(plugins.JStat)
	jstat.InitPlugin(g)

	listfiles := new(plugins.ListFiles)
	listfiles.InitPlugin(g)

	tailfile := new(plugins.TailFile)
	tailfile.InitPlugin(g)

	linkfile := new(plugins.LinkFile)
	linkfile.InitPlugin(g)

	removefile := new(plugins.RemoveFile)
	removefile.InitPlugin(g)

	g.RegisterPlugin(healthCheck)
	g.RegisterPlugin(uploadfile)
	g.RegisterPlugin(downloadfile)
	g.RegisterPlugin(uptimeinfo)
	g.RegisterPlugin(pidinfo)
	g.RegisterPlugin(getthreads)
	g.RegisterPlugin(jstat)
	g.RegisterPlugin(listfiles)
	g.RegisterPlugin(tailfile)
	g.RegisterPlugin(linkfile)
	g.RegisterPlugin(removefile)

	if g.GetConfig().GetProperty("plugin.mysql") == "true" {
		mysql := new(plugins.MysqlClient)
		mysql.InitPlugin(g)
		g.RegisterPlugin(mysql)
	}

}
