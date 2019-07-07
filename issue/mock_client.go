package issue

import (
	"github.com/stretchr/testify/mock"
	"os"
)

/**
 *  := 	create date: 06-Jul-2019
 *  := 	(C) CopyRight Shuza
 *  := 	shuza.ninja
 *  := 	shuza.sa@gmail.com
 *  := 	Code	:	Coffee	: Fun
 **/

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Init(f *os.File) error {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(error)
}

func (m *MockLogger) Error(url string, errorType string, err error) {
}

func (m *MockLogger) Info(url string, infoType string, message string) {
}

func (m *MockLogger) Infoln(message string) {

}
