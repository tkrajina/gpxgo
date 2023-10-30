
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
export class GPXTrackSegment {
    points: GPXPoint[];
}
export class GPXTrack {
    segments: GPXTrackSegment[];
}
export class GPXRoute {

}
export class GPXPoint {
    lat: number;
    lon: number;
    ele: number|undefined;
    ts: string|undefined;
}
export class GPXAttributes {
    nsattrs?: {[key: string]: };
}
export class GPX {
    xmlns?: string;
    xmlnsxsi?: string;
    xmlschemaloc?: string;
    attrs?: GPXAttributes;
    version?: string;
    creator?: string;
    name?: string;
    description?: string;
    authorname?: string;
    authoremail?: string;
    authorlink?: string;
    authorlinktext?: string;
    authorlinktype?: string;
    copyright?: string;
    copyrightyear?: string;
    copyrightlicense?: string;
    link?: string;
    linktext?: string;
    linktype?: string;
    time?: string|undefined;
    keywords?: string;
    waypoints?: GPXPoint[];
    routes?: GPXRoute[];
    tracks?: GPXTrack[];
    ext?: ExtensionNode[];
    metadataExt?: ExtensionNode[];
}