package plugins

import (
	"archive/zip"
	"bytes"
	"common"
	"net/http"
	"server"
)
import "io"
import "os"
import "strings"
import "strconv"
import "time"

type DownloadFile struct {
	me     string
	global *common.Global
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
		//global.InitAgent(new (server.TcpAgent))
	} else {
		u.global.GetAppServer().(*server.HttpAgent).GetRouter().HandleFunc("/download", u.downloadHandler)
	}
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
