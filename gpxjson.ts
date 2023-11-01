
export class GPXTrackSegment {
    points?: GPXPoint[];
    ext?: ExtensionNode[];
}
export class GPXTrack {
    name?: string;
    cmt?: string;
    desc?: string;
    src?: string;
    number?: number;
    type?: string;
    segments: GPXTrackSegment[];
    ext?: ExtensionNode[];
}
export class NullableInt {

}
export class GPXRoute {
    name?: string;
    cmt?: string;
    desc?: string;
    src?: string;
    number?: NullableInt;
    type?: string;
    pts?: GPXPoint[];
    ext?: ExtensionNode[];
}
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
    ele?: number|undefined;
    time?: string|undefined;
    magvar?: string;
    geoidheight?: string;
    name?: string;
    cmt?: string;
    desc?: string;
    src?: string;
    sym?: string;
    type?: string;
    fix?: string;
    sat?: number;
    hdop?: number;
    vdop?: number;
    pdop?: number;
    ageofdgpsdata?: number;
    dgpsid?: number;
    ext?: ExtensionNode[];
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
    desc?: string;
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
    rte?: GPXRoute[];
    trk: GPXTrack[];
    ext?: ExtensionNode[];
    metadataExt?: ExtensionNode[];
}