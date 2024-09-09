package argsConfig

import (
	"Simserver/Libs"
	"encoding/json"
	"io"
	"os"
	"sync"
)

type ArgConfigStruct struct {
	Ls           bool     `json:"ls"`
	Zip          bool     `json:"zip"`
	DownloadCode bool     `json:"downloadCode"`
	Upload       bool     `json:"upload"`
	ChatBoard    bool     `json:"chatBoard"`
	Https        []string `json:"https"`
	Log          bool     `json:"log"`
	Ip           string   `json:"ip"`
	Port         string   `json:"port"`
	View         bool     `json:"view"`
	Login        struct {
		Open     bool   `json:"open"`
		Account  string `json:"account"`
		Password string `json:"password"`
	} `json:"login"`
}

var (
	acs = ArgConfigStruct{
		Ls:           false,
		Zip:          false,
		DownloadCode: false,
		Upload:       false,
		ChatBoard:    false,
		Https:        []string{},
		Log:          false,
		Ip:           "0.0.0.0",
		Port:         "5000",
		View:         false,
		Login: struct {
			Open     bool   `json:"open"`
			Account  string `json:"account"`
			Password string `json:"password"`
		}{Open: false, Account: "", Password: ""},
	}
)

func loadConfig() {
	if Libs.LibsXExists("config.json") {
		if Libs.LibsXIsDir("config.json") {
			panic("'config.json' is not a file.")
		}
		configJson, err := os.Open("config.json")
		if err != nil {
			panic(err)
		}
		defer configJson.Close()
		byteValue, _ := io.ReadAll(configJson)
		err = json.Unmarshal(byteValue, &acs)
		if err != nil {
			panic(err)
		}
	} else {

		acsJson, err := json.MarshalIndent(acs, "", "    ")
		if err != nil {
			panic(err)
		}
		err = os.WriteFile("config.json", acsJson, 0664)
		if err != nil {
			panic(err)
		}
	}
}
func ArgConfigInit() ArgConfigStruct {
	sync.OnceFunc(func() {
		loadConfig()
	})
	return acs
}
