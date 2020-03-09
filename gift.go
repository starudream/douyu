package douyu

import (
	"io/ioutil"
	"net/http"

	"github.com/go-sdk/logx"
	"github.com/go-sdk/utilx/json"
)

const (
	giftURL = "https://gift.douyucdn.cn/api/gift/v2/web/list?rid="
)

type GiftResp struct {
	Data  GiftData `json:"data"`
	Error int64    `json:"error"`
}

type GiftData struct {
	Tid      string `json:"tid"`
	GiftList []Gift `json:"giftList"`
}

type Gift struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

var (
	GiftMap = map[int64]string{}
)

func InitGift(rid string) {
	resp, err := http.Get(giftURL + rid)
	if err != nil {
		logx.Errorf("gift: http get error: %v", err)
		return
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("gift: read response error: %v", err)
		return
	}
	giftResp := &GiftResp{}
	err = json.Unmarshal(bs, giftResp)
	if err != nil {
		logx.Errorf("gift: json unmarshal response error: %v", err)
		return
	}
	if giftResp.Error != 0 {
		return
	}
	for _, gift := range giftResp.Data.GiftList {
		GiftMap[gift.Id] = gift.Name
	}
}
