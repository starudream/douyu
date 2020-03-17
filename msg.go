package douyu

import (
	"github.com/go-sdk/utilx/dt"
)

type Message struct {
	Type  string    `json:"type,omitempty"`  // 类型
	Rid   string    `json:"rid,omitempty"`   // 房间 id
	NN    string    `json:"nn,omitempty"`    // 发送者昵称
	NL    dt.IntStr `json:"nl,omitempty"`    // 贵族等级
	Level dt.IntStr `json:"level,omitempty"` // 用户等级
	BNN   string    `json:"bnn,omitempty"`   // 勋章昵称
	BL    dt.IntStr `json:"bl,omitempty"`    // 勋章等级
	Txt   string    `json:"txt,omitempty"`   // 弹幕内容
	BG    dt.IntStr `json:"bg,omitempty"`    // 是否是大礼物
	GFid  dt.IntStr `json:"gfid,omitempty"`  // 礼物 id
	Pid   dt.IntStr `json:"pid,omitempty"`   // 礼物 id
	GFCnt dt.IntStr `json:"gfcnt,omitempty"` // 礼物个数
	Hits  dt.IntStr `json:"hits,omitempty"`  // 礼物连击数
	DRid  dt.IntStr `json:"drid,omitempty"`  // 爵位房间 id
	Unk   string    `json:"unk,omitempty"`   // 爵位姓名
}
