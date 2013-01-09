package cometd

import (
	//"time"
	"log"
	"strings"
)

var (
	channels map[string]map[string]*Session = make(map[string]map[string]*Session)
	services map[string]map[string]*Session = make(map[string]map[string]*Session)
)

func Subscribe(cometdChannel string, session *Session) {
	if _, preset := channels[cometdChannel]; ! preset {
		channels[cometdChannel] = map[string]*Session{}
	}
	
	if _, preset := channels[cometdChannel][session.ClientId]; ! preset {
		channels[cometdChannel][session.ClientId] = session
	}	
}

func Unsubscribe(cometdChannel string, session *Session) {
	if _, preset := channels[cometdChannel]; preset {
		if _, preset = channels[cometdChannel][session.ClientId]; preset {
			delete(channels[cometdChannel], session.ClientId)
		}
		if len(channels[cometdChannel]) == 0 {
			delete(channels, cometdChannel)
		}
	}
}

func Publisher(cometdChannel string, message *MetaMessage) {
	
	for key := range channels[cometdChannel] {
		log.Print("Publisher channel ", cometdChannel, " session: ", channels[cometdChannel][key])
		channels[cometdChannel][key].Message <- message
	}
}

func ChannelExists(channel string) (bool) {
	if  ChannelIsService(channel) {
		return true
	}

	log.Print("ChannelExists channel: ", channel)
	log.Print("ChannelExists channels: ", channels)
	_, preset := channels[channel]
	return preset
}

func ChannelIsService(channel string) bool {
	 return strings.HasPrefix(channel, "/service/")
}
