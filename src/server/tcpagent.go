package server

import (
	"bytes"
	"common"
	"messages"
	"net"
)

type TcpAgent struct {
	me       string
	global   *common.Global
	handlers map[int]common.RequestHandler
	msg2req  map[string]int

	MessageConst *messages.MessageConst
	//decoder  *coder.MessageDecoder
	//encoder  *coder.MessageEncoder
}

func (s *TcpAgent) InitServer(g *common.Global) {
	s.me = "TcpAgent"
	s.global = g
	s.handlers = make(map[int]common.RequestHandler)
	s.MessageConst = new(messages.MessageConst)
	//s.decoder = new (coder.MessageDecoder)
	//s.encoder = new (coder.MessageEncoder)
	//s.decoder.Init(g)
	//s.encoder.Init(g)
	s.initMessagesMap()

}
func (s *TcpAgent) Type() string {
	return "TCP"
}

func (s *TcpAgent) initMessagesMap() {
	s.msg2req = map[string]int{}
	s.msg2req["healthCheck"] = 1
	s.msg2req["ln"] = 10000

	//s.registerDecoder(s.MessageConst.LinkFileRequest(),1,new(messages.LinkFileRequest))
	//s.registerEncoder(s.MessageConst.LinkFileResponse(),1,new(messages.LinkFileResponse))
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
		go s.handleConnection(tcpConn)
	}
}

func (s *TcpAgent) StopServer() {

}

func (s *TcpAgent) handleConnection(conn *net.TCPConn) {
	//tcp连接的地址
	//ipStr := conn.RemoteAddr().String()
	header := make([]byte, 14)
	n, _ := conn.Read(header)
	if n < 14 || header[0] != 33 || header[1] != 65 || header[2] != 82 || header[3] != 81 { // Header
		conn.Close()
		return
	}
	bytetools := s.global.GetByteTools()
	//requestSeq := bytetools.BytesToInt(header[4:6]) // requestSeq
	messageId := bytetools.BytesToShort(header[6:8])
	version := bytetools.BytesToShort(header[8:10])
	length := bytetools.BytesToInt(header[10:14])
	if length > 0 {
		tempBuff := make([]byte, length)
		conn.Read(tempBuff)
		result, err := s.onHandleMessage(messageId, version, tempBuff[0:len(tempBuff)])
		/*ret := make(map[string]interface{})
		if err != 0 {
			ret["success"]= false
			ret["error"] = err
		} else {
			ret["success"]= true
			ret["data"]=result
		}*/
		encoder := s.global.GetMessageCoder().GetEncoder(messageId, version)
		retBytes := encoder.Encode(messageId, version, 2, result)
		var resultBuf bytes.Buffer
		if err != 0 {
			resultBuf.Write(bytetools.BoolToBytes(false))
			resultBuf.Write(bytetools.ShortToBytes(err))
		} else {
			resultBuf.Write(bytetools.BoolToBytes(true))
			retLen := int32(len(retBytes))
			resultBuf.Write(bytetools.IntToBytes(retLen))
			if retLen > 0 {
				resultBuf.Write(retBytes)
			}
		}

		var msgBuf bytes.Buffer
		msgBuf.Write([]byte("!ARP"))
		msgBuf.Write(header[4:10]) // requestSeq,messageId,version
		resultBytes := resultBuf.Bytes()
		msgBuf.Write(bytetools.IntToBytes(int32(len(resultBytes))))
		msgBuf.Write(resultBytes)

		conn.Write(msgBuf.Bytes())
	}
	conn.Close()
}

func (s *TcpAgent) RegisterHandler(req string, h common.RequestHandler) {
	s.global.GetLog().InfoA("Global", "Register Handler ["+h.GetName()+"] for URI ["+req+"]")
	msgId := s.msg2req[req]
	s.handlers[msgId] = h
	//s.router.HandleFunc("/"+req, s.handlerFunc)
}

func (s *TcpAgent) onHandleMessage(messageId, version int, datas []byte) (map[string]interface{}, int) {
	decoder := s.global.GetMessageCoder().GetDecoder(messageId, version)
	if decoder == nil {
		return nil, -1
	}
	params := decoder.Decode(messageId, version, 1, datas)
	if params != nil {
		handler := s.handlers[messageId]
		result, err := handler.Execute(params)
		return result, err
	}
	return nil, -1
}
