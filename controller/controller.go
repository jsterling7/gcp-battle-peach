package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jsterling7/gcp-battle-peach/gameEngine"
	"github.com/jsterling7/gcp-battle-peach/model"
)

func Root(responseWriter http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" || request.Method != "GET" {
		http.NotFound(responseWriter, request)
		return
	}

	fmt.Fprint(responseWriter, "Welcome to Joshua Sterling's Battle Peach Microservice")
}

func Battle(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.NotFound(responseWriter, request)
		return
	}

	//close response body at end of function
	defer request.Body.Close()

	//read response body
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
	}

	gameState := model.GameState{}

	err = json.Unmarshal(body, &gameState)

	if err != nil {
	}

	action := gameEngine.Play(gameState)

	fmt.Fprint(responseWriter, action)
}
