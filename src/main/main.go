// devopsagent project main.go
package main

// import "fmt"
import "net/http"
import "github.com/gorilla/mux"
import "common"
import "plugins"

func main() {
	global := new(common.Global)
	r := mux.NewRouter()
	global.InitGlobal(r)
	port := global.GetConfig().GetProperty("port")

	loadPlugins(global)
	http.Handle("/", r)
	global.GetLog().Info("Start Server at [" + port + "]")
	http.ListenAndServe(":"+port, nil)
}

func loadPlugins(g *common.Global) {
	healthCheck := new(plugins.HealthCheck)
	healthCheck.InitPlugin(g)

	uploadfile := new(plugins.UploadFile)
	uploadfile.InitPlugin(g)

	uptimeinfo := new(plugins.UptimeInfo)
	uptimeinfo.InitPlugin(g)

	pidinfo := new(plugins.GetPidInfo)
	pidinfo.InitPlugin(g)

	getthreads := new(plugins.GetProcessThreads)
	getthreads.InitPlugin(g)

	jstat := new(plugins.JStat)
	jstat.InitPlugin(g)

	g.RegisterPlugin(healthCheck)
	g.RegisterPlugin(uploadfile)
	g.RegisterPlugin(uptimeinfo)
	g.RegisterPlugin(pidinfo)
	g.RegisterPlugin(getthreads)
	g.RegisterPlugin(jstat)
}
