package plugins

import "common"
import "os/exec"

type GetProcessThreads struct {
	me     string
	global *common.Global
}

func (u *GetProcessThreads) GetName() string {
	return u.me
}

func (u *GetProcessThreads) InitPlugin(g *common.Global) {
	u.me = "GetProcessThreads"
	u.global = g
	g.GetLog().InfoA(u.me, "InitPlugin")
	g.RegisterHandler("threads", u)
	//g.GetRouter().HandleFunc("/threads", u.threadsHandler)
}

func (u *GetProcessThreads) StartPlugin() {

}

func (u *GetProcessThreads) StopPlugin() {

}
func (h *GetProcessThreads) GetRequestParams() []string {
	var params []string
	params = append(params, "pid")
	return params
}

func (h *GetProcessThreads) Execute(params map[string]interface{}) (map[string]interface{}, int) {
	pid := params["pid"]
	ret, err := h.ExecuteWithParams(pid.(string))
	if err != nil {
		return nil, -1
	}
	return ret, 0
}

func (u *GetProcessThreads) ExecuteWithParams(pid string) (map[string]interface{}, error) {
	context := make(map[string]interface{})
	//params := "-p " + pid + " |wc -l"
	//u.global.GetCmdTools().ExecuteWithCallback("pstree", params, context, true, u.ExecuteCallback)
	params := "hH p " + pid + "|wc -l"
	u.global.GetCmdTools().ExecuteWithCallback("ps", params, context, true, u.ExecuteCallback)
	return context, nil
}

func (u *GetProcessThreads) ExecuteCallback(cmd *exec.Cmd, line string, context interface{}) bool {
	if len(line) > 0 {
		c := context.(map[string]string)
		c["count"] = line
	}

	return false
}
