# Simserver

It's a file server. You can transfer the file with the machine.

## Getting Started

You can get binary in release page or build the latest version from source code.

## Usage
```
Tips:
 h  - show this help
 v  - get version
Mode:
 ls  - open ls function
 upload  - allow user upload data to host
 zip  - allow zip dir for downloadGroup (DANGER!)
 https  - use https with crt and key 
 log  - put log in file
 downloadCode  - use downloadGroup code to downloadGroup a group file with setting
 view  - use view in running
 login  - add basic auth for all resources
Args:
 p / port  - use the port
 ip  - use the ip
 config  - use 'config.json' args
Task:
 RFN  - restore files' name
```

explain for commands.
* Normally, you can download file with `ip:port/api/dl/n/filePath`
    - `ip` It is the ip you set or default `0.0.0.0` (below omit)
    - `port` It is the port you set or default `5000` (below omit)
    - `filePath` It is a path for your file.
      * `C:\\Users\\Administrator\\Desktop\\a.txt` Example for win.
      * `/home/root/desktop/a.txt` Example for Linux.
* `-h` The above help tips are displayed.
* `-v` Print the binary version.
* `-ls` Open ls mode.
  - `ip:port/api/dl/ls/folderPath` Open this URL to view the file list corresponding to the path.
    * `folderPath` It is a path for your folder.
      - `C:\\Users\\Administrator\\Desktop` Example for win.
      - `/home/root/desktop` Example for Linux.

* `-upload` Open upload mode.
  - `ip:port/upload/` Open this URL to get the upload page. It can upload lots of files once.
    * It will change upload file's extension name.
* `-zip` Open zip mode.
  - `ip:port/dl/zip/folderPath` Open this URL to download folder with zip format package.
* `-https` use https.
  - `-https server.crt server.key` Example for Win
  - `-https server.pem server.key` Example for Linux
    * You can get crt or key from openssl.
* `-log` Open log mode.
  - It will save run logs into `ftps.log`.
* `-downloadCode` Open downloadCode mode.
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
* `-view` Open view mode.
  - use view-plugins.
  - Build-in plugins:
    * h5player :  xgplayer@2.9.6
    * photoViewer : easy to show photos which in a folder
    * l : view of folder(origin `dls`)
* `-login` add basic auth for all resources.
  - `-login account password`
    - account: string without blank string.
    - password: string without blank string.
* `-ip ipstr` Set the listen ip.
  - `-ip 0.0.0.0`,`-ip 127.0.0.1` Example.
* `-p portstr` or `-port portstr` Set the listen port.
  - `-p 5000`, `-port 5050` Example.
* `-config` load `config.json` args for running.
  - `config.json` Example.
    * ```json
      {
          "ls": true,
          "zip": false,
          "log": true,
          "upload": true,
          "downloadCode": true,
          "https": ["Simserver.cer", "Simserver.pk"],
          "ip": "0.0.0.0",
          "port": "5000",
          "view": true
      }
      ```
    * You only have to set some args you want, the other will be set to default.
* `-RFN` reset files which in upload folder to origin's name.

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