package server

import (
	"bufio"
	"common"
	"fmt"
	"io"
	"net"
	"time"
)

type TcpAgent struct {
	me       string
	global   *common.Global
	handlers map[string]common.RequestHandler
}

func (s *TcpAgent) InitServer(g *common.Global) {
	s.me = "TcpAgent"
	s.global = g
	s.handlers = make(map[string]common.RequestHandler)
}
func (s *TcpAgent) StartServer() {
	var tcpAddr *net.TCPAddr
	port := s.global.GetConfig().GetProperty("port")
	//通过ResolveTCPAddr实例一个具体的tcp断点
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:"+port)
	//打开一个tcp断点监听
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)

	defer tcpListener.Close()
	s.global.GetLog().Info("Start Tcp Agent at [" + port + "]")
	//循环接收客户端的连接，创建一个协程具体去处理连接
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			//fmt.Println(err)
			continue
		}
		//fmt.Println("A client connected :" +tcpConn.RemoteAddr().String())
		go s.tcpPipe(tcpConn)
	}
}

func (s *TcpAgent) StopServer() {

}

func (s *TcpAgent) tcpPipe(conn *net.TCPConn) {
	//tcp连接的地址
	ipStr := conn.RemoteAddr().String()
	defer func() {
		fmt.Println(" Disconnected : " + ipStr)
		conn.Close()
	}()

	//获取一个连接的reader读取流
	reader := bufio.NewReader(conn)
	i := 0
	//接收并返回消息
	for {
		message, err := reader.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}
		fmt.Println(string(message))
		time.Sleep(time.Second * 3)
		msg := time.Now().String() + conn.RemoteAddr().String() + " Server Say hello! \n"
		b := []byte(msg)
		conn.Write(b)
		i++
		if i > 10 {
			break
		}
	}
}

func (s *TcpAgent) RegisterHandler(req string, h common.RequestHandler) {
	s.global.GetLog().InfoA("Global", "Register Handler ["+h.GetName()+"] for URI ["+req+"]")
	s.handlers[req] = h
	//s.router.HandleFunc("/"+req, s.handlerFunc)
}
