package configuration

type Configuration interface {
	Get(path string) string
	Unmarshal(path string, obj interface{})
}
