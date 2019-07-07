package main

import (
	"assignment-onefootball/issue"
	"assignment-onefootball/player"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"plugin"
	"strings"
)

/**
 *  := 	create date: 05-Jul-2019
 *  := 	(C) CopyRight Shuza
 *  := 	shuza.ninja
 *  := 	shuza.sa@gmail.com
 *  := 	Code	:	Coffee	: Fun
 **/

func main() {
	configPath := flag.String("config", ".", "path of config.json file")
	flag.Parse()

	teamNames, err := ParseTargetTeamNames(*configPath, "json", "config")
	if err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile("assignment_log.txt", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}

	resultFile, err := os.OpenFile("result.txt", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}

	//	Initiate log client to store log message
	logClient := issue.StdLog{}
	logClient.Init(logFile)

	//	Initiate log client to store result
	outputClient := issue.StdLog{}

	pluginList, err := ParseSource(*configPath, "json", "config")
	if err != nil {
		panic(err)
	}

	playerList := make([]player.Player, 0)

	for i, v := range pluginList {
		if err := v.Init(&logClient); err != nil {
			logClient.Error(fmt.Sprintf("Plugin No - %d", i), "Plugin Init", err)
			continue
		}
		players := v.FindPlayer(teamNames)
		playerList = append(playerList, players...)
	}

	outputClient.Init(resultFile)
	for i, v := range playerList {
		str := fmt.Sprintf("%d.  %v", i+1, v)
		fmt.Println(str)
		outputClient.Infoln(str)
	}

}

func ParseTargetTeamNames(path string, configType string, fileName string) ([]string, error) {
	viper.AddConfigPath(path)
	viper.SetConfigType(configType)
	viper.SetConfigName(fileName)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	targetName := viper.GetStringSlice("target_team")
	return targetName, nil
}

func ParseSource(path string, configType string, fileName string) ([]player.ISource, error) {
	viper.AddConfigPath(path)
	viper.SetConfigType(configType)
	viper.SetConfigName(fileName)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	sources := make([]player.ISource, 0)
	srcConfig := viper.GetStringMapString("source")
	for k, v := range srcConfig {
		plug, err := plugin.Open(v)
		if err != nil {
			//	invalid plugin file
			continue
		}

		sym, err := plug.Lookup(strings.ToUpper(k))
		if err != nil {
			//	plugin not export correct access key
			continue
		}

		source, ok := sym.(player.ISource)
		if ok {
			sources = append(sources, source)
		}
	}

	return sources, nil
}
