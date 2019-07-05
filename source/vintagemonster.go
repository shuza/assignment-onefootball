package source

import (
	"assignment-onefootball/player"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
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
	baseUrl       string
	startFileNo   uint8
	endFileNo     uint8
	clientTimeOut time.Duration
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

func (s *VintageMonster) Init() error {
	s.baseUrl = "https://vintagemonster.onefootball.com/api/teams/en"
	s.startFileNo = 0
	s.endFileNo = 3
	//	s.endFileNo = ^uint8(0)
	s.clientTimeOut = 3 * time.Second
	fmt.Printf("Iterate from  %v  to  %v\n", s.startFileNo, s.endFileNo)
	return nil
}

func (s *VintageMonster) generateUrl(fileNo uint8) string {
	return fmt.Sprintf("%s/%d.json", s.baseUrl, fileNo)
}

func (s *VintageMonster) FindPlayer(targetTeamNames []string) []player.Player {
	//	Sort team names so that we can use binary search
	sort.Strings(targetTeamNames)

	playerList := make([]player.Player, 0)
	playerCh := make(chan player.Player)
	wg := sync.WaitGroup{}

	for i := s.startFileNo; i <= s.endFileNo; i++ {
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
		log.Warnf("URL %s doGet Error : %v\n", url, err)
		return
	}
	if !response.isOk() {
		//	API don't provide data so skip
		log.Warnf("URL %s status not ok\n", url)
		return
	}

	//	Use binary search to find the appropriate position for current team name in target team list
	//	Check if it actually present in the target team list
	if index := sort.SearchStrings(targetTeamNames, response.Data.Team.Name); targetTeamNames[index] != response.Data.Team.Name {
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
	client := http.Client{Timeout: s.clientTimeOut}
	resp, err := client.Get(url)
	if err != nil {
		return vintageBaseResponse{}, err
	}
	defer resp.Body.Close()

	var response vintageBaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return vintageBaseResponse{}, err
	}

	return response, nil
}
