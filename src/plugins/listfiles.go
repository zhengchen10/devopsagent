package plugins

import "common"
import "os/exec"
import "strings"

type ListFiles struct {
	me     string
	global *common.Global
}

func (u *ListFiles) GetName() string {
	return u.me
}

func (u *ListFiles) InitPlugin(g *common.Global) {
	u.me = "ListFiles"
	u.global = g
	g.GetLog().InfoA(u.me, "InitPlugin")
	//g.GetRouter().HandleFunc("/listfiles", u.listFilesHandler)
	g.RegisterHandler("listfiles", u)
}

func (u *ListFiles) StartPlugin() {

}

func (u *ListFiles) StopPlugin() {

}

func (h *ListFiles) GetRequestParams() []string {
	var params []string
	params = append(params, "path")
	return params
}

func (h *ListFiles) Execute(params map[string]string) (map[string]interface{}, int) {
	path := params["path"]
	ret, err := h.ExecuteWithParams(path)
	if err != nil {
		return nil, -1
	}
	result := make(map[string]interface{})
	result["files"] = ret
	return result, 0
}

/*func (u *ListFiles) listFilesHandler(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	ret, err := u.Execute(path)
	if err != nil {
		u.global.GetHttpTools().WriteError(w, -1)
	}
	u.global.GetHttpTools().WriteList(w, ret)
}*/

func (u *ListFiles) ExecuteWithParams(path string) ([]map[string]interface{}, error) {
	context := make(map[string]interface{})
	var list []map[string]interface{}
	context["dataList"] = list
	//params := "-p " + pid + " |wc -l"
	//u.global.GetCmdTools().ExecuteWithCallback("pstree", params, context, true, u.ExecuteCallback)
	params := path + " -F"
	u.global.GetCmdTools().ExecuteWithCallback("ls", params, context, true, u.ExecuteCallback)
	return context["dataList"].([]map[string]interface{}), nil
}

func (u *ListFiles) ExecuteCallback(cmd *exec.Cmd, line string, context interface{}) bool {
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
