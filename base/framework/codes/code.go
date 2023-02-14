package codes

import "github.com/jinvei/microservice/base/api/proto/v1/dto"

type Code int

var codeMap map[Code]string

func init() {
	codeMap = make(map[Code]string)
}

// should call in init func
func Register(c map[Code]string) {
	for k, v := range c {
		codeMap[k] = v
	}
}

func (c Code) ToStatus() dto.Status {
	return dto.Status{
		Code: int64(c),
		Msg:  codeMap[c],
	}
}

func (c Code) Str() string {
	return codeMap[c]
}
