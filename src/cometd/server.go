package cometd

import (
    "net/http"
    "code.google.com/p/go.net/websocket"
    "encoding/json"
	"io/ioutil"
	"bytes"
	"log"
)

const (
	VERSION = "1.0"
	MINIMUM_VERSION = "1.0"
)


/*

func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("asf", r.Method)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
	}
	
	//body = r.URL.Query().Get("message")
	
	var request MetaMessage
	err = json.Unmarshal(body, &request);
	if err != nil {
		fmt.Printf("error unmarshal request - %s\n", err)
	}
	
	response := MetaMessage{}
	var session *Session
	
	
	switch {
		case r.URL.Path == "/cometd/handshake" || request.Channel == "/meta/handshake":
			response, session= metaHandshake(request)
		     
		case request.Channel == "/meta/connect":
			response, session = metaConnect(request)
			
		case request.Channel == "/meta/disconnect":
			response, session = metaDisconnect(request)
			
		default:
			http.NotFound(w, r)
			return 
			//response, session= metaHandshake(request)
			//response, session = subscribe(request)
	}
	
	messageJson, err := json.Marshal(response)
	if err == nil {
	}

    fmt.Fprintf(w, string(messageJson))
    
    return 
    
    message := Listener(session)
    message.response(&response)
}
*/

func parseMessages(in string) (messagesConverted []MetaMessage, err error) {
	log.Print("parseMessages: ", in)
	//var messagesMap []map[string]interface{}
	 
	messages, err := ioutil.ReadAll(bytes.NewBufferString(in))
	if err != nil {
		return 
	}
	
	err = json.Unmarshal(messages, &messagesConverted);
	if err != nil {
		return 
	}
	
	log.Print("parseMessages: length ", len(messagesConverted)) 
	
	return 
}

func jsonEncode(messages []*MetaMessage) (encoded []byte, err error) {
	jsonBuff, err := json.Marshal(messages)
	if err != nil {
		return 
	}
	
	messagesMap := []map[string]interface{}{}
	err = json.Unmarshal(jsonBuff, &messagesMap)
	if err != nil {
		return 
	}
	
	for key := range messagesMap { 
		message := messagesMap[key]
		for key2, value := range message {
			if value == "" || value == nil || value == false {
				delete(message, key2)
			}
		}
		messagesMap[key] = message
	}
	
	encoded, err = json.Marshal(messagesMap)
	
	return 
}

func handlerWebSocket(ws *websocket.Conn) {
	listener := NewListener()
	channelWs := make(chan *websocket.Conn)

	go func (channelWs <- chan *websocket.Conn) {
		ws := <- channelWs
		for {
	    	var messages []MetaMessage
	        err := websocket.JSON.Receive(ws, &messages)
	        if err != nil {
	        	log.Println("Receive: ", err.Error())
	        	return 
	        }
	        log.Print("RECEIVE: ", messages)
	        
	        for key := range messages {
	        	go listener.Send(&messages[key])
	        }
		}
	}(channelWs)
	
	

    go func (channelWs <- chan *websocket.Conn) {
    	ws := <- channelWs 
    	for {
	        response := listener.Receive()
	        message, err := jsonEncode(response)
	        if err != nil {
	        	return  
	        }
	        
	        log.Print("SEND: ", string(message))
	        
	        err = websocket.Message.Send(ws, string(message))
	        //err = websocket.JSON.Send(ws, response)
	        if err != nil {
	        	log.Println("Send: ", err)
	        	return 
	    	}
    	}
    }(channelWs)
    
    channelWs <- ws
    channelWs <- ws
    channelWs <- ws
    channelWs <- ws
}


/*
func handlerHandshake(w http.ResponseWriter, r *http.Request) {
	log.Print("handlerHandshake: Method", r.Method)
	if r.Method == "OPTIONS" {
		return
		
	} else if r.Method != "GET" {
		http.NotFound(w, r)
		return ;
	}
	
	messages, err := parseMessages(r.URL.Query().Get("message"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	
	response, err := switcMeta(messages)	
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
	
	messageJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

    fmt.Fprintf(w, string(messageJson))
}

func handlerTeste(w http.ResponseWriter, r *http.Request) {
	log.Print("handlerTeste: Method", r.URL.String())
}
*/
func Run() {
	
	http.Handle("/cometd", websocket.Handler(handlerWebSocket))
	
	// long-polling
    //http.HandleFunc("/cometd/handshake", handlerHandshake)
    //http.HandleFunc("/cometd/connect", handlerTeste)
    
    http.ListenAndServe(":8080", nil)
}