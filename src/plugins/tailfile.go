package plugins

import (
	"common"
)
import "net/http"
import "os/exec"
import "strings"
import "strconv"

type TailFile struct {
	me     string
	global *common.Global
}

func (u *TailFile) GetName() string {
	return u.me
}

func (u *TailFile) InitPlugin(g *common.Global) {
	u.me = "TailFile"
	u.global = g
	g.GetLog().InfoA(u.me, "InitPlugin")
	g.GetRouter().HandleFunc("/tailfile", u.tailFileHandler)
}

func (u *TailFile) StartPlugin() {

}

func (u *TailFile) StopPlugin() {

}

func (u *TailFile) tailFileHandler(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	lines := r.FormValue("lines")
	line := 50
	if len(lines) > 0 {
		line, err := strconv.Atoi(lines)
		if err != nil {
			line = 50
		}
		if line > 0 {

		}
	}
	ret, err := u.Execute(path, line)
	if err != nil {
		u.global.GetHttpTools().WriteError(w, -1)
		return
	}
	context := make(map[string]interface{})
	context["content"] = ret
	u.global.GetHttpTools().WriteData(w, context)
}

func (u *TailFile) Execute(path string, lines int) (string, error) {
	params := path + " -n " + strconv.Itoa(lines)
	ret, err := u.global.GetCmdTools().Execute("tail", params, true)
	return ret, err
}

func (u *TailFile) ExecuteCallback(cmd *exec.Cmd, line string, context interface{}) bool {
	if len(line) > 0 {
		c := context.(map[string]interface{})
		list := c["dataList"].([]map[string]interface{})

		item := make(map[string]interface{})
		if strings.HasSuffix(line, "/") {
			file := strings.TrimSuffix(line, "/")
			item["n"] = file
			item["d"] = 1
		} else if strings.HasSuffix(line, "*") {
			file := strings.TrimSuffix(line, "*")
			item["n"] = file
			item["d"] = 0
		} else if strings.HasSuffix(line, "@") {
			file := strings.TrimSuffix(line, "@")
			item["n"] = file
			item["d"] = 0
		} else {
			item["n"] = line
			item["d"] = 0
		}
		list = append(list, item)
		c["dataList"] = list
	}

	return false
}
