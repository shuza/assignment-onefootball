package player

/**
 *  := 	create date: 05-Jul-2019
 *  := 	(C) CopyRight Shuza
 *  := 	shuza.ninja
 *  := 	shuza.sa@gmail.com
 *  := 	Code	:	Coffee	: Fun
 **/

type ISource interface {
	Init() error
	FindPlayer(targetTeamNames []string) []Player
}
