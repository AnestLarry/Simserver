package argsConfig

import (
	"Simserver/Libs"
	"encoding/json"
	"io/ioutil"
	"os"
)

type ArgConfigStruct struct {
	Ls bool `json:"ls"`
	Dls bool `json:"dls"`
	Zip bool `json:"zip"`
	DownloadCode bool `json:"downloadCode"`
	Upload bool `json:"upload"`
	UploadText bool `json:"uploadText"`
	Https []string `json:"https"`
	Log bool `json:"log"`
	Ip string `json:"ip"`
	Port string `json:"port"`
}

func ArgConfigInit() ArgConfigStruct {
	if !Libs.LibsXIsFile__20201("config.json")  {
		panic("'config.json' is not a file.")
	}
	configJson,err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer configJson.Close()
	byteValue, _ := ioutil.ReadAll(configJson)
	var acs ArgConfigStruct
	err = json.Unmarshal(byteValue, &acs)
	if err != nil {
		panic(err)
	}
	return acs
}