package main

import (
	"bytes"
	"configrd"
	"encoding/json"
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
		resp, err := http.Get(conf["CollectorAddr"])
		if err != nil {
			logger.Log.Logln(logger.LEVEL_WARNING, "Unable to connect to collector,", err)
			return
		}
		jdecoder := json.NewDecoder(resp.Body)
		err = jdecoder.Decode(&fresp)
		if err != nil {
			logger.Log.Logln(logger.LEVEL_WARNING, "Unable to unmarshal response,", err)
			return
		}
		for _, item := range fresp.Noti {
			msg := SendMessageParam{
				chat_id:    chatid,
				parse_mode: "Markdown",
			}
			msg.text = "# " + item.Heading + "\n## " + item.Tm.Format(time.UnixDate) + "\n" + item.Content
			msgencoded, _ := json.Marshal(msg)
			resp, err := http.Post("https://api.telegram.org/"+conf["BotId"]+"/sendMessage", "application/json", bytes.NewReader(msgencoded))
			if err != nil {
				logger.Log.Logln(logger.LEVEL_WARNING, "Unable to post to telegram API,", err)
			}
			logger.Log.Logln(logger.LEVEL_INFO, resp.Body)
		}
	}
}
