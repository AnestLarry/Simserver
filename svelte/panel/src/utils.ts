export const http = {
  Get: (theUrl: string) => {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("GET", theUrl, false); // false for synchronous request
    xmlHttp.send(null);
    return xmlHttp.responseText;
  }
}
export const client = {
  ImgList: (url: string, prefix: string) => {
    const imgExts = ["png", "jpg", "jpeg", "webp", "gif", "bmp", "avif", "heif", "svg"];
    try {
      var response = http.Get(url);
      var data: LSResponse = JSON.parse(response);
      if (data["fileList"] === null) {
        return [];
      }
      var resList: LSItem[] = [];
      data["fileList"].sort(client.sortFunction("NameLenOrder"));
      data["fileList"].filter(x => {
        return imgExts.filter(e => x.Name.toLowerCase().endsWith(e)).length > 0 ? x : null;
      }).forEach(x => {
        var i = x;
        i.Name = prefix + i.Name;
        resList.push(i);
      });
      return resList;
    } catch (error) {
      console.warn("ImgList API call failed (ls feature may be disabled):", error);
      return [];
    }
  },
  FolderList: (url: string) => {
    try {
      var response = http.Get(url);
      var data: LSResponse = JSON.parse(response);
      if (data["folderList"] === null) {
        return [];
      }
      data["folderList"].sort(client.sortFunction("NameOrder"));
      return data["folderList"];
    } catch (error) {
      console.warn("FolderList API call failed (ls feature may be disabled):", error);
      return [];
    }
  },
  sortFunction: (x: string) => {

    var res = {
      NameOrder: (a: LSItem, b: LSItem) => a.Name.localeCompare(b.Name),
      NameReverse: (a: LSItem, b: LSItem) => -a.Name.localeCompare(b.Name),
      NameLenOrder: (a: LSItem, b: LSItem) => a.Name.length == b.Name.length ? a.Name.localeCompare(b.Name) : (a.Name.length - b.Name.length),
      NameLenReverse: (a: LSItem, b: LSItem) => -(a.Name.length == b.Name.length ? a.Name.localeCompare(b.Name) : (a.Name.length - b.Name.length)),
      TimeOrder: (a: LSItem, b: LSItem) => Number(a.ModTime - b.ModTime),
      TimeReverse: (a: LSItem, b: LSItem) => Number(b.ModTime - a.ModTime),
    }[x];
    if (res != undefined) {
      return res;
    } else {
      throw `sortFunction: "x" has not match with "${x}"`;
    }
  },
};