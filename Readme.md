# Simserver

It's a file server. You can transfer the file with the machine.

## Getting Started

You can get binary in release page or build the latest version from source code.

## Usage
```
Usage of Simserver:
  -RFN value
        restore files' name
  -chatBoard
        enable chatBoard mode
  -config value
        use 'config.json' args
  -dC
        enable download_code mode
  -https string
        use HTTPS with cer and key
        example: "cer.cer key.pvk"
  -ip string
        set the ip listen (default "0.0.0.0")
  -log
        enable log writing to file
  -login string
        add account password auth for all resource.
        example: "admin:admin"
  -ls
        enable ls mode
  -port string
        set the port listen (default "5000")
  -secureExt
        set secureExt mode(default: true)
  -uT
        enable upload text mode
  -upload
        enable upload mode
  -version value
        show the version
  -view
        enable view mode
  -zip
        enable zip mode
```

explain for commands.

* Normally, you can download file with `ip:port/api/dl/n/filePath`
    - `ip` It is the ip you set or default `0.0.0.0` (below omit)
    - `port` It is the port you set or default `5000`
    - `filePath` It is a path for your file.
      * `C:\\Users\\Administrator\\Desktop\\a.txt` Example for win.
      * `/home/root/desktop/a.txt` Example for Linux.

* `-h` The above help tips are displayed.

* `-version` Print the binary version.

* `-RFN` Reset files which in upload folder to origin's name without secure ext.

* `-chatBoard` Enable chatBoard feature in view mode.

* `-config` load `config.json` args for running.
  - `config.json` Example.
    * ```json
      {
          "download": {
              "ls": false,
              "zip": false,
              "downloadCode": false
          },
          "upload": {
              "enable": false,
              "secureExt": true
          },
          "security": {
              "https": ["Simserver.cer", "Simserver.pvk"],
              "log": true,
              "login": {
                  "enable": false,
                  "account": "",
                  "password": ""
              }
          },
          "view": {
              "enable": true,
              "chatBoard": true
          },
          "ip": "0.0.0.0",
          "port": "5000"
      }
      ```
    * You only have to set some args you want, the other will be set to default.

* `-downloadCode` Enable downloadCode mode.
    - `ip:port/dl/downloadCode/Code` Provides to download a combined zip file of several specific files you set.(`downloadCodes.json` is needed.)
    * `downloadCodes.json` content:
      - ```json
        [
          {
            "Code": "abc",
            "Name": "abc",
            "Files": [
              "D:/a.txt",
              "D:/binary/b.exe",
              "D:/c.pdf"
            ]
          },
          {
            "Code": "cp",
            "Files": [
              "/home/root/b",
              "/home/root/Desktop/c.pdf"
            ]
          }
        ]
        ```
      * `Code` It is `ip:port/dl/downloadCode/Code`'s code.
      * `Name` The file name is displayed in browser.
      * `Files` The file group you need to fill in the code package.

* `-https` use https.
  - `-https server.crt server.key` Example for Win
  - `-https server.pem server.key` Example for Linux
    * You can get crt or key from openssl.

* `-ip ipstr` Set the listen ip.
  - `-ip 0.0.0.0`,`-ip 127.0.0.1` Example.

* `-log` Enable log mode.
  - It will save run logs into `ftps.log`.

* `-login` add basic auth for all resources.
  - `-login account password`
    - account: string without blank string.
    - password: string without blank string.

* `-ls` Enable ls mode.
  - `ip:port/api/dl/ls/folderPath` Open this URL to view the file list corresponding to the path.
    * `folderPath` It is a path for your folder.
      - `C:\\Users\\Administrator\\Desktop` Example for win.
      - `/home/root/desktop` Example for Linux.


* `-port portstr` Set the listen port.
  - `-p 5000`, `-port 5050` Example.

* `-secureExt` set secureExt mode(default: true).
  - It will add `_dat` to the end of file name.

* `-upload` Enable upload mode.
  - `ip:port/upload/` Open this URL to get the upload page. It can upload lots of files once.
    * It will change upload file's extension name.

* `-uT` Enable upload text feature in upload mode.

* `-view` Enable view mode.
  - use view-plugins.
  - Build-in plugins:
    * l: An easy-to-use file manager
    * upload: Upload files 

* `-zip` Enable zip mode.
  - `ip:port/dl/zip/folderPath` Open this URL to download folder with zip format package.

## Built With

* Gin
* Go 1.22

## Authors

* **Anest Larry** 

See also the list of [contributors](https://github.com/AnestLarry/Simserver/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Thanks
[![](./Resources/icon-goland.svg)](https://www.jetbrains.com/)