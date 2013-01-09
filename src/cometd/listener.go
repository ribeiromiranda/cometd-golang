package cometd

import (
	"time"
	"log"
)


type Listener struct {
	listenerRead chan *MetaMessage
	listenerWrite chan *MetaMessage
}

func NewListener() (l Listener) {
	l.listenerRead = make(chan *MetaMessage)
	l.listenerWrite = make(chan *MetaMessage)

	return
}

func (l *Listener) Send(message *MetaMessage) {
	//go func () {
		var messageResponseData *MetaMessage
		messageResponse, session, err := SwitchMeta(message)
		if err == nil && messageResponse.wait {
			for {
				out := false
				select {
					case messageResponseData = <- session.Message:
						l.listenerWrite <- messageResponseData
					case <- time.After(10 * time.Second):
						out = true
				}
				
				if out {
					break
				}
				
			}
		} else if err != nil {
			log.Print("Send error: ", err)
		}
		
		l.listenerWrite <- messageResponse
	//}()
}

func (l *Listener) Receive() (message []*MetaMessage) {
	message = []*MetaMessage{<- l.listenerWrite};
	return 
}

/*
func Listener() (l Listener) {

	go func () {
		for {
			select {
				case m := <- listenerReceive:
				case m := <- listenerDelivery:
				//case <- time.After(5 * time.Second):
				//	return
			}
		}
	}()
	
	return 
}
*/