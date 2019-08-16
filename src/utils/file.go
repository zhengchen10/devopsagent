package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type File struct {
	Path string
}

//使用os.OpenFile()相关函数打开文件对象，并使用文件对象的相关方法进行文件写入操作
//清空一次文件
func (f *File) CreateNew() bool {
	fileObj, err := os.OpenFile(f.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Failed to open the file", err.Error())
		return false
	}
	defer fileObj.Close()
	return true
}

func (f *File) AppendContent(content string) bool {
	fileObj, err := os.OpenFile(f.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to open the file", err.Error())
		return false
	}
	defer fileObj.Close()
	if _, err := io.WriteString(fileObj, content); err == nil {

	}
	return true
}

func (f *File) ReadContent() string {
	b, err := ioutil.ReadFile(f.Path)
	if err != nil {
		fmt.Print(err)
		return ""
	}
	fmt.Println(b)
	str := string(b)
	return str
}

func (f *File) Delete() bool {
	err := os.Remove(f.Path) //删除文件test.txt
	if err != nil {
		//输出错误详细信息
		fmt.Printf("%s", err)
		return false
	} else {
		//如果删除成功则输出 file remove OK!
		//fmt.Print("file remove OK!")
		return true
	}
}
