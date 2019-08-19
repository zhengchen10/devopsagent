package plugins

import "common"
import "net/http"
import "os/exec"
import "strings"

type JStat struct {
	me     string
	global *common.Global
}

func (u *JStat) GetName() string {
	return u.me
}

func (u *JStat) InitPlugin(g *common.Global) {
	u.me = "JStat"
	u.global = g
	g.GetLog().InfoA(u.me, "InitPlugin")
	g.GetRouter().HandleFunc("/jstat", u.jstatHandler)
}

func (u *JStat) StartPlugin() {

}

func (u *JStat) StopPlugin() {

}

func (u *JStat) jstatHandler(w http.ResponseWriter, r *http.Request) {
	pid := r.FormValue("pid")
	ret, err := u.Execute(pid)
	if err != nil {
		u.global.GetHttpTools().WriteError(w, -1)
	}
	//context := make(map[string]string)
	//context["count"] = ret
	u.global.GetHttpTools().WriteData(w, ret)
}

func (u *JStat) Execute(pid string) (map[string]string, error) {
	context := make(map[string]string)
	//params := "-p " + pid + " |wc -l"
	//u.global.GetCmdTools().ExecuteWithCallback("pstree", params, context, true, u.ExecuteCallback)
	params := "-gcutil " + pid + " 1 1"
	ret, err := u.global.GetCmdTools().Execute("jstat", params, true)
	if err == nil {
		pos := strings.Index(ret, "\n")
		sub := u.global.GetStringTools().SubString(ret, pos+1, len(ret))
		items := strings.Split(sub, " ")
		index := 0
		for i := 0; i < len(items); i++ {
			if len(items[i]) > 0 {
				index = index + 1
				if index == 1 {
					context["S0"] = items[i]
				}
				if index == 2 {
					context["S1"] = items[i]
				}
				if index == 3 {
					context["E"] = items[i]
				}
				if index == 4 {
					context["O"] = items[i]
				}
				if index == 5 {
					context["M"] = items[i]
				}
				if index == 7 {
					context["YGC"] = items[i]
				}
				if index == 8 {
					context["YGCT"] = items[i]
				}
				if index == 9 {
					context["FGC"] = items[i]
				}
				if index == 10 {
					context["FGCT"] = items[i]
				}
				if index == 11 {
					context["GCT"] = u.global.GetStringTools().Trim(items[i]) // clear \n
				}
			}
		}
	}
	return context, nil
}

func (u *JStat) ExecuteCallback(cmd *exec.Cmd, line string, context interface{}) bool {
	if len(line) > 0 {
		c := context.(map[string]string)
		c["count"] = line
	}

	return false
}
