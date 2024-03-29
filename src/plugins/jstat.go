package plugins

import (
	"bytes"
	"common"
	"strconv"
)
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
	g.RegisterHandler(g.GetMessages().JSTAT_TEXT(), u)
	g.GetMessageCoder().RegisterDecoder(g.GetMessages().JSTAT(), 1, u)
	g.GetMessageCoder().RegisterEncoder(g.GetMessages().JSTAT(), 1, u)
}

func (u *JStat) StartPlugin() {

}

func (u *JStat) StopPlugin() {

}

/*
func (u *JStat) jstatHandler(w http.ResponseWriter, r *http.Request) {
	pid := r.FormValue("pid")
	ret, err := u.Execute(pid)
	if err != nil {
		u.global.GetHttpTools().WriteError(w, -1)
		return
	}
	//context := make(map[string]string)
	//context["count"] = ret
	u.global.GetHttpTools().WriteData(w, ret)
}*/

func (h *JStat) GetRequestParams() []string {
	var params []string
	params = append(params, "pid")
	return params
}

func (h *JStat) Execute(params map[string]interface{}) (map[string]interface{}, int) {
	pid := params["pid"]
	ret, err := h.ExecuteWithParams(pid.(string))
	if err != nil {
		return nil, -1
	}
	return ret, 0
}

func (u *JStat) ExecuteWithParams(pid string) (map[string]interface{}, error) {
	context := make(map[string]interface{})
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

func (r *JStat) Decode(messageId, version, msgType int, data []byte) map[string]interface{} {
	byteTools := r.global.GetByteTools()
	ret := make(map[string]interface{})
	pos := 0

	ret["pid"] = byteTools.ReadString(data, &pos)
	return ret
}

func (r *JStat) Encode(messageId, version, msgType int, msg map[string]interface{}) []byte {
	var buffer bytes.Buffer
	byteTools := r.global.GetByteTools()
	s0 := msg["S0"].(string)
	s0v, _ := strconv.ParseFloat(s0, 32)
	buffer.Write(byteTools.FloatToBytes((float32)(s0v)))

	s1 := msg["S1"].(string)
	s1v, _ := strconv.ParseFloat(s1, 32)
	buffer.Write(byteTools.FloatToBytes((float32)(s1v)))

	e := msg["E"].(string)
	ev, _ := strconv.ParseFloat(e, 32)
	buffer.Write(byteTools.FloatToBytes((float32)(ev)))

	o := msg["O"].(string)
	ov, _ := strconv.ParseFloat(o, 32)
	buffer.Write(byteTools.FloatToBytes((float32)(ov)))

	m := msg["M"].(string)
	mv, _ := strconv.ParseFloat(m, 32)
	buffer.Write(byteTools.FloatToBytes((float32)(mv)))

	ygc := msg["YGC"].(string)
	ygcv, _ := strconv.ParseInt(ygc, 10, 32)
	buffer.Write(byteTools.IntToBytes(int32(ygcv)))

	ygct := msg["YGCT"].(string)
	ygctv, _ := strconv.ParseFloat(ygct, 32)
	buffer.Write(byteTools.FloatToBytes((float32)(ygctv)))

	fgc := msg["FGC"].(string)
	fgcv, _ := strconv.ParseInt(fgc, 10, 32)
	buffer.Write(byteTools.IntToBytes(int32(fgcv)))

	fgct := msg["FGCT"].(string)
	fgctv, _ := strconv.ParseFloat(fgct, 32)
	buffer.Write(byteTools.FloatToBytes((float32)(fgctv)))

	gct := msg["GCT"].(string)
	gctv, _ := strconv.ParseFloat(gct, 32)
	buffer.Write(byteTools.FloatToBytes((float32)(gctv)))

	return buffer.Bytes()
}
