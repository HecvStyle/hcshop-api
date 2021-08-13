package response

import (
	"fmt"
	"time"
)

type JsonTime time.Time

// MarshalJson 这个是为了对time类型做json处理时候，保持指定的格式
func (j JsonTime) MarshalJson() ([]byte, error) {
	var tmp = fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02"))
	return []byte(tmp), nil
}

type UserResponse struct {
	Id       int32  `json:"id"`
	Nickname string `json:"nickname"`
	Mobile   string `json:"mobile"`
	//Birthday string `json:"birthday"`
	Gender   string   `json:"gender"`
	Birthday JsonTime `json:"birthday"`
}
