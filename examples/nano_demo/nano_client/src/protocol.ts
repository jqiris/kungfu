/* eslint-disable guard-for-in */
/* eslint-disable @typescript-eslint/prefer-optional-chain */
/* eslint-disable no-param-reassign */
/* eslint-disable max-params */
import * as net from "net"
import protobuf from "protocol-buffers";
import fs from "fs";
import mitt from 'mitt';
declare let Buffer: any;

let PKG_HEAD_BYTES = 4;
let MSG_FLAG_BYTES = 1;
let MSG_ROUTE_CODE_BYTES = 2;
let MSG_ROUTE_LEN_BYTES = 1;

let MSG_ROUTE_CODE_MAX = 0xffff;

let MSG_COMPRESS_ROUTE_MASK = 0x1;
let MSG_COMPRESS_GZIP_MASK = 0x1;
let MSG_COMPRESS_GZIP_ENCODE_MASK = 1 << 4;
let MSG_TYPE_MASK = 0x7;

let RES_OK = 200;
let RES_FAIL = 500;
let RES_OLD_CLIENT = 501;

let JS_CLIENT_TYPE = 'js-socket';
let JS_CLIENT_VERSION = '0.0.1';
export class Pomelo {
    private params;
    private socket;
    private reqId = 1;
    private routeMap = {};
    private callbacks = {};
    private packet = new Package();
    private message = new Message();
    private protos = null;
    private data: any;
    private dict = null;
    private handlers = {}
    private emitter = mitt();
    private heartbeatInterval = 0;
    private heartbeatTimeout = 0;
    private nextHeartbeatTimeout = 0;
    private gapThreshold = 100; // heartbeat gap threashold
    private heartbeatId = null;
    private heartbeatTimeoutId = null;
    private handshakeBuffer = {
        'sys': {
            type: JS_CLIENT_TYPE,
            version: JS_CLIENT_VERSION
        },
        'user': {}
    };
    private handshakeCallback = null;
    private initCallback = null;

    public init(params, cb) {
        this.initCallback = cb;
        this.handshakeBuffer.user = params.user;
        this.handshakeCallback = params.handshakeCallback;
        this.emitter.on("*", (t: string, e) => { this["on" + t](e) });
        this.params = params
        if (params.useProtos && params.protoPath.length > 0) {
            this.protos = protobuf(fs.readFileSync(params.protoPath))
        }
        let host = params.host
        let port = params.port
        let url = host + ":" + port

        this.socket = new net.Socket();
        this.socket.setNoDelay(true);
        this.socket.connect(port, host)
        this.socket.on('connect', () => {
            console.log('connect to url:', url);
        });
        this.socket.on('ready', () => {
            console.log('connect ready');
            let obj = this.packet.encode(PacketType.TYPE_HANDSHAKE, strencode(JSON.stringify(this.handshakeBuffer)));
            this.send(obj);
        });
        this.socket.on('data', (data) => {
            this.processPackage(this.packet.decode(data));
            // new package arrived, update the heartbeat timeout
            if (this.heartbeatTimeout) {
                this.nextHeartbeatTimeout = Date.now() + this.heartbeatTimeout;
            }
        })
        this.socket.on('error', (error) => {
            this.emitter.emit('Error', error)
        })
        this.socket.on("close", () => {
            this.disconnect();
        })
        // 注册事件
        let handshake = (data) => {
            data = JSON.parse(strdecode(data));
            if (data.code === RES_OLD_CLIENT) {
                this.emitter.emit('Error', 'client version not fullfill');
                return;
            }

            if (data.code !== RES_OK) {
                this.emitter.emit('Error', 'handshake fail');
                return;
            }
            console.log('handshake callback:', data)
            this.handshakeInit(data);

            let obj = this.packet.encode(PacketType.TYPE_HANDSHAKE_ACK);
            this.send(obj);
            if (this.initCallback) {
                this.initCallback();
                this.initCallback = null;
            }
        }

        let onData = (data) => {
            let msg = this.message.decode(data);

            if (msg.id > 0) {
                msg.route = this.routeMap[msg.id];
                delete this.routeMap[msg.id];
                if (!msg.route) {
                    return;
                }
            }

            msg.body = this.deCompose(msg);

            this.processMessage(msg);
        }
        let onKick = (data) => {
            console.log("onKick:", data)
        }

        let heartbeat = (data) => {
            console.log("heartbeat:", data)
            if (!this.heartbeatInterval) {
                // no heartbeat
                return;
            }

            let obj = this.packet.encode(PacketType.TYPE_HEARTBEAT);
            if (this.heartbeatTimeoutId) {
                clearTimeout(this.heartbeatTimeoutId);
                this.heartbeatTimeoutId = null;
            }

            if (this.heartbeatId) {
                // already in a heartbeat interval
                return;
            }

            this.heartbeatId = setTimeout(() => {
                this.heartbeatId = null;
                this.send(obj);

                this.nextHeartbeatTimeout = Date.now() + this.heartbeatTimeout;
                this.heartbeatTimeoutId = setTimeout(heartbeatTimeoutCb, this.heartbeatTimeout);
            }, this.heartbeatInterval);
        };

        let heartbeatTimeoutCb = () => {
            let gap = this.nextHeartbeatTimeout - Date.now();
            if (gap > this.gapThreshold) {
                this.heartbeatTimeoutId = setTimeout(heartbeatTimeoutCb, gap);
            } else {
                console.error('server heartbeat timeout');
                this.emitter.emit("Error", 'heartbeat timeout');
                this.disconnect();
            }
        };
        this.handlers[PacketType.TYPE_HANDSHAKE] = handshake;
        this.handlers[PacketType.TYPE_HEARTBEAT] = heartbeat;
        this.handlers[PacketType.TYPE_DATA] = onData;
        this.handlers[PacketType.TYPE_KICK] = onKick;

    }
    public disconnect() {
        if (this.socket) {
            this.socket.destroy();
            this.socket = null;
            console.log("socket close");
        }
        if (this.heartbeatId) {
            clearTimeout(this.heartbeatId);
            this.heartbeatId = null;
        }
        if (this.heartbeatTimeoutId) {
            clearTimeout(this.heartbeatTimeoutId);
            this.heartbeatTimeoutId = null;
        }
    }


