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
 dls  - add downloadGroup links with the ls function's list
 upload  - allow user upload files to host
 uploadText  - allow user fill textarea to save text in txt
 zip  - allow zip dir for downloadGroup (DANGER!)
 https  - use https with crt and key 
 log  - put log in file
 downloadCode  - use downloadGroup code to downloadGroup a group file with setting
 view  - use view in running
Args:
 p / port  - use the port
 ip  - use the ip
 config  - use 'config.json' args
Task:
 RFN  - restore files' name
```

explain for commands.
* Normally, you can download file with `ip:port/dl/n/filePath`
    - `ip` It is the ip you set or default `0.0.0.0` (below omit)
    - `port` It is the port you set or default `5000` (below omit)
    - `filePath` It is a path for your file.
      * `C:\\Users\\Administrator\\Desktop\\a.txt` Example for win.
      * `/home/root/desktop/a.txt` Example for Linux.
* `-h` The above help tips are displayed.
* `-v` Print the binary version.
* `-ls` Open ls mode.
  - `ip:port/dl/ls/folderPath` Open this URL to view the file list corresponding to the path.
    * `folderPath` It is a path for your folder.
      - `C:\\Users\\Administrator\\Desktop` Example for win.
      - `/home/root/desktop` Example for Linux.
* `-dls` Open dls mode.(dls: `ls` with download function)
    - `ip:port/dl/dls/folderPath` Open this URL to get the file list corresponding with download link to the path.
* `-upload` Open upload mode.
  - `ip:port/upload/` Open this URL to get the upload page. It can upload lots of files once.
    * It will change upload file's extension name.
* `-uploadText` Open upload text mode.
  - `ip:port/upload/text` Open this URL to get the upload text page. It can upload text without a txt file.
* `-zip` Open zip mode.
  - `ip:port/dl/zip/folderPath` Open this URL to download folder with zip format package.
* `https` use https.
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
  - Default plugins:
    * h5player :  xgplayer@2.9.6
* `-ip ipstr` Set the listen ip.
  - `-ip 0.0.0.0`,`-ip 127.0.0.1` Example.
* `-p portstr` or `-port portstr` Set the listen port.
  - `-p 5000`, `-port 5050` Example.
* `-config` load `config.json` args for running.
  - `config.json` Example.
    * ```json
      {
      "ls": true,
      "dls": true,
      "downloadCode": true,
      "zip": true,
      "upload": true,
      "uploadText": true,
      "ip": "0.0.0.0",
      "port": "5000",
      "https": ["Simserver.cer", "Simserver.pk"]
      }
      ```
    * You only have to set some args you want, the other will be set to default.
* `-RSUN` reset files which in upload folder to origin's name.

## Built With

* Gin
* Go 1.16
* xgplayer 2.9.6

## Authors

* **Anest Larry** 

See also the list of [contributors](https://github.com/AnestLarry/Simserver/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Thanks
[![](./Resources/icon-goland.svg)](https://www.jetbrains.com/)