package plugins

import "common"
import "strings"

type UptimeInfo struct {
	me     string
	global *common.Global
}

func (u *UptimeInfo) GetName() string {
	return u.me
}

func (u *UptimeInfo) InitPlugin(g *common.Global) {
	u.me = "UptimeInfo"
	u.global = g
	g.GetLog().InfoA(u.me, "InitPlugin")
	//g.GetRouter().HandleFunc("/uptime", u.uptimeHandler)
	g.RegisterHandler("uptime", u)
}

func (u *UptimeInfo) StartPlugin() {

}

func (u *UptimeInfo) StopPlugin() {

}

func (u *UptimeInfo) GetRequestParams() []string {
	var params []string
	return params
}

func (u *UptimeInfo) Execute(params map[string]interface{}) (map[string]interface{}, int) {
	ret, err := u.ExecuteWithParams()
	if err != nil {
		return nil, -1
	}
	return ret, 0
}

/*
func (u *UptimeInfo) uptimeHandler(w http.ResponseWriter, r *http.Request) {
	ret, err := u.Execute()
	if err != nil {
		u.global.GetHttpTools().WriteError(w, -1)
		return
	}

	u.global.GetHttpTools().WriteData(w, ret)
}*/

func (u *UptimeInfo) ExecuteWithParams() (map[string]interface{}, error) {
	var ret = make(map[string]interface{})
	out, err := u.global.GetCmdTools().Execute("uptime", "", true)
	if err == nil {
		startIndex := strings.LastIndex(out, "load average: ") + 14
		endIndex := len(out)
		stringTools := u.global.GetStringTools()
		substring := stringTools.SubString(out, startIndex, endIndex)
		items := strings.Split(substring, ",")
		ret["1min"] = stringTools.Trim(items[0])
		ret["5min"] = stringTools.Trim(items[1])
		ret["15min"] = stringTools.Trim(items[2])
		return ret, nil
	}

	return ret, err
}
