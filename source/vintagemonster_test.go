package source

import (
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
