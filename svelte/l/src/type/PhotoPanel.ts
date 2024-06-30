interface PhotoPanelConfig {
    imgList: LSItem[],
    listenIndex: number,
    showIndex: number,
    drawerHidden: boolean,
    showDrawer: () => void,
    fresh: (f: boolean) => void,
    setShowIndex: (i: number) => void,
    getListGroup: () => any,
}