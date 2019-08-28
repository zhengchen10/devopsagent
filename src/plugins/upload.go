package plugins

import (
	"bytes"
	"common"
	"net"
	"strconv"
)
import "server"
import "net/http"
import "io"
import "os"
import "utils"
import "fmt"

type UploadFile struct {
	authSecret  string
	uploadPath  string
	global      *common.Global
	uploadFiles map[string]interface{}
}

func (u *UploadFile) GetName() string {
	return "UploadFile"
}

func (u *UploadFile) InitPlugin(g *common.Global) {
	u.global = g
	g.GetLog().InfoA("UploadFile", "InitPlugin")

	//g.GetRouter().HandleFunc("/upload", u.uploadHandler)
	agent := u.global.GetConfig().GetProperty("agent")
	if agent == "TCP" {
		g.RegisterHandler("upload", u)
		g.RegisterHandler("uploadpackage", u)
		g.GetMessageCoder().RegisterDecoder(g.GetMessages().UPLOAD_FILE(), 1, u)
		g.GetMessageCoder().RegisterEncoder(g.GetMessages().UPLOAD_FILE(), 1, u)
		g.GetMessageCoder().RegisterDecoder(g.GetMessages().UPLOAD_PACKAGE(), 1, u)
		g.GetMessageCoder().RegisterEncoder(g.GetMessages().UPLOAD_PACKAGE(), 1, u)
		//global.InitAgent(new (server.TcpAgent))
	} else {
		u.global.GetAppServer().(*server.HttpAgent).GetRouter().HandleFunc("/upload", u.uploadHandler)
	}
	u.authSecret = g.GetConfig().GetProperty("auth")
	u.uploadPath = g.GetConfig().GetProperty("uploadPath")
	u.uploadFiles = make(map[string]interface{})
}

func (u *UploadFile) StartPlugin() {

}

func (u *UploadFile) StopPlugin() {

}

func (u *UploadFile) uploadHandler(w http.ResponseWriter, r *http.Request) {
	check := u.global.GetHttpTools().CheckHeader(w, r, u.authSecret)
	if check == false {
		return
	}

	r.ParseMultipartForm(128 << 20) // max memory is set to 32MB
	clientfd, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		u.global.GetHttpTools().WriteError(w, -2)
		return
	}
	if handler == nil {

	}
	defer clientfd.Close()

	versionStr := r.FormValue("version")
	pool := r.FormValue("pool")
	var localpath string
	if len(pool) != 0 {
		workpath := u.global.GetConfig().GetProperty(pool + ".workpath")
		localpath = fmt.Sprintf("%s/%s_%s.jar", workpath, pool, versionStr)
	} else {
		localpath = fmt.Sprintf("%s/%s", u.uploadPath, handler.Filename)
	}
	file := utils.File{Path: localpath}

	file.Delete()

	localfd, err := os.OpenFile(localpath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		u.global.GetHttpTools().WriteError(w, -3)
		return
	}
	defer localfd.Close()

	io.Copy(localfd, clientfd)
	md5Str := r.FormValue("md5")
	md5Value := u.global.GetMD5().EncodeFile(localpath)
	if md5Value != md5Str {
		u.global.GetHttpTools().WriteError(w, -4)
		return
	}
	if md5Str == "" || versionStr == "" {
		return
	}
	u.global.GetHttpTools().WriteSuccess(w)
}

func (u *UploadFile) GetRequestParams() []string {
	var params []string
	params = append(params, "version")
	params = append(params, "pool")
	params = append(params, "md5")
	params = append(params, "file")
	return params
}

