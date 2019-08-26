package common

type RequestHandler interface {
	GetName() string
	GetRequestParams() []string
	Execute(params map[string]interface{}) (map[string]interface{}, int)
}
