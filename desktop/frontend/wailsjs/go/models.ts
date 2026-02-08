export namespace main {
	
	export class DesktopConfig {
	    APIUrl: string;
	    Token: string;
	
	    static createFrom(source: any = {}) {
	        return new DesktopConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.APIUrl = source["APIUrl"];
	        this.Token = source["Token"];
	    }
	}

}

