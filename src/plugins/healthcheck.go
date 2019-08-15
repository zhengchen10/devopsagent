package plugins

import "common"
import "net/http"
import "io"

type HealthCheck struct {
}

func (h *HealthCheck) GetName() string {
	return "HealthCheck"
}

func (h *HealthCheck) InitPlugin(g *common.Global) {
	g.GetLog().InfoA("HealthCheck", "InitPlugin")
	g.GetRouter().HandleFunc("/healthCheck", h.healthCheckHandler)
}

func (h *HealthCheck) StartPlugin() {

}

func (h *HealthCheck) StopPlugin() {

}

func (h *HealthCheck) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "{\"success\":true}")
}
