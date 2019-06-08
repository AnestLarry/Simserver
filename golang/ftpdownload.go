package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	Command  string
	Otherkey string
	Path     string
)

func main() {
	if Exists("ftpdownload_setting.ini") {
		f, err := os.OpenFile("ftpdownload_setting.ini", os.O_RDONLY, 0600)
		defer f.Close()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			contentByte, _ := ioutil.ReadAll(f)
			Path = string(contentByte)
		}
	} else {
		f, err := os.Create("ftpdownload_setting.ini")
		defer f.Close()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			_, _ = f.Write([]byte("http://127.0.0.1:5000/download/"))
			Path = "http://127.0.0.1:5000/download/"
		}
	}

	for {
		fmt.Print(Path + "  $ ")
		fmt.Scanln(&Command, &Otherkey)
		if Command != "" {
			switch Command {
			case "cd":
				f, err := os.Create("ftpdownload_setting.ini")
				defer f.Close()
				if err != nil {
					fmt.Println(err.Error())
				} else {
					_, err = f.Write([]byte(Otherkey))
					Path = Otherkey
				}
			case "cdd":
				f, err := os.Create("ftpdownload_setting.ini")
				defer f.Close()
				if err != nil {
					fmt.Println(err.Error())
				} else {
					_, err = f.Write([]byte(Path + Otherkey))
					Path = Path + Otherkey
				}
			case "pwd":
				fmt.Println(Path)
			default:
				FileNameIndex := -1
				if strings.LastIndex(Path+Command, "/") > strings.LastIndex(Path+Command, "\\") {
					FileNameIndex = strings.LastIndex(Path+Command, "/")
				} else {
					FileNameIndex = strings.LastIndex(Path+Command, "\\")
				}
				FileName := Path + Command
				httpGet(FileName, FileName[FileNameIndex+1:])
			}
			Command = ""
		}
	}
}

func httpGet(url string, filename string) {
	resp, err := http.Get(url)
	if err != nil {
	}

	defer resp.Body.Close()

	out, _ := os.Create(filename)
	defer out.Close()

	_, _ = io.Copy(out, resp.Body)
	fmt.Println(filename + " is Finished.")
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
