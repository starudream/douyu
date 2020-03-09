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
	r string
	l int64
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

	flag.StringVar(&r, "r", "", "")
	flag.Int64Var(&l, "l", 30, "")
	flag.Parse()

	if r == "" {
		os.Exit(1)
	}
	if l < 0 {
		l = 0
	}
	if l > 120 {
		l = 120
	}
}

func main() {
	logx.Debugf("room: %s, level: %d", r, l)

	c := douyu.NewClient().SetRoomId(r)

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
				if msg.Level < dt.IntStr(l) {
					continue
				}
				format := "弹幕 %" + length(msg.NN, 30) + "s |%3d| %" + length(msg.BNN, 6) + "s |%3d|: %s"
				logx.Infof(format, msg.NN, msg.Level, msg.BNN, msg.BL, msg.Txt)
			case "dgb":
				if msg.BG == 0 {
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
