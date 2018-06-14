package main

import (
	"bytes"
	"configrd"
	"encoding/json"
	"io"
	"logger"
	"net/http"
	"os"
	"strconv"
	"time"
)

var conf map[string]string

func main() {
	var confile = configrd.Config(os.Args[1])
	conf = confile.ReadConfig()
	_, mpok := conf["BotId"]
	if !mpok {
		panic("BotId has to be exist in config map!")
	}
	_, mpok = conf["ChatId"]
	if !mpok {
		panic("ChatId has to be exist in config map!")
	}
	_, mpok = conf["CollectorAddr"]
	if !mpok {
		panic("CollectorAddr has to be exist in config map!")
	}
	_, mpok = conf["Timeout"]
	if !mpok {
		panic("Timeout has to be exist in config map!")
	}
	timeout, err := strconv.Atoi(conf["Timeout"])
	if err != nil {
		panic("Timeout is not a valid integer!")
	}
	chatid, cerr := strconv.Atoi(conf["ChatId"])
	if cerr != nil {
		panic("ChatId is not a valid integer!")
	}
	// LOOP START!
	for {
		fresp := FetchRepsonse{
			Serv: make([]Service, 0),
			Stat: make([]Status, 0),
			Noti: make([]Notification, 0),
		}
		<-time.After(time.Duration(timeout) * time.Millisecond)
		// logger.Log.Logln(logger.LEVEL_DEBUG, "Looped!")
		resp, err := http.Get(conf["CollectorAddr"])
		if err != nil {
			logger.Log.Logln(logger.LEVEL_WARNING, "Unable to connect to collector,", err)
			continue
		}
		// logger.Log.Logln(logger.LEVEL_DEBUG, resp)
		jdecoder := json.NewDecoder(resp.Body)
		err = jdecoder.Decode(&fresp)
		// logger.Log.Logln(logger.LEVEL_INFO, "decode completed")
		if err != nil {
			logger.Log.Logln(logger.LEVEL_WARNING, "Unable to unmarshal response,", err)
			continue
		}
		logger.Log.Logln(logger.LEVEL_INFO, "fresp", fresp)
		for _, item := range fresp.Noti {
			// logger.Log.Logln(logger.LEVEL_INFO, item)
			msg := make(map[string]string)
			msg["chat_id"] = strconv.Itoa(chatid)
			msg["parse_mode"] = "Markdown"
			msg["text"] = "*" + item.Heading + "*\n_" + item.Tm.Format(time.UnixDate) + "_\n" + item.Content
			msgencoded, _ := json.Marshal(msg)
			// logger.Log.Logln(logger.LEVEL_INFO, "to be sent:", string(msgencoded))
			resp, err := http.Post("https://api.telegram.org/"+conf["BotId"]+"/sendMessage", "application/json", bytes.NewReader(msgencoded))
			if err != nil {
				logger.Log.Logln(logger.LEVEL_WARNING, "Unable to post to telegram API,", err)
				continue
			}
			buffer := make([]byte, 999)
			io.ReadFull(resp.Body, buffer)
			logger.Log.Logln(logger.LEVEL_INFO, "tg response:", string(buffer))
		}
	}
}
