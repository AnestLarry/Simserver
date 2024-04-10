interface LSItem {
    Name: string,
    ModTime: number,
    Size: bigint,
}
interface FileListPanelConfig {
    fileList: LSItem[],
    folderList: LSItem[],
    urlStack: string[],
    listenIndex: number,
    fresh: (f: boolean) => void,
}