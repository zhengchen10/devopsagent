package plugins

import "common"
import "fmt"
import "os/exec"
import "strings"

type GetPidInfo struct {
	me     string
	global *common.Global
}

func (u *GetPidInfo) GetName() string {
	return u.me
}

func (u *GetPidInfo) InitPlugin(g *common.Global) {
	u.me = "PidInfo"
	u.global = g
	g.GetLog().InfoA(u.me, "InitPlugin")
	//g.GetRouter().HandleFunc("/pid", u.pidHandler)
	g.RegisterHandler("pid", u)
}

func (u *GetPidInfo) StartPlugin() {

}

func (u *GetPidInfo) StopPlugin() {

}

/*
func (u *GetPidInfo) pidHandler(w http.ResponseWriter, r *http.Request) {
	pname := r.FormValue("pname")
	ret, err := u.Execute(pname)
	if err != nil {
		u.global.GetHttpTools().WriteError(w, -1)
	}
	context := make(map[string]interface{})
	context["pid"] = ret
	u.global.GetHttpTools().WriteData(w, context)
}*/

func (h *GetPidInfo) GetRequestParams() []string {
	var params []string
	params = append(params, "pname")
	return params
}

func (h *GetPidInfo) Execute(params map[string]string) (map[string]interface{}, int) {
	pname := params["pname"]

	ret, err := h.ExecuteWithParams(pname)
	if err != nil {
		return nil, -1
	}
	result := make(map[string]interface{})
	result["pid"] = ret
	return result, 0

}

func (u *GetPidInfo) ExecuteWithParams(pname string) (string, error) {
	context := make(map[string]interface{})
	item := strings.Split(pname, ",") // 支持按名字查询进程，多名字按,分割
	params := "-ef"
	for i := 0; i < len(item); i++ {
		params = params + "|grep " + item[i] + " "
	}
	//params = params + "| awk '{print $2}'"
	out := u.global.GetCmdTools().ExecuteWithCallback("ps", params, context, true, u.ExecuteCallback)

	if out {

	}

	return context["pid"].(string), nil
}

func (u *GetPidInfo) ExecuteCallback(cmd *exec.Cmd, line string, context interface{}) bool {
	fmt.Println("Callback - " + line)
	c := context.(map[string]interface{})
	commd := u.global.GetStringTools().SubString(line, 48, len(line))
	fmt.Println("Callback cmd - " + commd)
	if strings.HasPrefix(commd, "sh -c ps -ef") == true {
		return false
	}
	if strings.HasPrefix(commd, "grep") == true {
		return false
	}

	item := strings.Split(line, " ")

	index := 0
	var value2 string
	var value7 string
	for i := 0; i < len(item); i++ {
		if len(item[i]) > 0 {
			index = index + 1
			if index == 2 {
				value2 = item[i]
			}
			if index == 7 {
				value7 = item[i] // check time field
				break
			}
		}
	}
	if len(value7) == 8 {
		c["pid"] = value2
	}
	return false
}
