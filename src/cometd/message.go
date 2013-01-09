package cometd

import(
	"fmt"
	"errors"
)

type MetaMessage struct {
    Channel string `json:"channel"`
    Version string `json:"version"`
    MinimumVersion string `json:"minimumVersion"`
    SupportedConnectionTypes []string `json:"supportedConnectionTypes"`
    ClientId string `json:"clientId"`
	Successful bool `json:"successful"`
    AuthSuccessful bool `json:"authSuccessful"`
    Advice *Advice `json:"advice"`
    ConnectionType string `json:"connectionType"`
    Error string `json:"error"`
    Timestamp string `json:"timestamp"`
    Data interface{} `json:"data"`
    Id string `json:"id"`
    Subscription string `json:"subscription"`
    wait bool
}

type Advice struct {
	Reconnect string `json:"reconnect"` 
	Timeout int `json:"timeout"`
	Interval int `json:"interval"`
}

func metaHandshake(messageRequest *MetaMessage) (*MetaMessage) {
	messageResponse := &MetaMessage{}
	session, err := NewSession()
	if err == nil {
		messageResponse.Id = messageRequest.Id
		messageResponse.ClientId = session.ClientId
		messageResponse.Successful = true
		
	} else {
		messageResponse.Successful = false
		messageResponse.Error = "Implements Message Error"
	}

	messageResponse.Channel = "/meta/handshake"
	messageResponse.Version = VERSION
	messageResponse.MinimumVersion = MINIMUM_VERSION
	messageResponse.SupportedConnectionTypes = []string{"long-polling", "callback-polling"}
	messageResponse.SupportedConnectionTypes = []string{"websocket"}
	//messageResponse.Advice = Advice{Reconnect: "retry", Interval:0, Timeout:20000}
	
	return messageResponse
}

func metaConnect(messageRequest *MetaMessage) (messageResponse *MetaMessage) {
	session, err := GetSession(messageRequest.ClientId);
	
	messageResponse = &MetaMessage{}
	messageResponse.Channel = "/meta/connect"
	messageResponse.Id = messageRequest.Id
	
	switch {
		case err != nil:
			messageResponse.ClientId = messageRequest.ClientId
			messageResponse.Advice = &Advice{Reconnect:"handshake", Interval:0, Timeout:20000}
			messageResponse.Successful = false
			messageResponse.Error = fmt.Sprintf("%s", err)

		case ! session.Connected:
			session.Connected = true
		
			messageResponse.ClientId = messageRequest.ClientId
			messageResponse.Advice = &Advice{Reconnect:"retry", Interval:0, Timeout:20000}
			messageResponse.Successful = true
			messageResponse.Error = ""

		default:
			messageResponse.Successful = true
			messageResponse.Error = ""
			messageResponse.wait = true
	}
	

	return
}

func metaDisconnect(messageRequest *MetaMessage) (messageResponse *MetaMessage) {
	messageResponse = &MetaMessage{}
	messageResponse.Channel = messageRequest.Channel
	messageResponse.ClientId = messageRequest.ClientId
	
	if _, err := GetSession(messageRequest.ClientId); err != nil {
		messageResponse.Successful = false
		textError := fmt.Sprintf("%s", err)
		messageResponse.Error = textError
		return 
		
	} else {
		messageResponse.Id = messageRequest.Id
		DestroySession(messageRequest.ClientId)
		messageResponse.Successful = true
		messageResponse.Error = "" 
	}
	
	return 
}

func metaSubscribe(messageRequest *MetaMessage) (messageResponse *MetaMessage) {
	messageResponse = &MetaMessage{}
	session, err := GetSession(messageRequest.ClientId);
	
	messageResponse.Id = messageRequest.Id
	messageResponse.Subscription = messageRequest.Subscription
	messageResponse.Channel = messageRequest.Channel
	messageResponse.Successful = err == nil
	
	if ! messageResponse.Successful {
		return
	}
	
	if messageResponse.Subscription != "" {
		err = errors.New("Value Subscription invalid")
	}
	
	Subscribe(messageRequest.Subscription, session)
	
	return
}

func metaUnsubscribe(messageRequest *MetaMessage) (messageResponse *MetaMessage) {
	messageResponse = &MetaMessage{}
	session, err := GetSession(messageRequest.ClientId);
	
	messageResponse.Id = messageRequest.Id
	messageResponse.Subscription = messageRequest.Subscription
	messageResponse.Channel = messageRequest.Channel
	messageResponse.Successful = err == nil
	
	if ! messageResponse.Successful {
		return
	}
	
	if messageResponse.Subscription != "" {
		err = errors.New("Value Subscription invalid")
	}
	
	Unsubscribe(messageRequest.Subscription, session)
	
	return
}

func metaPublish(messageRequest *MetaMessage) (messageResponse *MetaMessage) {
	messageResponse = &MetaMessage{}
	_, err := GetSession(messageRequest.ClientId);
	
	messageResponse.Id = messageRequest.Id
	messageResponse.Channel = messageRequest.Channel
	messageResponse.Successful = err == nil
	
	if ! messageResponse.Successful {
		return
	}
	
	if ! ChannelIsService(messageRequest.Channel) {
		metaDelivery(messageRequest)
	}
	
	//Unsubscribe(messageRequest.Subscription, session)
	
	return 
}

func metaDelivery(messageRequest *MetaMessage) {
	messageDelivery := &MetaMessage{}
	messageDelivery.Id = messageRequest.Id
	messageDelivery.Channel = messageRequest.Channel
	messageDelivery.Data = messageRequest.Data
	
	Publisher(messageRequest.Channel, messageDelivery)
}

func SwitchMeta(message *MetaMessage) (messageResponse *MetaMessage, session *Session, err error) {
	err = nil
	switch {
		case message.Channel == "/meta/handshake":
			messageResponse = metaHandshake(message)
			
		case message.Channel == "/meta/connect":
			messageResponse = metaConnect(message)

		case message.Channel == "/meta/disconnect":
			messageResponse = metaDisconnect(message)
		
		case message.Channel == "/meta/subscribe":
			messageResponse = metaSubscribe(message)
			
		case message.Channel == "/meta/unsubscribe":
			messageResponse = metaUnsubscribe(message)

		case ChannelExists(message.Channel):
			messageResponse = metaPublish(message)

		default:
			err = errors.New("Invalid channel")
	}
	
	if err == nil {
		session, err = GetSession(message.ClientId)
	}

	return 
}