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
	g.RegisterPlugin(healthCheck)
}
