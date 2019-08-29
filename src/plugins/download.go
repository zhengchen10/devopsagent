package plugins

import (
	"archive/zip"
	"bytes"
	"common"
	"fmt"
	"net"
	"net/http"
	"server"
)
import "io"
import "os"
import "strings"
import "strconv"
import "time"

type DownloadFile struct {
	me            string
	global        *common.Global
	downloadFiles map[string]interface{}
	downloadPath  string
}

const fs_maxbufsize = 4096 // 4096 bits = default page size on OSX

func (u *DownloadFile) GetName() string {
	return u.me
}

func (u *DownloadFile) InitPlugin(g *common.Global) {
	u.me = "DownloadFile"
	u.global = g
	g.GetLog().InfoA(u.me, "InitPlugin")
	//g.RegisterHandler("download", u)
	agent := u.global.GetConfig().GetProperty("agent")
	if agent == "TCP" {
		g.RegisterHandler(u.global.GetMessages().DOWNLOAD_FILE_TEXT(), u)
		g.RegisterHandler(u.global.GetMessages().DOWNLOAD_PACKAGE_TEXT(), u)
		g.GetMessageCoder().RegisterDecoder(g.GetMessages().DOWNLOAD_FILE(), 1, u)
		g.GetMessageCoder().RegisterEncoder(g.GetMessages().DOWNLOAD_FILE(), 1, u)
		g.GetMessageCoder().RegisterDecoder(g.GetMessages().DOWNLOAD_PACKAGE(), 1, u)
		g.GetMessageCoder().RegisterEncoder(g.GetMessages().DOWNLOAD_PACKAGE(), 1, u)
	} else {
		u.global.GetAppServer().(*server.HttpAgent).GetRouter().HandleFunc("/download", u.downloadHandler)
	}
	u.downloadFiles = make(map[string]interface{})
	u.downloadPath = g.GetConfig().GetProperty("default.DownloadPath")
}

func (u *DownloadFile) StartPlugin() {

}

func (u *DownloadFile) StopPlugin() {

}

func (u *DownloadFile) min(x int64, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func (u *DownloadFile) downloadHandler(w http.ResponseWriter, r *http.Request) {
	file := r.FormValue("name")
	path := r.FormValue("path")
	uses_gzip := r.FormValue("compress")

	filepath := path + file
	// Opening the file handle
	f, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "404 Not Found : Error while opening the file.", 404)
		return
	}

	defer f.Close()

	// Checking if the opened handle is really a file
	statinfo, err := f.Stat()
	if err != nil {
		http.Error(w, "500 Internal Error : stat() failure.", 500)
		return
	}

	if statinfo.IsDir() { // If it's a directory, open it !
		//handleDirectory(f, w, req)
		return
	}

	if (statinfo.Mode() &^ 07777) == os.ModeSocket { // If it's a socket, forbid it !
		http.Error(w, "403 Forbidden : you can't access this resource.", 403)
		return
	}

	// Manages If-Modified-Since and add Last-Modified (taken from Golang code)
	if t, err := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since")); err == nil && statinfo.ModTime().Unix() <= t.Unix() {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("Last-Modified", statinfo.ModTime().Format(http.TimeFormat))

	// Manage Content-Range (TODO: Manage end byte and multiple Content-Range)
	/*if r.Header.Get("Range") != "" {
		start_byte := u.parseRange(r.Header.Get("Range"))

		if start_byte < statinfo.Size() {
			f.Seek(start_byte, 0)
		} else {
			start_byte = 0
		}

		w.Header().Set("Content-Range",
			fmt.Sprintf("bytes %d-%d/%d", start_byte, statinfo.Size()-1, statinfo.Size()))
	}*/

	// Manage gzip/zlib compression
	output_writer := w.(io.Writer)

	//tmpfile, err := ioutil.TempFile("", "tmpfile")
	buf := make([]byte, fs_maxbufsize)
	if uses_gzip == "zip" {
		zipbuf := new(bytes.Buffer)
		zw := zip.NewWriter(zipbuf)
		fw, err := zw.Create(file)
		n := 0
		for err == nil {
			n, err = f.Read(buf)
			fw.Write(buf[0:n])
		}
		zw.Close()
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment;filename="+file+".zip")
		w.Header().Set("Content-Length", strconv.Itoa(zipbuf.Len()))
		output_writer.Write(zipbuf.Bytes())

	} else {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment;filename="+file)
		w.Header().Set("Content-Length", strconv.FormatInt(statinfo.Size(), 10))

		n := 0
		for err == nil {
			n, err = f.Read(buf)
			output_writer.Write(buf[0:n])
		}
	}
}

func (u *DownloadFile) parseCSV(data string) []string {
	splitted := strings.SplitN(data, ",", -1)
	data_tmp := make([]string, len(splitted))

	for i, val := range splitted {
		data_tmp[i] = strings.TrimSpace(val)
	}

	return data_tmp
}

