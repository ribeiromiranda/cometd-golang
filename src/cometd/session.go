package cometd


import (
	"fmt"
	"io"
	"crypto/rand"
	"errors"
)
const (
	SessionIDLength  = 30
    SessionIDCharset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

var (
	sessions map[string]*Session = make(map[string]*Session)
)

type Session struct {
	ClientId string
	Connected bool
	Message chan *MetaMessage
}

func NewSession() (*Session, error) {
	id := NewSessionId()
	s := Session{ClientId: id, Connected: false, Message: make(chan *MetaMessage)}
	
	//for  {
	//	if _, present := sessions[s.ClientId]; ! present {
	//		break
	//	}
	//	s = Session{ClientId: NewSessionID()}
	//}
	sessions[s.ClientId] = &s
	return &s, nil
}

func DestroySession(clientId string) error {
	if _, err := GetSession(clientId); err != nil {
		return err
	}
	
	delete(sessions, clientId)
	return nil
}

func GetSession(sessionId string) (s *Session, err error) {
	var present bool
	s, present = sessions[sessionId]
	
	if ! present {
		err = errors.New(fmt.Sprintf("Session with id '%s' does not exist.", sessionId))
	}
	
	return s, err
}

func NewSessionId() string {
	b := make([]byte, SessionIDLength)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
    	return ""
	}

    for i := 0; i < SessionIDLength; i++ {
    	b[i] = SessionIDCharset[b[i]%uint8(len(SessionIDCharset))]
	}

    return string(b)
}