export namespace clents {
	
	export class Vout {
	    value: number;
	
	    static createFrom(source: any = {}) {
	        return new Vout(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.value = source["value"];
	    }
	}
	export class RespSendToAddress {
	    GasUsed: number;
	    vot: Vout[];
	
	    static createFrom(source: any = {}) {
	        return new RespSendToAddress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.GasUsed = source["GasUsed"];
	        this.vot = this.convertValues(source["vot"], Vout);
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

export namespace peer {
	
	export class AddressInfo {
	    index: number;
	    address: string;
	    balance: number;
	    balance_frozen: number;
	    balance_lockup: number;
	    peer_type: number;
	
	    static createFrom(source: any = {}) {
	        return new AddressInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.address = source["address"];
	        this.balance = source["balance"];
	        this.balance_frozen = source["balance_frozen"];
	        this.balance_lockup = source["balance_lockup"];
	        this.peer_type = source["peer_type"];
	    }
	}
	export class Peer {
	    id: string;
	    group: string;
	    name: string;
	    ip: string;
	    port: number;
	    username: string;
	    password: string;
	    status: boolean;
	    highest_block: number;
	    current_block: number;
	    snapshot_height: number;
	    total_balance: number;
	    types: {[key: number]: number};
	    used_time: number;
	    is_del: boolean;
	    error: string;
	    updated_at: number;
	    default_address?: AddressInfo;
	    addresses: AddressInfo[];
	
	    static createFrom(source: any = {}) {
	        return new Peer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.group = source["group"];
	        this.name = source["name"];
	        this.ip = source["ip"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.status = source["status"];
	        this.highest_block = source["highest_block"];
	        this.current_block = source["current_block"];
	        this.snapshot_height = source["snapshot_height"];
	        this.total_balance = source["total_balance"];
	        this.types = source["types"];
	        this.used_time = source["used_time"];
	        this.is_del = source["is_del"];
	        this.error = source["error"];
	        this.updated_at = source["updated_at"];
	        this.default_address = this.convertValues(source["default_address"], AddressInfo);
	        this.addresses = this.convertValues(source["addresses"], AddressInfo);
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

