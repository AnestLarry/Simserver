interface PhotoPanelConfig {
    imgList: LSItem[];
    listenIndex: number;
    showIndex: number;
    drawerHidden: boolean;
    fresh: (f: boolean) => void;
    setShowIndex: (i: number) => void;
    getListGroup:()=>any;
}