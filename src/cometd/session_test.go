package cometd


import (
    "testing"
)


func TestNewSessionID(t *testing.T) {

	id := NewSessionId()
	if id == "" {
		t.Fail()
	}
	
}