package common

type AppPlugin interface {
	GetName() string
	InitPlugin(g *Global)
	StartPlugin()
	StopPlugin()
}
