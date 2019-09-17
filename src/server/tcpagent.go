package server

import (
	"bytes"
	"common"
	"net"
)

type TcpAgent struct {
	me       string
	global   *common.Global
	handlers map[int]common.RequestHandler
	msg2req  map[string]int
}

func (s *TcpAgent) InitServer(g *common.Global) {
	s.me = "TcpAgent"
	s.global = g
	s.handlers = make(map[int]common.RequestHandler)
	s.initMessagesMap()

}
func (s *TcpAgent) Type() string {
	return "TCP"
}

func (s *TcpAgent) initMessagesMap() {
	s.msg2req = map[string]int{}
	messages := s.global.GetMessages()
	s.msg2req[messages.LINK_FILE_TEXT()] = messages.LINK_FILE()
	s.msg2req[messages.JSTAT_TEXT()] = messages.JSTAT()
	s.msg2req[messages.UPLOAD_FILE_TEXT()] = messages.UPLOAD_FILE()
	s.msg2req[messages.UPLOAD_PACKAGE_TEXT()] = messages.UPLOAD_PACKAGE()
	s.msg2req[messages.DOWNLOAD_FILE_TEXT()] = messages.DOWNLOAD_FILE()
	s.msg2req[messages.DOWNLOAD_PACKAGE_TEXT()] = messages.DOWNLOAD_PACKAGE()
	s.msg2req[messages.MYSQL_CLIENT_TEXT()] = messages.MYSQL_CLIENT()
}

func (s *TcpAgent) StartServer() {
	var tcpAddr *net.TCPAddr
	port := s.global.GetConfig().GetProperty("port")
	//通过ResolveTCPAddr实例一个具体的tcp断点
	tcpAddr, _ = net.ResolveTCPAddr("tcp", ":"+port)
	//打开一个tcp断点监听
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)

	defer tcpListener.Close()
	s.global.GetLog().Info("Start Tcp Agent at [" + port + "]")
	//循环接收客户端的连接，创建一个协程具体去处理连接
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			continue
		}
		go s.handleConnection(tcpConn)
	}
}

func (s *TcpAgent) StopServer() {

}

func (s *TcpAgent) handleConnection(conn *net.TCPConn) {
	var closeListener common.CloseListener
	bytetools := s.global.GetByteTools()
	header := make([]byte, 14)
	for {
		n, _ := conn.Read(header)
		if n < 14 || header[0] != 33 || header[1] != 65 || header[2] != 82 || header[3] != 81 { // Header
			break
		}
		//requestSeq := bytetools.BytesToInt(header[4:6]) // requestSeq
		pos := 6
		messageId := bytetools.BytesToShort(header, &pos)
		version := bytetools.BytesToShort(header, &pos)
		length := bytetools.BytesToInt(header, &pos)
		if length > 0 {
			tempBuff := make([]byte, length)
			readBytes := 0
			readError := false
			for {
				rl, _ := conn.Read(tempBuff[readBytes:length])
				readBytes += rl
				if rl == 0 {
					readError = true
					break
				}
				if readBytes == length {
					break
				}
			}
			if readError {
				break
			}
			result, err := s.onHandleMessage(conn, messageId, version, tempBuff[0:len(tempBuff)])
			encoder := s.global.GetMessageCoder().GetEncoder(messageId, version)
			retBytes := encoder.Encode(messageId, version, 2, result)
			var resultBuf bytes.Buffer
			if err != 0 {
				resultBuf.Write(bytetools.BoolToBytes(false)) // Success : false
				resultBuf.Write(bytetools.ShortToBytes(err))  // Error :
			} else {
				resultBuf.Write(bytetools.BoolToBytes(true))  // Success : true
				retLen := int32(len(retBytes))                // Data :
				resultBuf.Write(bytetools.IntToBytes(retLen)) //
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
			keep := result["_keep_connection"]
			if keep == nil || !(keep.(bool)) {
				break
			}
			if result["_closeListener"] != nil {
				closeListener = result["_closeListener"].(common.CloseListener)
			}
		} else {
			break
		}
	}
	if closeListener != nil {
		closeListener.OnConnectClose(conn)
	}
	conn.Close()
}

func (s *TcpAgent) RegisterHandler(req string, h common.RequestHandler) {
	s.global.GetLog().InfoA("Global", "Register Handler ["+h.GetName()+"] for URI ["+req+"]")
	msgId := s.msg2req[req]
	s.handlers[msgId] = h
	//s.router.HandleFunc("/"+req, s.handlerFunc)
}

func (s *TcpAgent) onHandleMessage(conn *net.TCPConn, messageId, version int, datas []byte) (map[string]interface{}, int) {
	decoder := s.global.GetMessageCoder().GetDecoder(messageId, version)
	if decoder == nil {
		return nil, -1
	}
	params := decoder.Decode(messageId, version, 1, datas)
	params["_conn"] = conn
	params["_messageId"] = messageId
	if params != nil {
		handler := s.handlers[messageId]
		result, err := handler.Execute(params)
		return result, err
	}
	return nil, -1
}

func (s *TcpAgent) ReadData(context map[string]interface{}, start, length int) []byte {
	conn := context["_conn"].(*net.TCPConn)
	tempBuff := make([]byte, length)
	conn.Read(tempBuff)
	return tempBuff
}
