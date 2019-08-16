package utils

import "fmt"
import "os/exec"
import "bufio"
import "io"

type CmdTools struct {
	log *Log
}

func (c *CmdTools) SetLogger(log *Log) {
	c.log = log
}

func (c *CmdTools) Execute(commond, param string) (ret string, err error) {
	c.log.DebugA("CmdTools", "Execute commond ["+commond+"] params ["+param+"]")
	cmd := exec.Command("sh", "-c", commond+" "+param)
	out, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	ret = string(out)
	c.log.DebugA("CmdTools", "Execute commond ["+commond+"] output \r\n"+ret)
	return ret, nil
}

/**
 * 执行命令并逐行回调
 */
func (c *CmdTools) ExecuteWithCallback(commond, param string, context interface{},
	callback func(cmd *exec.Cmd, line string, context interface{}) bool) bool {
	c.log.DebugA("CmdTools", "Execute commond ["+commond+"] params ["+param+"]")
	cmd := exec.Command("sh", "-c", commond+" "+param)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	for {
		line, isprefix, err2 := reader.ReadLine()
		if err2 != nil || io.EOF == err2 {
			break
		}
		if !isprefix {
			cancel := callback(cmd, string(line), context)
			if cancel {
				break
			}
		}
	}
	cmd.Wait()
	return true
}
