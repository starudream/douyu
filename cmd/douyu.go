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
	R  string // 房间号
	M  bool   // 显示弹幕
	ML int64  // 弹幕屏蔽等级
	G  bool   // 显示礼物
	U  bool   // 显示进入提醒
	UL int64  // 进入提醒屏蔽等级
	H  bool   // 帮助
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
	flag.BoolVar(&M, "m", false, "弹幕")
	flag.Int64Var(&ML, "ml", 30, "弹幕屏蔽等级")
	flag.BoolVar(&G, "g", false, "礼物")
	flag.BoolVar(&U, "u", false, "进入提醒")
	flag.Int64Var(&UL, "ul", 50, "进入提醒屏蔽等级")
	flag.BoolVar(&H, "h", false, "帮助")
	flag.Parse()

	if R == "" || H {
		flag.Usage()
		os.Exit(1)
	}
	if ML < 0 {
		ML = 0
	}
	if ML > 120 {
		ML = 120
	}
	if UL < 0 {
		UL = 0
	}
	if UL > 120 {
		UL = 120
	}
}

func main() {
	logx.Infof("房间号：%s，是否开启弹幕：%t（%d），是否开启礼物：%t，是否开启进入提醒：%t（%d）", R, M, ML, G, U, UL)

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
			nl := douyu.NobleMap[int64(msg.NL)]
			switch msg.Type {
			case "chatmsg":
				if !M || msg.Level < dt.IntStr(ML) {
					continue
				}
				format := "弹幕 %" + length(msg.NN, 30) + "s |%3d| %" + length(nl, 4) + "s | %" + length(msg.BNN, 6) + "s |%3d|: %s"
				logx.Infof(format, msg.NN, msg.Level, nl, msg.BNN, msg.BL, msg.Txt)
			case "dgb":
				if !G || msg.BG == 0 {
					continue
				}
				format := "礼物 %" + length(msg.NN, 30) + "s |%3d| %" + length(nl, 4) + "s | %" + length(msg.BNN, 6) + "s |%3d|: %v %d 个，共 %d 个"
				logx.Infof(format, msg.NN, msg.Level, nl, msg.BNN, msg.BL, douyu.GetGift(int64(msg.GFid), int64(msg.Pid)), msg.GFCnt, msg.Hits)
			case "uenter":
				if !U || msg.Level < dt.IntStr(UL) {
					continue
				}
				format := "进入 %" + length(msg.NN, 30) + "s |%3d| %" + length(nl, 4) + "s |"
				logx.Infof(format, msg.NN, msg.Level, nl)
			case "rnewbc":
				if strconv.FormatInt(int64(msg.DRid), 10) != R {
					continue
				}
				format := "续费 %" + length(msg.Unk, 30) + "s |%3d| %" + length(nl, 4) + "s |"
				logx.Infof(format, msg.Unk, msg.Level, nl)
			case "anbc":
				if strconv.FormatInt(int64(msg.DRid), 10) != R {
					continue
				}
				format := "开通 %" + length(msg.Unk, 30) + "s |%3d| %" + length(nl, 4) + "s |"
				logx.Infof(format, msg.Unk, msg.Level, nl)
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
