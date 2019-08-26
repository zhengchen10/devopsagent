package plugins

import (
	"common"
)

type RemoveFile struct {
	me     string
	global *common.Global
}

func (u *RemoveFile) GetName() string {
	return u.me
}

func (u *RemoveFile) InitPlugin(g *common.Global) {
	u.me = "RemoveFile"
	u.global = g
	g.GetLog().InfoA(u.me, "InitPlugin")
	g.RegisterHandler("rm", u)
}

func (u *RemoveFile) StartPlugin() {

}

func (u *RemoveFile) StopPlugin() {

}

func (u *RemoveFile) GetRequestParams() []string {
	var params []string
	params = append(params, "path")
	return params
}

func (u *RemoveFile) Execute(params map[string]interface{}) (map[string]interface{}, int) {
	path := params["path"]
	ret, _ := u.ExecuteWithParams(path.(string))
	if ret == nil {
		return nil, -1
	}
	return ret, 0
}

func (u *RemoveFile) ExecuteWithParams(path string) (map[string]interface{}, error) {
	context := make(map[string]interface{})
	if !u.global.GetFileTools().IsExist(path) {
		return nil, nil
	}
	params := "-f " + path
	if u.global.GetFileTools().IsDir(path) {
		params = "-rf " + path
	}
	_, err := u.global.GetCmdTools().Execute("rm", params, true)
	if err == nil {
		return context, nil
	}
	return nil, nil
}