    public processPackage(msg) {
        this.handlers[msg.type](msg.body);
    }
    public processMessage(msg) {
        if (!msg.id) {
            // server push message
            this.emitter.emit(msg.route, msg.body);
            return;
        }

        // if have a id then find the callback function with the request
        let cb = this.callbacks[msg.id];

        delete this.callbacks[msg.id];
        if (typeof cb !== 'function') {
            return;
        }

        cb(msg.body);
    }
    public onError(data) {
        console.log("received error:", data);
    }
    public deCompose(msg) {
        let protos = this.data.protos ? this.data.protos.server : {};
        let abbrs = this.data.abbrs;
        let route = msg.route;

        if (msg.compressRoute) {
            if (!abbrs[route]) {
                return {};
            }

            route = msg.route = abbrs[route];
        }
        let proto = protos[route]
        if (proto && this.protos[proto]) {
            return this.protos[proto].decode(msg.body);
        } else {
            return JSON.parse(strdecode(msg.body));
        }
    };


    public handshakeInit(data) {
        if (data.sys && data.sys.heartbeat) {
            this.heartbeatInterval = data.sys.heartbeat * 1000; // heartbeat interval
            this.heartbeatTimeout = this.heartbeatInterval * 2; // max heartbeat timeout
        } else {
            this.heartbeatInterval = 0;
            this.heartbeatTimeout = 0;
        }

        this.initData(data);

        if (typeof this.handshakeCallback === 'function') {
            this.handshakeCallback(data.user);
        }
    };
    public initData(data) {
        if (!data || !data.sys) {
            return;
        }
        this.data = this.data || {};
        let dict = data.sys.dict;
        let protos = data.sys.protos;

        // Init compress dict
        if (dict) {
            this.data.dict = dict;
            this.data.abbrs = {};

            for (let route in dict) {
                this.data.abbrs[dict[route]] = route;
            }
        }

        // Init protobuf protos
        if (protos) {
            this.data.protos = {
                server: protos.server || {},
                client: protos.client || {}
            };
        }
    }
    public request(route, msg, cb) {
        if (!route) {
            return;
        }
        this.reqId++;
        this.sendMessage(this.reqId, route, msg);
        this.callbacks[this.reqId] = cb;
        this.routeMap[this.reqId] = route;
    }

