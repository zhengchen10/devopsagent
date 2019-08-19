package plugins

import "common"
import "net/http"
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
	g.GetRouter().HandleFunc("/threads", u.threadsHandler)
}

func (u *GetProcessThreads) StartPlugin() {

}

func (u *GetProcessThreads) StopPlugin() {

}

func (u *GetProcessThreads) threadsHandler(w http.ResponseWriter, r *http.Request) {
	pid := r.FormValue("pid")
	ret, err := u.Execute(pid)
	if err != nil {
		u.global.GetHttpTools().WriteError(w, -1)
	}
	context := make(map[string]string)
	context["count"] = ret
	u.global.GetHttpTools().WriteData(w, context)
}

func (u *GetProcessThreads) Execute(pid string) (string, error) {
	context := make(map[string]string)
	//params := "-p " + pid + " |wc -l"
	//u.global.GetCmdTools().ExecuteWithCallback("pstree", params, context, true, u.ExecuteCallback)
	params := "hH p " + pid + "|wc -l"
	u.global.GetCmdTools().ExecuteWithCallback("ps", params, context, true, u.ExecuteCallback)
	return context["count"], nil
}

func (u *GetProcessThreads) ExecuteCallback(cmd *exec.Cmd, line string, context interface{}) bool {
	if len(line) > 0 {
		c := context.(map[string]string)
		c["count"] = line
	}

	return false
}
