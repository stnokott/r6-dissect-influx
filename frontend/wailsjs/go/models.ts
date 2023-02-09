export namespace config {
	
	export class InfluxConfigJson {
	    host: string;
	    port: number;
	    org: string;
	    bucket: string;
	    token: string;
	
	    static createFrom(source: any = {}) {
	        return new InfluxConfigJson(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.host = source["host"];
	        this.port = source["port"];
	        this.org = source["org"];
	        this.bucket = source["bucket"];
	        this.token = source["token"];
	    }
	}
	export class GameConfigJson {
	    install_dir: string;
	
	    static createFrom(source: any = {}) {
	        return new GameConfigJson(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.install_dir = source["install_dir"];
	    }
	}
	export class Config {
	    game: GameConfigJson;
	    influx_db: InfluxConfigJson;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.game = this.convertValues(source["game"], GameConfigJson);
	        this.influx_db = this.convertValues(source["influx_db"], InfluxConfigJson);
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
	

}