    private sendMessage(reqId, route, msg) {
        let type = reqId ? MessageType.TYPE_REQUEST : MessageType.TYPE_NOTIFY;

        // compress message by protobuf
        let protos = this.data.protos ? this.data.protos.client : {};
        let proto = protos[route]
        if (proto && this.protos[proto]) {
            msg = this.protos[proto].encode(msg);
        } else {
            msg = strencode(JSON.stringify(msg));
        }


        let compressRoute = false;
        if (this.dict?.[route]) {
            route = this.dict[route];
            compressRoute = true;
        }

        msg = this.message.encode(reqId, type, compressRoute, route, msg);
        let packet = this.packet.encode(PacketType.TYPE_DATA, msg);
        this.send(packet);
    }

    private send(packet) {
        this.socket.write(packet);
    }
}

// 包类型
export enum PacketType {
    Type_None,
    TYPE_HANDSHAKE,
    TYPE_HANDSHAKE_ACK,
    TYPE_HEARTBEAT,
    TYPE_DATA,
    TYPE_KICK
}



export class Package {
    public isValidType(type: PacketType): boolean {
        return type >= PacketType.TYPE_HANDSHAKE && type <= PacketType.TYPE_KICK
    }

    public encode(type: PacketType, body?: Buffer) {
        let length = body ? body.length : 0;
        let buffer = Buffer.alloc(PKG_HEAD_BYTES + length);
        let index = 0;
        buffer[index++] = type & 0xff;
        buffer[index++] = (length >> 16) & 0xff;
        buffer[index++] = (length >> 8) & 0xff;
        buffer[index++] = length & 0xff;
        if (body) {
            copyArray(buffer, index, body, 0, length);
        }
        return buffer;
    }

    public decode(buffer: Buffer) {
        let offset = 0;
        let bytes = Buffer.from(buffer);
        let length = 0;
        let rs = [];
        while (offset < bytes.length) {
            let type = bytes[offset++];
            length = ((bytes[offset++]) << 16 | (bytes[offset++]) << 8 | bytes[offset++]) >>> 0;
            if (!this.isValidType(type) || length > bytes.length) {
                return { 'type': type }; // return invalid type, then disconnect!
            }
            let body = length ? Buffer.alloc(length) : null;
            if (body) {
                copyArray(body, 0, bytes, offset, length);
            }
            offset += length;
            rs.push({ 'type': type, 'body': body });
        }
        return rs.length === 1 ? rs[0] : rs;
    }
}


// 消息类型
export enum MessageType {
    TYPE_REQUEST,
    TYPE_NOTIFY,
    TYPE_RESPONSE,
    TYPE_PUSH
}

export class Message {
    public encode(id: number, type: MessageType, compressRoute: boolean, route: number | string | Buffer, msg: Buffer, compressGzip?: boolean) {
        // caculate message max length
        let idBytes = msgHasId(type) ? caculateMsgIdBytes(id) : 0;
        let msgLen = MSG_FLAG_BYTES + idBytes;

        if (msgHasRoute(type)) {
            if (compressRoute) {
                if (typeof route !== 'number') {
                    throw new Error('error flag for number route!');
                }
                msgLen += MSG_ROUTE_CODE_BYTES;
            } else {
                msgLen += MSG_ROUTE_LEN_BYTES;
                if (route) {
                    route = strencode(route as string);
                    if ((route as string).length > 255) {
                        throw new Error('route maxlength is overflow');
                    }
                    msgLen += (route as string).length;
                }
            }
        }

        if (msg) {
            msgLen += msg.length;
        }

        let buffer = Buffer.alloc(msgLen);
        let offset = 0;

        // add flag
        offset = encodeMsgFlag(type, compressRoute, buffer, offset, compressGzip);

        // add message id
        if (msgHasId(type)) {
            offset = encodeMsgId(id, buffer, offset);
        }

        // add route
        if (msgHasRoute(type)) {
            offset = encodeMsgRoute(compressRoute, route, buffer, offset);
        }

        // add body
        if (msg) {
            offset = encodeMsgBody(msg, buffer, offset);
        }

        return buffer;
    }

