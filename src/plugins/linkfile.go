package plugins

import (
	"common"
)

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
	g.GetMessageCoder().RegisterDecoder(g.GetMessages().LINK_FILE(), 1, u)
	g.GetMessageCoder().RegisterEncoder(g.GetMessages().LINK_FILE(), 1, u)
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

func (u *LinkFile) Execute(params map[string]interface{}) (map[string]interface{}, int) {
	src := params["src"]
	dest := params["dest"]
	ret, err := u.ExecuteWithParams(src.(string), dest.(string))
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

func (r *LinkFile) Decode(messageId, version, msgType int, data []byte) map[string]interface{} {
	byteTools := r.global.GetByteTools()
	ret := make(map[string]interface{})

	srclen := byteTools.BytesToShort(data[0:2])
	src := data[2 : 2+srclen]
	destlen := byteTools.BytesToShort(data[2+srclen : 4+srclen])
	dest := data[4+srclen : 4+srclen+destlen]
	ret["src"] = string(src)
	ret["dest"] = string(dest)
	return ret
}

func (r *LinkFile) Encode(messageId, version, msgType int, msg map[string]interface{}) []byte {
	ret := make([]byte, 0)
	return ret
}
