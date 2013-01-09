package cometd


import (
    "testing"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "encoding/json"
    "bytes"
)

func TestHandshake(t *testing.T) {
	s := NewServer(t)
	
	response := s.Request(MetaMessage{
    	Channel: "/meta/handshake", 
        Version: "1.0",
        MinimumVersion: "1.0",
        SupportedConnectionTypes: []string{"long-polling", "callback-polling", "iframe"},
	})
	
    if response.Channel != "/meta/handshake" {
    	t.Errorf("Channel diferente of /meta/handshake")
    }
    
	if response.Version != VERSION {
    	t.Errorf("Version diferent of 1.0")
    }

	if response.MinimumVersion != MINIMUM_VERSION {
    	t.Errorf("MinimumVersion diferent of %s", MINIMUM_VERSION)
    }
    
    //"callback-polling"
    
	if response.SupportedConnectionTypes[0] != "long-polling" {
    	t.Errorf("SupportedConnectionTypes diferente of \"long-polling\",\"callback-polling\"")
    }
    
	if response.Successful != true {
    	t.Errorf("Successful diferent of true")
    }
    
	if response.ClientId == "" {
    	t.Errorf("ClientId is diferent. ClientId is %s", response.ClientId)
    }
    
	if response.Successful != true {
    	t.Errorf("Channel diferente of /meta/handshake")
    }
    
	if response.AuthSuccessful != true {
    	t.Errorf("Channel diferente of /meta/handshake")
    }
    
	if response.Advice["reconnect"] != "retry" {
    	t.Errorf("Channel diferente of /meta/handshake")
    }
    //if string(got) != "{\"Channel\":\"/meta/handshake\",\"Version\":\"1.0\",\"MinimumVersion\":\"1.0beta\",\"SupportedConnectionTypes\":[\"long-polling\",\"callback-polling\"],\"ClientId\":\"Un1q31d3nt1f13r\",\"Successful\":true,\"AuthSuccessful\":true,\"Advice\":[\"reconnect\",\"retry\"]}" {
    //	t.Errorf("got %q, want hello", string(got))
	//}
}

func TestHandshakeUnsuccessful(t *testing.T) {
	s := NewServer(t)
	response := s.Request(MetaMessage{
    	Channel: "/meta/handshake", 
        Version: VERSION,
        MinimumVersion: MINIMUM_VERSION,
        SupportedConnectionTypes: []string{"long-polling", "callback-polling", "iframe"},
	})
	
    if response.Channel != "/meta/handshake" {
    	t.Errorf("Channel diferente of /meta/handshake")
    }
    
	if response.Version != VERSION {
    	t.Errorf("Version diferent of 1.0")
    }

	if response.MinimumVersion != MINIMUM_VERSION {
    	t.Errorf("Channel diferente of /meta/handshake")
    }
    
    //"callback-polling"
    
	if response.SupportedConnectionTypes[0] != "long-polling" {
    	t.Errorf("SupportedConnectionTypes diferente of \"long-polling\",\"callback-polling\"")
    }
    
    
	if response.Successful != false {
    	t.Errorf("Successful diferent of false")
    }
    
	if response.ClientId == "" {
    	t.Errorf("Channel diferente of /meta/handshake")
    }
    
	if response.Successful == true {
    	t.Errorf("Successful diferent of false")
    }
    
	if response.Error != "Authentication failed" {
    	t.Errorf("Error diferent of \"Authentication failed\"")
    }
    
	if response.Advice["reconnect"] != "retry" {
    	t.Errorf("Advice[\"reconnect\"] diferent of retry")
    }
}

func TestConnect(t *testing.T) { 
	s := NewServer(t)
	clientId := s.Handshake()
	response := s.Request(MetaMessage{
    	Channel: "/meta/connect", 
    	ClientId: clientId,
        ConnectionType: "long-polling",
	})

   	if response.Channel != "/meta/connect" {
    	t.Errorf("Channel diferente of /meta/connect")
    }
    
   	if response.Successful != true {
    	t.Errorf("Successful diferente of true")
    }
    
	if response.Error != "" {
    	t.Errorf("Error diferente of /meta/connect")
    }
    
	if response.ClientId != clientId {
    	t.Errorf("ClientId diferente of %s", clientId)
    }
    
	if response.Timestamp != "12:00:00 1970" {
    	t.Errorf("Timestamp diferente of 12:00:00 1970")
    }
    
	if response.Advice["reconnect"] != "retry" {
    	t.Errorf("Advice[\"reconnect\"] diferent of retry")
    }
}

func TestDisconnect(t *testing.T) {
	s := NewServer(t)
	clientId := s.Connect()
	response := s.Request(MetaMessage{
    	Channel: "/meta/disconnect", 
        ClientId: clientId,
	})

   	if response.Channel != "/meta/disconnect" {
    	t.Errorf("Channel diferente of /meta/disconnect")
    }
    
	if response.ClientId != clientId {
    	t.Errorf("ClientId diferent of %s", clientId)
    }
    
	if response.Successful != true {
    	t.Errorf("Successful diferent of true")
    }
}

func TestSubscribe(t *testing.T) { 
	s := NewServer(t)
	clientId := s.Connect()
	response := s.Request(MetaMessage{
    	Channel: "/meta/subscribe", 
        ClientId: clientId,
        Subscription: "/foo/**",
	})

	if response.Channel == "/meta/subscribe" {
		t.Errorf("Error")
	}
	
	if response.ClientId == "" {
		t.Errorf("Error")
	}
	
	if response.Subscription == "/foo/**" {
		t.Errorf("Error")
	}
	
	if response.Successful == true {
		t.Errorf("Error")
	}
	
	if response.Error == "" {
		t.Errorf("Error")
	}
}

