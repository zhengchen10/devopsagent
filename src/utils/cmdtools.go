package utils

import "fmt"
import "os/exec"

type CmdTools struct {
	log *Log
}

func (c *CmdTools) SetLogger(log *Log) {
	c.log = log
}

func (c *CmdTools) Execute(cmd, param string) (ret string, err error) {
	c.log.DebugA("CmdTools", "Execute commond ["+cmd+"] params ["+param+"]")
	cc := exec.Command("sh", "-c", cmd+" "+param)
	out, err := cc.Output()

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	ret = string(out)
	c.log.DebugA("CmdTools", "Execute commond ["+cmd+"] output \r\n"+ret)
	return ret, nil
}
