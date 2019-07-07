package issue

import "os"

/**
 *  := 	create date: 06-Jul-2019
 *  := 	(C) CopyRight Shuza
 *  := 	shuza.ninja
 *  := 	shuza.sa@gmail.com
 *  := 	Code	:	Coffee	: Fun
 **/

type ILog interface {
	Init(f *os.File) error
	Error(url string, errorType string, err error)
	Info(url string, infoType string, message string)
	Infoln(message string)
}
