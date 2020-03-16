package douyu

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-sdk/logx"
	"github.com/go-sdk/utilx/json"
)

const (
	giftURL = "https://gift.douyucdn.cn/api/prop/v1/web/single?pid="
)

type GiftResp struct {
	Data  GiftData `json:"data"`
	Error int64    `json:"error"`
}

type GiftData struct {
	Id        int64  `json:"id"`        // 礼物 pid
	Name      string `json:"name"`      // 名称
	Price     int64  `json:"price"`     // 价格，单位：分
	PriceType int64  `json:"priceType"` // 金钱类型
}

var (
	GiftMap     = map[int64]GiftData{}
	giftDefault = GiftData{}
	giftMu      = sync.Mutex{}
)

func GetGift(pid int64) GiftData {
	if v, ok := GiftMap[pid]; ok {
		return v
	}
	resp, err := http.Get(giftURL + strconv.FormatInt(pid, 10))
	if err != nil {
		logx.Errorf("gift: http get error: %v", err)
		return giftDefault
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("gift: read response error: %v", err)
		return giftDefault
	}
	giftResp := &GiftResp{}
	err = json.Unmarshal(bs, giftResp)
	if err != nil {
		logx.Errorf("gift: json unmarshal response error: %v", err)
		return giftDefault
	}
	if giftResp.Error != 0 {
		return giftDefault
	}
	giftMu.Lock()
	defer giftMu.Unlock()
	GiftMap[pid] = giftResp.Data
	return giftResp.Data
}
