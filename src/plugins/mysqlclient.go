package plugins

import (
	"bytes"
	"common"
	"os/exec"
)

type MysqlClient struct {
	me     string
	global *common.Global
}

func (u *MysqlClient) GetName() string {
	return u.me
}

func (u *MysqlClient) InitPlugin(g *common.Global) {
	u.me = "mysql"
	u.global = g
	g.GetLog().InfoA(u.me, "InitPlugin")
	//g.GetRouter().HandleFunc("/listfiles", u.listFilesHandler)

	agent := u.global.GetConfig().GetProperty("agent")
	if agent == "TCP" {
		g.GetMessageCoder().RegisterDecoder(g.GetMessages().MYSQL_CLIENT(), 1, u)
		g.GetMessageCoder().RegisterEncoder(g.GetMessages().MYSQL_CLIENT(), 1, u)
	}
	g.RegisterHandler(u.global.GetMessages().MYSQL_CLIENT_TEXT(), u)
}

func (u *MysqlClient) StartPlugin() {

}

func (u *MysqlClient) StopPlugin() {

}

func (h *MysqlClient) GetRequestParams() []string {
	var params []string
	return params
}

func (r *MysqlClient) Decode(messageId, version, msgType int, data []byte) map[string]interface{} {
	byteTools := r.global.GetByteTools()
	ret := make(map[string]interface{})
	pos := 0

	ret["userName"] = byteTools.ReadString(data, &pos)
	ret["databaseName"] = byteTools.ReadString(data, &pos)
	len := byteTools.BytesToShort(data, &pos)
	cmds := []string{}
	for i := 0; i < len; i++ {
		cmd := byteTools.ReadString(data, &pos)
		cmds = append(cmds, cmd)
	}
	ret["commands"] = cmds
	return ret
}

func (m *MysqlClient) Encode(messageId, version, msgType int, msg map[string]interface{}) []byte {
	var buffer bytes.Buffer
	byteTools := m.global.GetByteTools()
	buffer.Write(byteTools.ShortToBytes(1))
	return buffer.Bytes()
}
func (m *MysqlClient) Execute(input map[string]interface{}) (map[string]interface{}, int) {
	context := make(map[string]interface{})

	user := input["userName"]
	//db := input["databaseName"];
	pass := m.global.GetConfig().GetProperty("plugin.mysql." + user.(string) + ".password")
	params := "-u " + user.(string) + " -P" + pass
	client := m.global.GetConfig().GetProperty("plugin.mysql.client")
	m.global.GetCmdTools().ExecuteWithCallback(client, params, context, false, m.ExecuteCallback)
	result := make(map[string]interface{})
	return result, 0
}
func (u *MysqlClient) ExecuteCallback(cmd *exec.Cmd, line string, context interface{}) bool {
	if len(line) > 0 {
		c := context.(map[string]string)
		c["count"] = line
	}
	return false
}
