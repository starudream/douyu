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
			Level: logx.InfoLevel,
		}),
		logx.NewFileWriter(logx.FileWriterConfig{
			Level:    logx.DebugLevel,
			Type:     logx.TextFileWriter,
			Filename: "danmu.log",
		}),
	))

	flag.StringVar(&R, "r", "", "房间号")
	flag.Int64Var(&L, "l", 30, "弹幕屏蔽等级")
	flag.BoolVar(&M, "m", false, "弹幕")
	flag.BoolVar(&G, "g", false, "礼物")
	flag.BoolVar(&H, "h", false, "帮助")
	flag.Parse()

	if R == "" || H {
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
	logx.Infof("房间号：%s，弹幕屏蔽等级：%d，是否开启弹幕：%t，是否开启礼物：%t", R, L, M, G)

	c := douyu.NewClient().SetRoomId(R)

	ch, _ := c.Start()
	defer c.Close()

	go func() {
		m := c.GetMessage()
		for {
			t := <-m
			msg := &douyu.Message{}
			bs := json.MustMarshal(t)
			logx.Debug(string(bs))
			_ = json.Unmarshal(bs, msg)
			nl := noble(msg.NL)
			switch msg.Type {
			case "chatmsg":
				if !M || msg.Level < dt.IntStr(L) {
					continue
				}
				format := "弹幕 %" + length(msg.NN, 30) + "s |%3d| | %" + length(nl, 4) + "s | %" + length(msg.BNN, 6) + "s |%3d|: %s"
				logx.Infof(format, msg.NN, msg.Level, nl, msg.BNN, msg.BL, msg.Txt)
			case "dgb":
				if !G || msg.BG == 0 {
					continue
				}
				g := douyu.GetGift(int64(msg.Pid))
				format := "礼物 %" + length(msg.NN, 30) + "s |%3d| | %" + length(nl, 4) + "s | %" + length(msg.BNN, 6) + "s |%3d|: %v %d 个，共 %d 个"
				logx.Infof(format, msg.NN, msg.Level, nl, msg.BNN, msg.BL, g.Name, msg.GFCnt, msg.Hits)
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

func noble(nl dt.IntStr) string {
	if v, ok := douyu.NobleMap[int64(nl)]; ok {
		return v
	}
	return ""
}