    public decode(buffer: Buffer) {
        let bytes = Buffer.from(buffer);
        let bytesLen = bytes.length || bytes.byteLength;
        let offset = 0;
        let id = 0;
        let route = null;

        // parse flag
        let flag = bytes[offset++];
        let compressRoute = flag & MSG_COMPRESS_ROUTE_MASK;
        let type = (flag >> 1) & MSG_TYPE_MASK;
        let compressGzip = (flag >> 4) & MSG_COMPRESS_GZIP_MASK;

        // parse id
        if (msgHasId(type)) {
            let m = 0;
            let i = 0;
            do {
                m = parseInt(bytes[offset], 10);
                id += (m & 0x7f) << (7 * i);
                offset++;
                i++;
            } while (m >= 128);
        }

        // parse route
        if (msgHasRoute(type)) {
            if (compressRoute) {
                route = (bytes[offset++]) << 8 | bytes[offset++];
            } else {
                let routeLen = bytes[offset++];
                if (routeLen) {
                    route = Buffer.alloc(routeLen);
                    copyArray(route, 0, bytes, offset, routeLen);
                    route = strdecode(route);
                } else {
                    route = '';
                }
                offset += routeLen;
            }
        }

        // parse body
        let bodyLen = bytesLen - offset;
        let body = Buffer.alloc(bodyLen);

        copyArray(body, 0, bytes, offset, bodyLen);

        return {
            'id': id, 'type': type, 'compressRoute': compressRoute,
            'route': route, 'body': body, 'compressGzip': compressGzip
        };
    }
}

export let strencode = function (str: string) {
    return Buffer.from(str)
}
export let strdecode = function (buffer: object) {
    return buffer.toString();
}

let copyArray = function (dest: Buffer, doffset: number, src: Buffer, soffset: number, length: number) {
    if (typeof src.copy === 'function') {
        // Buffer
        src.copy(dest, doffset, soffset, soffset + length);
    } else {
        // Uint8Array
        for (let index = 0; index < length; index++) {
            dest[doffset++] = src[soffset++];
        }
    }
};

let msgHasId = function (type: MessageType) {
    return type === MessageType.TYPE_REQUEST || type === MessageType.TYPE_RESPONSE;
};

let msgHasRoute = function (type: MessageType) {
    return type === MessageType.TYPE_REQUEST || type === MessageType.TYPE_NOTIFY ||
        type === MessageType.TYPE_PUSH;
};

let caculateMsgIdBytes = function (id: number) {
    let len = 0;
    do {
        len += 1;
        id >>= 7;
    } while (id > 0);
    return len;
};

let encodeMsgFlag = function (type: number, compressRoute: boolean, buffer: Buffer, offset: number, compressGzip: boolean) {
    if (type !== MessageType.TYPE_REQUEST && type !== MessageType.TYPE_NOTIFY &&
        type !== MessageType.TYPE_RESPONSE && type !== MessageType.TYPE_PUSH) {
        throw new Error('unkonw message type: ' + type);
    }

    buffer[offset] = (type << 1) | (compressRoute ? 1 : 0);

    if (compressGzip) {
        buffer[offset] = buffer[offset] | MSG_COMPRESS_GZIP_ENCODE_MASK;
    }

    return offset + MSG_FLAG_BYTES;
};

let encodeMsgId = function (id: number, buffer: Buffer, offset: number) {
    do {
        let tmp = id % 128;
        let next = Math.floor(id / 128);

        if (next !== 0) {
            tmp = tmp + 128;
        }
        buffer[offset++] = tmp;

        id = next;
    } while (id !== 0);

    return offset;
};

let encodeMsgRoute = function (compressRoute: boolean, _route: number | string | Buffer, buffer: Buffer, offset: number) {
    if (compressRoute) {
        let route = _route as number;
        if (route > MSG_ROUTE_CODE_MAX) {
            throw new Error('route number is overflow');
        }

        buffer[offset++] = (route >> 8) & 0xff;
        buffer[offset++] = route & 0xff;
    } else {
        let route = _route as Buffer;
        if (route) {
            buffer[offset++] = route.length & 0xff;
            copyArray(buffer, offset, route as Buffer, 0, route.length);
            offset += route.length;
        } else {
            buffer[offset++] = 0;
        }
    }

    return offset;
};

let encodeMsgBody = function (msg: Buffer, buffer: Buffer, offset: number) {
    copyArray(buffer, offset, msg, 0, msg.length);
    return offset + msg.length;
};
