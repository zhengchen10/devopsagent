package plugins

import "common"

type LinkFile struct {
	me     string
	global *common.Global
}

func (u *LinkFile) GetName() string {
	return u.me
}

func (u *LinkFile) InitPlugin(g *common.Global) {
	u.me = "LinkFile"
	u.global = g
	g.GetLog().InfoA(u.me, "InitPlugin")
	g.RegisterHandler("ln", u)
	//g.GetRouter().HandleFunc("/jstat", u.jstatHandler)
}

func (u *LinkFile) StartPlugin() {

}

func (u *LinkFile) StopPlugin() {

}

func (u *LinkFile) GetRequestParams() []string {
	var params []string
	params = append(params, "src")
	params = append(params, "dest")
	return params
}

func (u *LinkFile) Execute(params map[string]string) (map[string]interface{}, int) {
	src := params["src"]
	dest := params["dest"]
	ret, err := u.ExecuteWithParams(src, dest)
	if err != nil {
		return nil, -1
	}
	return ret, 0
}

func (u *LinkFile) ExecuteWithParams(src string, dest string) (map[string]interface{}, error) {
	context := make(map[string]interface{})
	//params := "-p " + pid + " |wc -l"
	//u.global.GetCmdTools().ExecuteWithCallback("pstree", params, context, true, u.ExecuteCallback)
	params := "-s " + src + " " + dest
	_, err := u.global.GetCmdTools().Execute("ln", params, true)
	if err == nil {

	}
	return context, nil
}
