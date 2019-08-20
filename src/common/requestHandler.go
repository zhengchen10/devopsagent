package common

type RequestHandler interface {
	GetName() string
	GetRequestParams() []string
	Execute(params map[string]string) (map[string]interface{}, int)
}
