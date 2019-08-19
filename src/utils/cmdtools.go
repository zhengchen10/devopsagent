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

func (c *CmdTools) Execute(command, param string, useSH bool) (ret string, err error) {
	c.log.DebugA("CmdTools", "Execute command ["+command+"] params ["+param+"]")
	var cmd *exec.Cmd
	if useSH {
		cmd = exec.Command("sh", "-c", command+" "+param)
	} else {
		cmd = exec.Command(command, param)
	}
	out, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	ret = string(out)
	c.log.DebugA("CmdTools", "Execute command ["+command+"] output \r\n"+ret)
	return ret, nil
}

/**
 * 执行命令并逐行回调
 */
func (c *CmdTools) ExecuteWithCallback(command, param string, context interface{}, useSH bool,
	callback func(cmd *exec.Cmd, line string, context interface{}) bool) bool {
	c.log.DebugA("CmdTools", "Execute command ["+command+"] params ["+param+"]")
	var cmd *exec.Cmd
	if useSH {
		cmd = exec.Command("sh", "-c", command+" "+param)
	} else {
		cmd = exec.Command(command, param)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}
	reader := bufio.NewReader(stdout)
	cmd.Start()

	//实时循环读取输出流中的一行内容
	for {
		line, isprefix, err2 := reader.ReadLine()
		if err2 != nil || io.EOF == err2 {
			//fmt.Println(err2)
			break
		}
		if !isprefix {
			cancel := callback(cmd, string(line), context)
			if cancel {
				break
			}
		}
		//fmt.Println("output : " + string(line))
	}
	cmd.Wait()
	return true
}
