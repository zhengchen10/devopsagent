package common

type AppServer interface {
	InitServer(g *Global)
	StartServer()
	StopServer()
	RegisterHandler(req string, h RequestHandler)
	Type() string
}
