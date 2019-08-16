package plugins

import "common"
import "net/http"
import "io"
import "os"
import "utils"
import "fmt"

type UploadFile struct {
	authSecret string
	uploadPath string
	global     *common.Global
}

func (u *UploadFile) GetName() string {
	return "UploadFile"
}

func (u *UploadFile) InitPlugin(g *common.Global) {
	u.global = g
	g.GetLog().InfoA("UploadFile", "InitPlugin")
	g.GetRouter().HandleFunc("/upload", u.uploadHandler)
	u.authSecret = g.GetConfig().GetProperty("auth")
	u.uploadPath = g.GetConfig().GetProperty("uploadPath")
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
