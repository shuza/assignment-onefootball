package player

import "fmt"

/**
 *  := 	create date: 05-Jul-2019
 *  := 	(C) CopyRight Shuza
 *  := 	shuza.ninja
 *  := 	shuza.sa@gmail.com
 *  := 	Code	:	Coffee	: Fun
 **/

type Player struct {
	Name    string
	Details PlayerInfo
}

func (p Player) String() string {
	return fmt.Sprintf("%s; %v", p.Name, p.Details)
}

type PlayerInfo struct {
	Age   string
	Teams []string
}

func (info PlayerInfo) String() string {
	str := ""
	for i, v := range info.Teams {
		if i > 0 {
			str = fmt.Sprintf("%s, %s", str, v)
		} else {
			str = v
		}
	}

	return fmt.Sprintf("%s; %s", info.Age, str)
}
