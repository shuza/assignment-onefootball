package main

import (
	"assignment-onefootball/source"
	"fmt"
)

/**
 *  := 	create date: 05-Jul-2019
 *  := 	(C) CopyRight Shuza
 *  := 	shuza.ninja
 *  := 	shuza.sa@gmail.com
 *  := 	Code	:	Coffee	: Fun
 **/

func main() {
	teamNames := getTargetTeamNames()
	t := source.VintageMonster{}
	t.Init()

	t.FindPlayer(teamNames)

	playerList := t.FindPlayer(teamNames)

	for i, v := range playerList {
		fmt.Printf("%v. %v\n", i+1, v)
	}

}

func getTargetTeamNames() []string {
	return []string{
		"Germany",
		"England",
		"France",
		"Spain",
		"Manchester United",
		"Arsenal",
		"Chelsea",
		"Barcelona",
		"Real Madrid",
		"Bayern Munich",
	}
}
