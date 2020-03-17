package douyu

import (
	"strconv"
	"sync"

	"github.com/go-sdk/logx"
)

const (
	giftURL1 = "https://gift.douyucdn.cn/api/gift/v2/web/list?rid="
	giftURL2 = "https://gift.douyucdn.cn/api/prop/v1/web/single?pid="
)

type GiftResp struct {
	Data  GiftData `json:"data"`
	Error int64    `json:"error"`
}

type GiftData struct {
	GiftList []GiftInfo `json:"giftList"`
	Name     string     `json:"name"`
	Price    int64      `json:"price"`
}

type GiftInfo struct {
	Id        int64         `json:"id"`
	Name      string        `json:"name"`
	PriceInfo GiftPriceInfo `json:"priceInfo"`
}

type GiftPriceInfo struct {
	Price int64 `json:"price"`
}

type Gift struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

var (
	GiftMap1 = map[int64]Gift{}
	GiftMap2 = map[int64]Gift{}
	giftMu   = sync.Mutex{}
)

func InitGift(rid string) {
	giftResp := &GiftResp{}
	err := httpGet(giftURL1+rid, giftResp)
	if err != nil {
		logx.Errorf("gift: http get error: %v", err)
		return
	}
	if giftResp.Error != 0 {
		return
	}
	for _, gift := range giftResp.Data.GiftList {
		GiftMap1[gift.Id] = Gift{Name: gift.Name, Price: gift.PriceInfo.Price}
	}
}

func GetGift(id, pid int64) Gift {
	if v, ok := GiftMap1[id]; ok {
		return v
	}
	if v, ok := GiftMap2[pid]; ok {
		return v
	}
	g := Gift{Name: strconv.FormatInt(pid, 10)}
	giftResp := &GiftResp{}
	err := httpGet(giftURL2+g.Name, giftResp)
	if err != nil {
		logx.Errorf("gift: http get error: %v", err)
		return g
	}
	if giftResp.Error != 0 {
		return g
	}
	if giftResp.Data.Name == "" {
		return g
	}
	giftMu.Lock()
	defer giftMu.Unlock()
	g = Gift{Name: giftResp.Data.Name, Price: giftResp.Data.Price}
	GiftMap2[pid] = g
	return g
}
