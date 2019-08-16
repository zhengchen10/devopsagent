package utils

import "io"
import "crypto/md5"
import "encoding/hex"
import "os"

type MD5 struct {
}

func (m *MD5) Encode(datas []byte) string {
	h := md5.New()
	h.Write(datas)
	return hex.EncodeToString(h.Sum(nil))
}

func (m *MD5) EncodeFile(path string) string {
	file, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return ""
	}
	h := md5.New()
	// 按字节读取
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if 0 == n {
			break
		}
		if n == 1024 {
			h.Write(buf)
		} else {
			sbuf := make([]byte, n)
			copy(sbuf, buf)
			h.Write(sbuf)
		}
	}
	return hex.EncodeToString(h.Sum(nil))
}

func (m *MD5) CheckMD5(datas []byte, md5String string) bool {
	newMd5 := m.Encode(datas)
	return md5String == newMd5
}
