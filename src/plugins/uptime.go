package plugins

import "common"
import "net/http"
import "io"
import "os"
import "utils"
import "fmt"

type UptimeInfo struct {
	global *common.Global
}

func (u *UptimeInfo) GetName() string {
	return "UptimeInfo"
}

func (u *UptimeInfo) InitPlugin(g *common.Global) {
	u.global = g
	g.GetLog().InfoA("UptimeInfo", "InitPlugin")
	g.GetRouter().HandleFunc("/uptime", u.uploadHandler)
}

func (u *UploadFile) StartPlugin() {

}

func (u *UploadFile) StopPlugin() {

}

func (u *UploadFile) uploadHandler(w http.ResponseWriter, r *http.Request) {
	u.global.GetHttpTools().WriteSuccess(w)
}
