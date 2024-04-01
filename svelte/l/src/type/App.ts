interface Panel {
    baseUrl: string;
    workUrl: string;
    pageMode: string;
    sortedBy: string;
    photo: Photo;
    workUrlListening: Array<(f: boolean) => void>;
    pushUrlStack: (x: string) => void;
    popUrlStack: () => void;
}
interface Photo {
    photoMode: string[];
    sizeRange: SizeRange;
}
interface SizeRange {
    isApply: boolean;
    widthRange: number;
    heightRange: number;
    updateRange: () => void;
}