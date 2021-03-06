package monkey

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const alwaysMonkeyingAround = 1.0
const neverMonkeyAround = 0.0
const cannedResponse = "hello, world"

func TestItLoadsFromYAML(t *testing.T) {

	yaml := `
---
# Writes a different body 50% of the time
- body: "This is wrong :( "
  frequency: 0.5

# Delays initial writing of response by a second 20% of the time
- delay: 1000
  frequency: 0.2

# Returns a 404 30% of the time
- status: 404
  frequency: 0.3

# Write 10,000,000 garbage bytes 10% of the time
- garbage: 10000000
  frequency: 0.09
`
	degegate := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	monkeyServer, err := NewServerFromYAML(degegate.Config.Handler, []byte(yaml))

	if err != nil {
		t.Fatalf("It didnt return a server from the YAML: %v", err)
	}

	if len(monkeyServer.(*server).behaviours) != 4 {
		t.Error("It didnt load all the behaviours from YAML")
	}

	monkeyServer, _ = NewServer(degegate.Config.Handler, "")

	if monkeyServer != degegate.Config.Handler {
		t.Error("It should just return the server as is if the config path is empty")
	}
}

func TestItMonkeysWithStatusCodesAndBodies(t *testing.T) {
	monkeyBehaviour := new(behaviour)
	monkeyBehaviour.Frequency = alwaysMonkeyingAround
	monkeyBehaviour.Status = http.StatusNotFound
	monkeyBehaviour.Body = "hello, monkey"

	testServer, request := makeTestServerAndRequest()

	monkeyServer := newServerFromBehaviour(testServer.Config.Handler, []behaviour{*monkeyBehaviour})

	w := httptest.NewRecorder()

	monkeyServer.ServeHTTP(w, request)

	if w.Code != monkeyBehaviour.Status {
		t.Error("Server shouldve returned a 404 because of monkey override")
	}

	if w.Body.String() != monkeyBehaviour.Body {
		t.Error("Server should've returned a different body because of monkey override")
	}
}

func TestItReturnsGarbage(t *testing.T) {
	monkeyBehaviour := new(behaviour)
	monkeyBehaviour.Frequency = alwaysMonkeyingAround
	monkeyBehaviour.Garbage = 1984

	testServer, request := makeTestServerAndRequest()

	monkeyServer := newServerFromBehaviour(testServer.Config.Handler, []behaviour{*monkeyBehaviour})

	w := httptest.NewRecorder()

	monkeyServer.ServeHTTP(w, request)

	if len(w.Body.Bytes()) != monkeyBehaviour.Garbage {
		t.Error("Server shouldve returned garbage")
	}
}

func TestItDoesntMonkeyAroundWhenFrequencyIsNothing(t *testing.T) {
	monkeyBehaviour := new(behaviour)
	monkeyBehaviour.Frequency = neverMonkeyAround
	monkeyBehaviour.Body = "blah blah"

	testServer, request := makeTestServerAndRequest()

	monkeyServer := newServerFromBehaviour(testServer.Config.Handler, []behaviour{*monkeyBehaviour})

	w := httptest.NewRecorder()

	monkeyServer.ServeHTTP(w, request)

	if w.Body.String() != cannedResponse {
		t.Error("Server shouldn't have been monkeyed with ")
	}
}

func makeTestServerAndRequest() (*httptest.Server, *http.Request) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cannedResponse))
	}))
	request, _ := http.NewRequest("GET", server.URL, nil)

	return server, request
}
