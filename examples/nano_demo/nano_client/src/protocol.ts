/* eslint-disable no-param-reassign */
/* eslint-disable max-params */
import * as net from "net"
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

export class Pomelo {
    private params;
    private socket;
    public init(params, cb) {
        this.params = params
        let host = params.host
        let port = params.port
        let url = host + ":" + port

        this.socket = new net.Socket();
        this.socket.connect(port, host, () => {
            console.log('connect to ' + url);
            if (cb) {
                cb(this.socket);
            }
        })
        this.socket.on('reconnect', () => {
            console.log('reconnect');
        });
        this.socket.on('data', (data) => {
            console.log("data:", data)
        })
        this.socket.on('error', (error) => {
            console.log("err:", error)
        })
        this.socket.on("disconnect", () => {
            this.socket.close();
            this.socket = null;
            console.log("socket disconnect");
        })

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
