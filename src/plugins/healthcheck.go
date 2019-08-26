package plugins

import "common"

type HealthCheck struct {
}

func (h *HealthCheck) GetName() string {
	return "HealthCheck"
}

func (h *HealthCheck) InitPlugin(g *common.Global) {
	g.GetLog().InfoA("HealthCheck", "InitPlugin")
	g.RegisterHandler("healthCheck", h)
	//g.GetRouter().HandleFunc("healthCheck", h.healthCheckHandler)
}

func (h *HealthCheck) StartPlugin() {

}

func (h *HealthCheck) StopPlugin() {

}

func (h *HealthCheck) GetRequestParams() []string {
	var params []string
	return params
}

func (h *HealthCheck) Execute(params map[string]interface{}) (map[string]interface{}, int) {
	ret := make(map[string]interface{})
	return ret, 0
}

/*func (h *HealthCheck) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "{\"success\":true}")
}*/
