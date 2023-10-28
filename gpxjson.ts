
export class NullableFloat {


    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}
export class GPXPoint {
    lat: number;
    lon: number;
    ele: NullableFloat;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.lat = source["lat"];
        this.lon = source["lon"];
        this.ele = this.convertValues(source["ele"], NullableFloat);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class GPXTrackSegment {
    points: GPXPoint[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.points = this.convertValues(source["points"], GPXPoint);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class GPXTrack {
    segments: GPXTrackSegment[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.segments = this.convertValues(source["segments"], GPXTrackSegment);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class GPX {
    tracks: GPXTrack[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.tracks = this.convertValues(source["tracks"], GPXTrack);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}