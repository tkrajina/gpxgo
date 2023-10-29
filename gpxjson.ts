
export class ExtensionNodeAttr {
    ns?: string;
    name?: string;
    val?: string;
}
export class ExtensionNode {
    ns?: string;
    name?: string;
    attrs?: ExtensionNodeAttr[];
    data?: string;
    nodes?: ExtensionNode[];
}
export class GPXPoint {
    lat: number;
    lon: number;
    ele: number|undefined;
}
export class GPXTrackSegment {
    points: GPXPoint[];
}
export class GPXTrack {
    segments: GPXTrackSegment[];
}
export class GPX {
    tracks: GPXTrack[];
    ext: ExtensionNode[];
}