func (u *UploadFile) Execute(params map[string]interface{}) (map[string]interface{}, int) {
	taskid := params["taskid"]
	if taskid == nil {
		taskid := u.global.Random().RandomString(32)
		context := make(map[string]interface{})
		context["taskid"] = taskid
		context["filename"] = params["filename"]
		context["version"] = params["version"]
		context["pool"] = params["pool"]

		context["start"] = params["start"]
		context["filelength"] = params["filelength"]
		u.uploadFiles[taskid] = context

		localpath := fmt.Sprintf("%s/%s", u.uploadPath, params["filename"].(string))
		start := params["start"].(int)
		if start == 0 {
			f, err3 := os.Create(localpath)
			if err3 != nil {
				return context, 0
			}
			context["_file_handler"] = f
		} else {
			f, err3 := os.OpenFile(localpath, os.O_RDWR|os.O_CREATE, 0666)
			if err3 != nil {
				return context, 0
			}
			context["_file_handler"] = f
		}
		context["_keep_connection"] = true
		context["_closeListener"] = u
		context["_conn"] = params["_conn"]
		return context, 0
	} else {
		context := make(map[string]interface{})
		context["taskid"] = taskid
		fileinfo := u.uploadFiles[taskid.(string)].(map[string]interface{})
		f := fileinfo["_file_handler"].(*os.File)
		data := params["data"].([]byte)
		start := params["start"].(int)
		filelength := fileinfo["filelength"].(int)
		f.Seek(int64(start), 0)
		f.Write(data)
		context["start"] = params["start"]
		if start+len(data) == filelength {
			context["_keep_connection"] = false
		} else {
			context["_keep_connection"] = true
		}
		context["_closeListener"] = u
		return context, 0
	}

}
func (u *UploadFile) OnErrorHeader(conn *net.TCPConn) {

}
func (u *UploadFile) OnConnectClose(conn *net.TCPConn) {
	for key, value := range u.uploadFiles {
		c := value.(map[string]interface{})["_conn"].(*net.TCPConn)
		if c == conn {
			f := value.(map[string]interface{})["_file_handler"].(*os.File)
			f.Close()
			delete(u.uploadFiles, key)
			break
		}
	}

}

func (u *UploadFile) Decode(messageId, version, msgType int, data []byte) map[string]interface{} {
	byteTools := u.global.GetByteTools()
	ret := make(map[string]interface{})

	if messageId == u.global.GetMessages().UPLOAD_FILE() {
		pos := 0
		fnl := byteTools.BytesToShort(data[pos : pos+2])
		ret["filename"] = string(data[pos+2 : pos+2+fnl])
		pos = pos + 2 + fnl
		fvl := byteTools.BytesToShort(data[pos : pos+2])
		ret["version"] = string(data[pos+2 : pos+2+fvl])
		pos = pos + 2 + fvl
		gl := byteTools.BytesToShort(data[pos : pos+2])
		ret["pool"] = string(data[pos+2 : pos+2+gl])
		pos = pos + 2 + gl
		start := byteTools.BytesToInt(data[pos : pos+4])
		ret["start"] = int(start)
		filelength := byteTools.BytesToInt(data[pos+4 : pos+8])
		ret["filelength"] = int(filelength)
		ret["index"] = pos

	} else if messageId == u.global.GetMessages().UPLOAD_PACKAGE() {
		pos := 0
		tl := byteTools.BytesToShort(data[pos : pos+2])
		ret["taskid"] = string(data[pos+2 : pos+2+tl])
		pos = pos + 2 + tl
		start := byteTools.BytesToInt(data[pos : pos+4])
		length := byteTools.BytesToInt(data[pos+4 : pos+8])
		ret["start"] = start
		ret["length"] = length
		ret["data"] = data[pos+8 : pos+8+length]
	}
	return ret
}

func (u *UploadFile) Encode(messageId, version, msgType int, msg map[string]interface{}) []byte {
	if messageId == u.global.GetMessages().UPLOAD_FILE() {
		var buffer bytes.Buffer
		byteTools := u.global.GetByteTools()
		taskid := msg["taskid"].(string)
		buffer.Write(byteTools.ShortToBytes(len(taskid)))
		buffer.Write([]byte(taskid))

		start := msg["start"].(int)
		buffer.Write(byteTools.IntToBytes(int32(start)))
		return buffer.Bytes()
	} else {
		var buffer bytes.Buffer
		byteTools := u.global.GetByteTools()
		taskid := msg["taskid"].(string)
		buffer.Write(byteTools.ShortToBytes(len(taskid)))
		buffer.Write([]byte(taskid))

		start := msg["start"].(int)
		buffer.Write(byteTools.IntToBytes(int32(start)))
		u.global.GetLog().InfoA("Upload", "upload data at "+strconv.Itoa(start))
		return buffer.Bytes()
	}
}
