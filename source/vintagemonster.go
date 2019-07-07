package main

import (
	"assignment-onefootball/issue"
	"assignment-onefootball/player"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

/**
 *  := 	create date: 05-Jul-2019
 *  := 	(C) CopyRight Shuza
 *  := 	shuza.ninja
 *  := 	shuza.sa@gmail.com
 *  := 	Code	:	Coffee	: Fun
 **/

type VintageMonster struct {
	baseUrl     string
	startFileNo uint16
	endFileNo   uint16
	rateLimit   uint16
	logClient   issue.ILog
	isBlocked   bool
}

type vintageBaseResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Data   struct {
		Team struct {
			Id      int    `json:"id"`
			Name    string `json:"name"`
			Players []struct {
				Id       string `json:"id"`
				FullName string `json:"name"`
				Age      string `json:"age"`
			} `json:"players"`
		} `json:"team"`
	} `json:"data"`
	Message string `json:"message"`
}

func (resp *vintageBaseResponse) isOk() bool {
	return strings.ToLower(resp.Status) == "ok"
}

//	Generate long wait time between 4 to 7 minute
func (s *VintageMonster) getLongWaitDuration() time.Duration {
	delta := rand.Intn(3)
	minute := 4 + delta
	return time.Duration(minute) * time.Minute
}

func (s *VintageMonster) generateUrl(fileNo uint16) string {
	return fmt.Sprintf("%s/%d.json", s.baseUrl, fileNo)
}

//	teamList should be sorted string list
//	Use binary search to find the appropriate position for current team name in target team list
//	Check if it actually present in the target team list
func (s *VintageMonster) isListedTeam(teamList []string, name string) bool {
	index := sort.SearchStrings(teamList, name)
	return index < len(teamList) && teamList[index] == name
}

func (s *VintageMonster) Init(logClient issue.ILog) error {
	s.baseUrl = "https://vintagemonster.onefootball.com/api/teams/en"
	s.startFileNo = 1
	s.endFileNo = 13660
	//	s.endFileNo = ^uint16(0)
	s.rateLimit = 100
	s.logClient = logClient
	return nil
}

func (s *VintageMonster) FindPlayer(targetTeamNames []string) []player.Player {
	processStartTime := time.Now()
	defer func() {
		processDuration := time.Now().Sub(processStartTime)
		s.logClient.Info("Duration", "total_process", fmt.Sprintf("%v", processDuration))
	}()

	//	Sort team names so that we can use binary search
	sort.Strings(targetTeamNames)

	playerList := make([]player.Player, 0)
	playerCh := make(chan player.Player, 10)

	var wg sync.WaitGroup

	for i := s.startFileNo; i <= s.endFileNo; i++ {
		if i != 0 && i%s.rateLimit == 0 {
			//	Reach rate limit
			//	Sleep sometime to bypass rate limit
			//	Sleep time < responseTime * rateLimit
			time.Sleep(5 * time.Second)
			resp, err := s.doGet(s.generateUrl(i))
			if err != nil && resp.Code == 403 {
				//	Blocked by server, wait for a long period and try again
				waitDuration := s.getLongWaitDuration()
				s.logClient.Info("sleep", "blocked", fmt.Sprintf("wait until %v unblocked", waitDuration))
				time.Sleep(s.getLongWaitDuration())
				i--
				continue
			}
		}

		url := s.generateUrl(i)
		wg.Add(1)
		go s.parseTeamPlayer(targetTeamNames, url, &wg, playerCh)
	}

	go s.closePlayerChan(&wg, playerCh)

	playerNames := make([]string, 0)
	playerMap := make(map[string]player.PlayerInfo)

	for p := range playerCh {
		if v, ok := playerMap[p.Name]; ok {
			p.Details.Teams = append(v.Teams, p.Details.Teams...)
		} else {
			//	Add new player name
			playerNames = append(playerNames, p.Name)
		}
		playerMap[p.Name] = p.Details
	}

	//	Sort player name and prepare player list
	sort.Strings(playerNames)
	for _, v := range playerNames {
		playerList = append(playerList, player.Player{
			Name:    v,
			Details: playerMap[v],
		})
	}

	return playerList
}

//	Wait until all goroutines finish their task
//	then close the channel
func (s *VintageMonster) closePlayerChan(wg *sync.WaitGroup, playerCh chan player.Player) {
	wg.Wait()
	close(playerCh)
}

//	Must call wg.Done() when method is finished
//	Otherwise deadlock will arise
func (s *VintageMonster) parseTeamPlayer(targetTeamNames []string, url string, wg *sync.WaitGroup, playerCh chan<- player.Player) {
	defer wg.Done()

	response, err := s.doGet(url)
	if err != nil {
		s.logClient.Error(url, "api_request", err)
		return
	}
	if !response.isOk() {
		//	API don't provide data so skip
		s.logClient.Info(url, "response_code", fmt.Sprintf("response code %s message %s", response.Status, response.Message))
		return
	}

	if !s.isListedTeam(targetTeamNames, response.Data.Team.Name) {
		return
	}

	for _, v := range response.Data.Team.Players {
		player := player.Player{
			Name: v.FullName,
			Details: player.PlayerInfo{
				Age:   v.Age,
				Teams: []string{response.Data.Team.Name},
			},
		}

		playerCh <- player
	}
}

//	Hit the URL and parse response
//	Use Decode() instead of Unmarshal()
func (s *VintageMonster) doGet(url string) (vintageBaseResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return vintageBaseResponse{}, err
	}
	defer resp.Body.Close()

	var response vintageBaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return vintageBaseResponse{
			Status: resp.Status,
			Code:   resp.StatusCode,
		}, err
	}

	return response, nil
}

var INSTANCEVINTAGEMONOSTER VintageMonster
