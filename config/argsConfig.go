package argsConfig

import (
	"Simserver/Libs"
	"encoding/json"
	"io"
	"os"
	"sync"
)

type downloadArgConfigStruct struct {
	Ls           bool `json:"ls"`
	Zip          bool `json:"zip"`
	DownloadCode bool `json:"downloadCode"`
}
type uploadArgConfigStruct struct {
	Enable    bool `json:"enable"`
	SecureExt bool `json:"secureExt"`
}
type securityLoginArgConfigStruct struct {
	Enable   bool   `json:"enable"`
	Account  string `json:"account"`
	Password string `json:"password"`
}
type securityArgConfigStruct struct {
	Https []string                     `json:"https"`
	Log   bool                         `json:"log"`
	Login securityLoginArgConfigStruct `json:"login"`
}
type viewArgConfigStruct struct {
	Enable    bool `json:"enable"`
	ChatBoard bool `json:"chatBoard"`
}
type ArgConfigStruct struct {
	Download downloadArgConfigStruct `json:"download"`
	Upload   uploadArgConfigStruct   `json:"upload"`
	Security securityArgConfigStruct `json:"security"`
	View     viewArgConfigStruct     `json:"view"`
	Ip       string                  `json:"ip"`
	Port     string                  `json:"port"`
}

var (
	acs = ArgConfigStruct{
		Download: downloadArgConfigStruct{
			Ls:           false,
			Zip:          false,
			DownloadCode: false,
		},
		Upload: uploadArgConfigStruct{
			Enable:    false,
			SecureExt: true,
		},
		Security: securityArgConfigStruct{
			Https: []string{},
			Log:   false,
			Login: securityLoginArgConfigStruct{
				Enable:   false,
				Account:  "",
				Password: "",
			},
		},
		View: viewArgConfigStruct{
			Enable:    false,
			ChatBoard: false,
		},
		Ip:   "0.0.0.0",
		Port: "5000",
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
	})()
	return acs
}
