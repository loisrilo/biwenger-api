package main

import (
	"encoding/json"
	"fmt"
	"github.com/magiconair/properties"
	"net/http"
)


const getMyPlayersUrl = "https://biwenger.as.com/api/v2/user?fields=*,lineup(type,playersID),players(*,fitness,team,owner),market(*,-userID),offers,-trophies"
const setMyLineUpUrl = "https://biwenger.as.com/api/v2/user?fields=*"

func GetMyPlayers(w http.ResponseWriter, r *http.Request) {

	headers := getDefaultHeaders(r)
	getlineUpBiwengerResponse := new(GetLineUpBiwengerResponse)
	doRequestAndGetStruct("GET", getMyPlayersUrl, headers, "", &getlineUpBiwengerResponse)
	playerIds := GetPlayerIdsFromPlayers(getlineUpBiwengerResponse.Data.Players)
	fmt.Fprintf(w, SendApiResponse(playerIds))
}

func GetPlayerIdsFromPlayers(players []PlayerBase) []int {
	var playerIds []int
	for _, player := range players {
		playerIds = append(playerIds, player.ID)
	}
	return playerIds
}

func SetLineUp(w http.ResponseWriter, r *http.Request) {

	lineUp := new(LineUp)
	getJsonBody(r, &lineUp)
	jsonLineUp, _ := json.Marshal(StartingEleven{*lineUp})
	headers := getDefaultHeaders(r)
	setLineUpBiwengerResponse := new(SetLineUpBiwengerResponse)
	doRequestAndGetStruct("PUT", setMyLineUpUrl, headers, string(jsonLineUp), &setLineUpBiwengerResponse)
	fmt.Fprintf(w, SendApiResponse(setLineUpBiwengerResponse))
}

func getJsonBody(r *http.Request, target interface{}) {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&target)
	if err != nil {
		panic(err)
	}
}

func getDefaultHeaders(r *http.Request) map[string]string {

	p := properties.MustLoadFile("application.properties", properties.UTF8)

	auth := r.Header.Get("authorization")

	var m = make(map[string]string)
	m["Content-Type"] = "application/json"
	m["authorization"] = auth
	m["x-lang"] = "en"
	m["x-league"] = p.GetString("leagueId", "")
	m["x-user"] = p.GetString("userId", "")
	m["x-version"] = p.GetString("biwengerVersion", "")
	return m
}

func structToJson(entity interface{}) string {
	json, err := json.Marshal(&entity)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(json)
}


type GetLineUpBiwengerResponse struct {
	Status int        `json:"status"`
	Data   LineUpData `json:"data"`
}

type SetLineUpBiwengerResponse struct {
	Status int            `json:"status"`
	Data   LineUpBaseData `json:"data"`
}

type StartingEleven struct {
	L LineUp `json:"lineup"`
}

type LineUp struct {
	Formation string `json:"type"`
	PlayerIds []int  `json:"playersID"`
}

type Owner struct {
	Date  int `json:"date"`
	Price int `json:"price"`
}

type PlayerBase struct {
	ID    int   `json:"id"`
	Owner Owner `json:"owner"`
}

type LineUpData struct {
	ID         int           `json:"id"`
	Name       string        `json:"name"`
	Group      string        `json:"group"`
	Icon       string        `json:"icon"`
	Points     int           `json:"points"`
	Balance    int           `json:"balance"`
	JoinDate   int           `json:"joinDate"`
	LineupDate int           `json:"lineupDate"`
	LineUp     LineUp        `json:"lineup"`
	Market     []interface{} `json:"market"`
	Players    []PlayerBase  `json:"players"`
	Offers     []interface{} `json:"offers"`
}

type LineUpBaseData struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Group      string `json:"group"`
	Icon       string `json:"icon"`
	Points     int    `json:"points"`
	Balance    int    `json:"balance"`
	JoinDate   int    `json:"joinDate"`
	LineupDate int    `json:"lineupDate"`
}
