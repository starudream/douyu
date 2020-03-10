package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/go-sdk/logx"
	"github.com/go-sdk/utilx/dt"
	"github.com/go-sdk/utilx/json"

	"github.com/starudream/douyu"
)

var (
	R string // 房间号
	L int64  // 弹幕屏蔽等级
	M bool   // 显示弹幕
	G bool   // 显示礼物
	H bool   // 帮助
)

func init() {
	logx.SetLogger(logx.NewWithWriters(
		logx.NewConsoleWriter(logx.ConsoleWriterConfig{
			Level: logx.DebugLevel,
		}),
		logx.NewFileWriter(logx.FileWriterConfig{
			Level:    logx.DebugLevel,
			Type:     logx.TextFileWriter,
			NoColor:  true,
			Filename: "danmu.log",
		}),
	))

	flag.StringVar(&R, "r", "", "房间号")
	flag.Int64Var(&L, "l", 30, "弹幕屏蔽等级")
	flag.BoolVar(&M, "m", true, "弹幕")
	flag.BoolVar(&G, "g", false, "礼物")
	flag.BoolVar(&H, "h", false, "帮助")
	flag.Parse()

	if R == "" {
		flag.Usage()
		os.Exit(1)
	}
	if L < 0 {
		L = 0
	}
	if L > 120 {
		L = 120
	}
}

func main() {
	c := douyu.NewClient().SetRoomId(R)

	ch, _ := c.Start()
	defer c.Close()

	go func() {
		m := c.GetMessage()
		for {
			t := <-m
			msg := &douyu.Message{}
			_ = json.Unmarshal(json.MustMarshal(t), msg)
			switch msg.Type {
			case "chatmsg":
				if !M || msg.Level < dt.IntStr(L) {
					continue
				}
				format := "弹幕 %" + length(msg.NN, 30) + "s |%3d| %" + length(msg.BNN, 6) + "s |%3d|: %s"
				logx.Infof(format, msg.NN, msg.Level, msg.BNN, msg.BL, msg.Txt)
			case "dgb":
				if !G || msg.BG == 0 {
					continue
				}
				format := "礼物 %" + length(msg.NN, 30) + "s |%3d| %" + length(msg.BNN, 6) + "s |%3d|: %v %d 个，共 %d 个"
				logx.Infof(format, msg.NN, msg.Level, msg.BNN, msg.BL, gift(msg.GFid), msg.GFCnt, msg.Hits)
			}
		}
	}()

	<-ch
}

func length(s string, def int64) string {
	for _, v := range s {
		if !((v >= '0' && v <= '9') || (v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') || v == '_') {
			def--
		}
	}
	return strconv.FormatInt(def, 10)
}

func gift(id dt.IntStr) interface{} {
	if v, ok := douyu.GiftMap[int64(id)]; ok {
		return v
	}
	return id
}