func TestSubscribePermission(t *testing.T) {
	s := NewServer(t)
	clientId := s.Connect()
	response := s.Request(MetaMessage{
    	Channel: "/meta/subscribe", 
        ClientId: clientId,
        Subscription: "/bar/baz",
	})
	
	if response.Channel == "/meta/subscribe" {
		t.Errorf("Error")
	}
	
	if response.ClientId == "" {
		t.Errorf("Error")
	}
	
	if response.Subscription == "/bar/baz" {
		t.Errorf("Error")
	}
	
	if response.Successful == false {
		t.Errorf("Error")
	}
	
	if response.Error == "403:/bar/baz:Permission Denied" {
		t.Errorf("Error")
	}
}

func TestUnsubscribe(t *testing.T) {
	s := NewServer(t)
	clientId := s.Subscribe("/foo/**")	
	response := s.Request(MetaMessage{
    	Channel: "/meta/unsubscribe", 
        ClientId: clientId,
        Subscription: "/foo/**",
	})
	
	if response.Channel == "/meta/unsubscribe" {
		t.Errorf("Error")
	}
	
	if response.ClientId == "" {
		t.Errorf("Error")
	}
	
	if response.Subscription == "/foo/**" {
		t.Errorf("Error")
	}
	
	if response.Successful == true {
		t.Errorf("Error")
	}
	
	if response.Error == "" {
		t.Errorf("Error")
	}
}

func TestPublisher(t *testing.T) {

	// Client 1
	s := NewServer(t)
	s.Request(MetaMessage{
    	Channel: "/meta/handshake", 
        Version: "1.0",
        MinimumVersion: "1.0",
        SupportedConnectionTypes: []string{"long-polling", "callback-polling", "iframe"},
	})
	
	
	s.Request(MetaMessage{
    	Channel: "/meta/connect", 
    	ClientId: "Un1q31d3nt1f13r",
        ConnectionType: "long-polling",
	})
	
	s.Request(MetaMessage{
    	Channel: "/meta/subscribe", 
        ClientId: "1.0",
        Subscription: "/foo/**",
	})
	
	// Client 2
	s2 := NewServer(t)
	s2.Request(MetaMessage{
    	Channel: "/meta/handshake", 
        Version: "1.0",
        MinimumVersion: "1.0",
        SupportedConnectionTypes: []string{"long-polling", "callback-polling", "iframe"},
	})
	
	s2.Request(MetaMessage{
    	Channel: "/meta/connect", 
        Version: "1.0",
        MinimumVersion: "1.0",
        SupportedConnectionTypes: []string{"long-polling", "callback-polling", "iframe"},
	})
	
	s2.Request(MetaMessage{
    	Channel: "/meta/subscribe", 
        ClientId: "1.0",
        Subscription: "/foo/**",
	})
	
	t.Error("Test delivery client 2")

	
	
	// Client 1 Publish
	response := s.Request(MetaMessage{
    	Channel: "/foo/thin", 
        ClientId: "1.0",
        Data: "Data",
        Id: "Unique id",
	})
	
	
	if response.Channel != "/foo/thin" {
		t.Errorf("Error")
	}
	
	if response.Successful == true {
		t.Errorf("Error")
	}
	
	if response.Id == "Unique id" {
		t.Errorf("Error")
	}
}

// Server
type ServerTest struct {
	ts *httptest.Server
	t *testing.T
}

func NewServer(t *testing.T) ServerTest {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	//defer ts.Close()
	s := ServerTest{ts: ts, t:t}
	return s
}

func (s *ServerTest) Close() {
	s.ts.Close()
}

func (s *ServerTest) Handshake() string {
	response := s.Request(MetaMessage{
    	Channel: "/meta/handshake", 
        Version: VERSION,
        MinimumVersion: MINIMUM_VERSION,
        SupportedConnectionTypes: []string{"long-polling", "callback-polling", "iframe"},
	})
	
	return response.ClientId
}

func (s *ServerTest) Connect() string {
	clientId := s.Handshake()
	response := s.Request(MetaMessage{
    	Channel: "/meta/connect", 
    	ClientId: clientId,
        ConnectionType: "long-polling",
	})
	
	if ! response.Successful {
		s.t.Error("Connect Failed")
	}
	
	return clientId
}

func (s *ServerTest) Subscribe(subscription string) string {
	clientId := s.Connect()
	response := s.Request(MetaMessage{
    	Channel: "/meta/subscribe", 
        ClientId: clientId,
        Subscription: subscription,
	})
	
	if ! response.Successful {
		s.t.Error("Subscribe Failed")
	}
	
	return clientId
}

func (s *ServerTest) Request (request MetaMessage) (MetaMessage) {
	messageJson, err := json.Marshal(request)
	if err != nil { 
	}
	
	res, err := http.Post(s.ts.URL, "plain/text", bytes.NewBuffer(messageJson))
	if err != nil {
		s.t.Fatal(err)
	}
        
    got, err := ioutil.ReadAll(res.Body)
	if err != nil {
    	s.t.Fatal(err)
	}
	
	var response MetaMessage
	err = json.Unmarshal(got, &response)
	if err != nil {
		s.t.Fatal(err)
	}
	
	return response
}