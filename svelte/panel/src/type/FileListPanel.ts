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
    filterCond: string,
    showFolders: boolean,
    showFiles: boolean,
    columnCount: string,
    fileListFiltered: () => LSItem[],
    folderListFiltered: () => LSItem[],
    fresh: (f: boolean) => void,
}