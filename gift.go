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
	GiftList []Gift `json:"giftList"`
	Name     string `json:"name"`
}

type Gift struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

var (
	GiftMap1 = map[int64]string{}
	GiftMap2 = map[int64]string{}
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
		GiftMap1[gift.Id] = gift.Name
	}
}

func GetGift(id, pid int64) string {
	if v, ok := GiftMap1[id]; ok {
		return v
	}
	if v, ok := GiftMap2[pid]; ok {
		return v
	}
	j := strconv.FormatInt(pid, 10)
	giftResp := &GiftResp{}
	err := httpGet(giftURL2+j, giftResp)
	if err != nil {
		logx.Errorf("gift: http get error: %v", err)
		return j
	}
	if giftResp.Error != 0 {
		return j
	}
	if giftResp.Data.Name == "" {
		return j
	}
	giftMu.Lock()
	defer giftMu.Unlock()
	GiftMap2[pid] = giftResp.Data.Name
	return giftResp.Data.Name
}
