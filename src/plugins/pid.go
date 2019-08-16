package plugins

import "common"
import "net/http"
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
	g.GetRouter().HandleFunc("/pid", u.pidHandler)
}

func (u *GetPidInfo) StartPlugin() {

}

func (u *GetPidInfo) StopPlugin() {

}

func (u *GetPidInfo) pidHandler(w http.ResponseWriter, r *http.Request) {
	pname := r.FormValue("pname")
	ret, err := u.Execute(pname)
	if err != nil {
		u.global.GetHttpTools().WriteError(w, -1)
	}
	context := make(map[string]string)
	context["pid"] = ret
	u.global.GetHttpTools().WriteData(w, context)
}

func (u *GetPidInfo) Execute(pname string) (string, error) {
	context := make(map[string]string)
	item := strings.Split(pname, ",") // 支持按名字查询进程，多名字按,分割
	params := "-ef"
	for i := 0; i < len(item); i++ {
		params = params + "|grep " + item[i] + " "
	}
	//params = params + "| awk '{print $2}'"
	out := u.global.GetCmdTools().ExecuteWithCallback("ps", params, context, u.ExecuteCallback)

	if out {

	}
	for key, value := range context {
		fmt.Println(key + ":" + value)
	}
	return context["pid"], nil
}

func (u *GetPidInfo) ExecuteCallback(cmd *exec.Cmd, line string, context interface{}) bool {
	fmt.Println("Callback - " + line)
	c := context.(map[string]string)
	commd := u.global.GetStringTools().SubString(line, 48, len(line))
	fmt.Println("Callback cmd - " + commd)
	if strings.HasPrefix(commd, "sh -c ps -ef") == false {
		item := strings.Split(line, " ")
		index := 0
		for i := 0; i < len(item); i++ {
			if len(item[i]) > 0 {
				index = index + 1
				if index == 2 {
					c["pid"] = item[i]
					break
				}
			}
		}
	}
	return false
}
