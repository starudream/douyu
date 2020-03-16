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
	Id    int64  `json:"id"`    // 礼物 pid
	Name  string `json:"name"`  // 名称
	Price int64  `json:"price"` // 价格，单位：分
}

var (
	GiftMap = map[int64]GiftData{}
	giftMu  = sync.Mutex{}
)

func GetGift(pid int64) GiftData {
	if v, ok := GiftMap[pid]; ok {
		return v
	}
	i := strconv.FormatInt(pid, 10)
	g := GiftData{Name: i}
	resp, err := http.Get(giftURL + i)
	if err != nil {
		logx.Errorf("gift: http get error: %v", err)
		return g
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("gift: read response error: %v", err)
		return g
	}
	giftResp := &GiftResp{}
	err = json.Unmarshal(bs, giftResp)
	if err != nil {
		logx.Errorf("gift: json unmarshal response error: %v", err)
		return g
	}
	if giftResp.Error != 0 {
		return g
	}
	giftMu.Lock()
	defer giftMu.Unlock()
	GiftMap[pid] = giftResp.Data
	return giftResp.Data
}
