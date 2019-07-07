package main

import (
	"assignment-onefootball/issue"
	"fmt"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

/**
 *  := 	create date: 06-Jul-2019
 *  := 	(C) CopyRight Shuza
 *  := 	shuza.ninja
 *  := 	shuza.sa@gmail.com
 *  := 	Code	:	Coffee	: Fun
 **/

func TestVintageBaseResponse_isOk(t *testing.T) {
	Convey("Testing response status", t, func() {
		response := vintageBaseResponse{Status: "ok"}
		Convey("Response status should be ok", func() {
			So(response.isOk(), ShouldBeTrue)
		})

		response = vintageBaseResponse{Status: "no_ok"}
		Convey("Response code should not be ok", func() {
			So(response.isOk(), ShouldBeFalse)
		})
	})
}

func TestVintageMonster_doGet(t *testing.T) {
	model := VintageMonster{}
	Convey("Given request for valid API response", t, func() {
		url := "https://vintagemonster.onefootball.com/api/teams/en/1.json"
		resp, err := model.doGet(url)
		Convey("Should return valid response", func() {
			So(err, ShouldBeNil)
			So(resp.isOk(), ShouldBeTrue)
		})
	})

	Convey("Given request for invalid API response", t, func() {
		url := "https://vintagemonster.onefootball.com/api/teams/en/0.json"
		resp, err := model.doGet(url)
		Convey("API request should be successful but invalid response", func() {
			So(err, ShouldBeNil)
			So(resp.isOk(), ShouldBeFalse)
		})
	})

	Convey("Given request to invalid URL", t, func() {
		url := "https://vintagemonster.onefootball.com/api/teams/en/0.json/tmp"
		_, err := model.doGet(url)
		Convey("API request should be failed", func() {
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given request to invalid URL", t, func() {
		url := "invalid-url"
		_, err := model.doGet(url)
		Convey("API request should be failed", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestVintageMonster_Init(t *testing.T) {
	model := VintageMonster{}
	Convey("Calling Init() method", t, func() {
		Convey("Init should return no error", func() {
			err := model.Init(&issue.StdLog{})
			So(err, ShouldBeNil)
		})
	})
}

func TestVintageMonster_generateUrl(t *testing.T) {
	model := VintageMonster{}
	model.Init(&issue.MockLogger{})
	Convey("Generating URL", t, func() {
		url := model.generateUrl(1)
		Convey("URL should be https://vintagemonster.onefootball.com/api/teams/en/2.json", func() {
			So(url, ShouldEqual, "https://vintagemonster.onefootball.com/api/teams/en/1.json")
		})
	})
}

func TestVintageMonster_getLongWaitDuration(t *testing.T) {
	model := VintageMonster{}
	model.Init(&issue.MockLogger{})
	Convey("Test random long wait time generator", t, func() {
		a := model.getLongWaitDuration()
		b := model.getLongWaitDuration()
		Convey("a and b should not be equal", func() {
			So(a, ShouldNotEqual, b)
		})
	})
}

func TestVintageMonster_FindPlayer(t *testing.T) {
	model := VintageMonster{}
	model.Init(&issue.MockLogger{})
	Convey("FindPlayer from file 2.json", t, func() {
		model.endFileNo = 2
		playerList := model.FindPlayer([]string{"Arsenal"})
		Convey("PlayerList should contain 33 player", func() {
			So(len(playerList), ShouldEqual, 33)
		})
	})

	Convey("FindPlayer from file 0.json to 150.json", t, func() {
		model.endFileNo = 150
		playerList := model.FindPlayer([]string{
			"Arsenal", "England", "France", "Spain", "Manchester United", "Arsenal", "Chelsea", "Barcelona", "Real Madrid", "Bayern Munich",
		})
		Convey("PlayerList should contain 55 player", func() {
			So(len(playerList), ShouldEqual, 241)
			for i, v := range playerList {
				str := fmt.Sprintf("%d.  %v", i+1, v)
				fmt.Println(str)
				logrus.Printf("%d.  %v\n", i+1, v)
			}
		})
	})
}

func TestVintageMonster_doGet_tmp(t *testing.T) {
	model := VintageMonster{}
	Convey("Given request for error url", t, func() {
		url := "https://vintagemonster.onefootball.com/api/teams/en/789.json"
		resp, err := model.doGet(url)
		Convey("Error url test", func() {
			So(err, ShouldBeNil)
			So(resp.isOk(), ShouldBeTrue)
		})
	})
}
