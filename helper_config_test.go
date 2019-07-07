package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

/**
 *  := 	create date: 07-Jul-2019
 *  := 	(C) CopyRight Shuza
 *  := 	shuza.ninja
 *  := 	shuza.sa@gmail.com
 *  := 	Code	:	Coffee	: Fun
 **/

func TestParseTargetTeamNames(t *testing.T) {
	Convey("Invalid config file path to parse target team list", t, func() {
		_, err := ParseTargetTeamNames("asd", "json", "config")
		Convey("Should return error", func() {
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Valid config file path to parse target team list", t, func() {
		data, err := ParseTargetTeamNames(".", "json", "config")
		Convey("Tram name should not empty", func() {
			So(err, ShouldBeNil)
			So(len(data), ShouldNotBeZeroValue)
		})
	})
}

func TestParseSource(t *testing.T) {
	Convey("Valid config file path to parse plugins", t, func() {
		plugList, err := ParseSource(".", "json", "config")
		Convey("plugin list should not be empty", func() {
			So(err, ShouldBeNil)
			So(len(plugList), ShouldEqual, 1)
		})
	})
}
