package issue

import (
	log "github.com/sirupsen/logrus"
	"os"
)

/**
 *  := 	create date: 06-Jul-2019
 *  := 	(C) CopyRight Shuza
 *  := 	shuza.ninja
 *  := 	shuza.sa@gmail.com
 *  := 	Code	:	Coffee	: Fun
 **/

type StdLog struct{}

func (l *StdLog) Init(f *os.File) error {
	log.SetOutput(f)

	return nil
}

func (l *StdLog) Error(url string, errorType string, err error) {
	log.Warnf("%s  ==  %s  :://  %s\n", url, errorType, err)
}

func (l *StdLog) Info(url string, infoType string, message string) {
	log.Infof("%s  ==  %s   :://   %s\n", url, infoType, message)
}

func (l *StdLog) Infoln(message string) {
	log.Println(message)
}