func (u *DownloadFile) parseRange(data string) int64 {
	stop := (int64)(0)
	part := 0
	for i := 0; i < len(data) && part < 2; i = i + 1 {
		if part == 0 { // part = 0 <=> equal isn't met.
			if data[i] == '=' {
				part = 1
			}
			continue
		}
		if part == 1 { // part = 1 <=> we've met the equal, parse beginning
			if data[i] == ',' || data[i] == '-' {
				part = 2 // part = 2 <=> OK DUDE.
			} else {
				if 48 <= data[i] && data[i] <= 57 { // If it's a digit ...
					// ... convert the char to integer and add it!
					stop = (stop * 10) + (((int64)(data[i])) - 48)
				} else {
					part = 2 // Parsing error! No error needed : 0 = from start.
				}
			}
		}
	}

	return stop
}

func (u *DownloadFile) GetRequestParams() []string {
	var params []string
	return params
}

func (u *DownloadFile) Execute(params map[string]interface{}) (map[string]interface{}, int) {
	result := make(map[string]interface{})
	taskid := params["taskid"]
	if taskid == nil {
		taskid := u.global.Random().RandomString(32)
		filename := params["filename"].(string)
		group := params["group"].(string)
		workpath := u.global.GetConfig().GetProperty(group + ".DownloadPath")
		localpath := fmt.Sprintf("%s/%s", workpath, filename)
		f, err := os.Open(localpath)
		if err != nil {
			return nil, -2
		}
		filelength, err := f.Seek(0, os.SEEK_END)
		result["taskid"] = taskid
		result["filelength"] = int(filelength)
		result["start"] = params["start"].(int)
		result["filename"] = params["filename"].(string)
		result["_file_handler"] = f
		result["_keep_connection"] = true
		result["_closeListener"] = u
		result["_conn"] = params["_conn"]
		u.downloadFiles[taskid] = result
	} else {
		result["taskid"] = taskid
		result["_closeListener"] = u
		fileinfo := u.downloadFiles[taskid.(string)].(map[string]interface{})
		start := params["start"].(int)
		length := params["length"].(int)
		dataBuff := make([]byte, length)
		f := fileinfo["_file_handler"].(*os.File)

		readed := 0
		success := false
		for {
			n, err := f.ReadAt(dataBuff[readed:length], int64(start))
			if n == 0 {
				break
			}
			if err != nil {
				break
			}
			readed += n
			if readed == length {
				success = true
				break
			}
		}
		result["start"] = start
		result["length"] = length
		result["data"] = dataBuff
		if success {
			result["_keep_connection"] = true
		} else {
			result["_keep_connection"] = false
		}
	}
	return result, 0
}

func (u *DownloadFile) Decode(messageId, version, msgType int, data []byte) map[string]interface{} {
	byteTools := u.global.GetByteTools()
	ret := make(map[string]interface{})

	if messageId == u.global.GetMessages().DOWNLOAD_FILE() {
		pos := 0
		ret["filename"] = byteTools.ReadString(data, &pos)
		ret["group"] = byteTools.ReadString(data, &pos)
		ret["needCompress"] = byteTools.BytesToBool(data, &pos)
		ret["start"] = byteTools.BytesToInt(data, &pos)
	} else if messageId == u.global.GetMessages().DOWNLOAD_PACKAGE() {
		pos := 0
		ret["taskid"] = byteTools.ReadString(data, &pos)
		start := byteTools.BytesToInt(data, &pos)
		length := byteTools.BytesToInt(data, &pos)
		ret["start"] = start
		ret["length"] = length
	}
	return ret
}

func (u *DownloadFile) Encode(messageId, version, msgType int, msg map[string]interface{}) []byte {
	var buffer bytes.Buffer
	byteTools := u.global.GetByteTools()
	if messageId == u.global.GetMessages().DOWNLOAD_FILE() {
		byteTools.WriteString(&buffer, msg["taskid"].(string))
		byteTools.WriteString(&buffer, msg["filename"].(string))
		buffer.Write(byteTools.IntToBytes(int32(msg["filelength"].(int))))
		buffer.Write(byteTools.IntToBytes(int32(msg["start"].(int))))
	} else if messageId == u.global.GetMessages().DOWNLOAD_PACKAGE() {
		buffer.Write(byteTools.IntToBytes(int32(msg["start"].(int))))
		buffer.Write(byteTools.IntToBytes(int32(msg["length"].(int))))
		buffer.Write(msg["data"].([]byte))
	}
	return buffer.Bytes()
}

func (u *DownloadFile) OnErrorHeader(conn *net.TCPConn) {

}
func (u *DownloadFile) OnConnectClose(conn *net.TCPConn) {
	for key, value := range u.downloadFiles {
		c := value.(map[string]interface{})["_conn"].(*net.TCPConn)
		if c == conn {
			f := value.(map[string]interface{})["_file_handler"].(*os.File)
			f.Close()
			delete(u.downloadFiles, key)
			break
		}
	}
}
