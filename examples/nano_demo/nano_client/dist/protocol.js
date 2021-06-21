"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    Object.defineProperty(o, k2, { enumerable: true, get: function() { return m[k]; } });
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.strdecode = exports.strencode = exports.Message = exports.MessageType = exports.Package = exports.PacketType = exports.Pomelo = void 0;
/* eslint-disable no-param-reassign */
/* eslint-disable max-params */
var net = __importStar(require("net"));
var PKG_HEAD_BYTES = 4;
var MSG_FLAG_BYTES = 1;
var MSG_ROUTE_CODE_BYTES = 2;
var MSG_ROUTE_LEN_BYTES = 1;
var MSG_ROUTE_CODE_MAX = 0xffff;
var MSG_COMPRESS_ROUTE_MASK = 0x1;
var MSG_COMPRESS_GZIP_MASK = 0x1;
var MSG_COMPRESS_GZIP_ENCODE_MASK = 1 << 4;
var MSG_TYPE_MASK = 0x7;
var Pomelo = /** @class */ (function () {
    function Pomelo() {
    }
    Pomelo.prototype.init = function (params, cb) {
        var _this = this;
        this.params = params;
        var host = params.host;
        var port = params.port;
        var url = host + ":" + port;
        this.socket = new net.Socket();
        this.socket.connect(port, host, function () {
            console.log('connect to ' + url);
            if (cb) {
                cb(_this.socket);
            }
        });
        this.socket.on('reconnect', function () {
            console.log('reconnect');
        });
        this.socket.on('data', function (data) {
            console.log("data:", data);
        });
        this.socket.on('error', function (error) {
            console.log("err:", error);
        });
        this.socket.on("disconnect", function () {
            _this.socket.close();
            _this.socket = null;
            console.log("socket disconnect");
        });
    };
    return Pomelo;
}());
exports.Pomelo = Pomelo;
// 包类型
var PacketType;
(function (PacketType) {
    PacketType[PacketType["Type_None"] = 0] = "Type_None";
    PacketType[PacketType["TYPE_HANDSHAKE"] = 1] = "TYPE_HANDSHAKE";
    PacketType[PacketType["TYPE_HANDSHAKE_ACK"] = 2] = "TYPE_HANDSHAKE_ACK";
    PacketType[PacketType["TYPE_HEARTBEAT"] = 3] = "TYPE_HEARTBEAT";
    PacketType[PacketType["TYPE_DATA"] = 4] = "TYPE_DATA";
    PacketType[PacketType["TYPE_KICK"] = 5] = "TYPE_KICK";
})(PacketType = exports.PacketType || (exports.PacketType = {}));
var Package = /** @class */ (function () {
    function Package() {
    }
    Package.prototype.isValidType = function (type) {
        return type >= PacketType.TYPE_HANDSHAKE && type <= PacketType.TYPE_KICK;
    };
    Package.prototype.encode = function (type, body) {
        var length = body ? body.length : 0;
        var buffer = Buffer.alloc(PKG_HEAD_BYTES + length);
        var index = 0;
        buffer[index++] = type & 0xff;
        buffer[index++] = (length >> 16) & 0xff;
        buffer[index++] = (length >> 8) & 0xff;
        buffer[index++] = length & 0xff;
        if (body) {
            copyArray(buffer, index, body, 0, length);
        }
        return buffer;
    };
    Package.prototype.decode = function (buffer) {
        var offset = 0;
        var bytes = Buffer.from(buffer);
        var length = 0;
        var rs = [];
        while (offset < bytes.length) {
            var type = bytes[offset++];
            length = ((bytes[offset++]) << 16 | (bytes[offset++]) << 8 | bytes[offset++]) >>> 0;
            if (!this.isValidType(type) || length > bytes.length) {
                return { 'type': type }; // return invalid type, then disconnect!
            }
            var body = length ? Buffer.alloc(length) : null;
            if (body) {
                copyArray(body, 0, bytes, offset, length);
            }
            offset += length;
            rs.push({ 'type': type, 'body': body });
        }
        return rs.length === 1 ? rs[0] : rs;
    };
    return Package;
}());
exports.Package = Package;
// 消息类型
var MessageType;
(function (MessageType) {
    MessageType[MessageType["TYPE_REQUEST"] = 0] = "TYPE_REQUEST";
    MessageType[MessageType["TYPE_NOTIFY"] = 1] = "TYPE_NOTIFY";
    MessageType[MessageType["TYPE_RESPONSE"] = 2] = "TYPE_RESPONSE";
    MessageType[MessageType["TYPE_PUSH"] = 3] = "TYPE_PUSH";
})(MessageType = exports.MessageType || (exports.MessageType = {}));
var Message = /** @class */ (function () {
    function Message() {
    }
    Message.prototype.encode = function (id, type, compressRoute, route, msg, compressGzip) {
        // caculate message max length
        var idBytes = msgHasId(type) ? caculateMsgIdBytes(id) : 0;
        var msgLen = MSG_FLAG_BYTES + idBytes;
        if (msgHasRoute(type)) {
            if (compressRoute) {
                if (typeof route !== 'number') {
                    throw new Error('error flag for number route!');
                }
                msgLen += MSG_ROUTE_CODE_BYTES;
            }
            else {
                msgLen += MSG_ROUTE_LEN_BYTES;
                if (route) {
                    route = exports.strencode(route);
                    if (route.length > 255) {
                        throw new Error('route maxlength is overflow');
                    }
                    msgLen += route.length;
                }
            }
        }
        if (msg) {
            msgLen += msg.length;
        }
        var buffer = Buffer.alloc(msgLen);
        var offset = 0;
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
    };
    Message.prototype.decode = function (buffer) {
        var bytes = Buffer.from(buffer);
        var bytesLen = bytes.length || bytes.byteLength;
        var offset = 0;
        var id = 0;
        var route = null;
        // parse flag
        var flag = bytes[offset++];
        var compressRoute = flag & MSG_COMPRESS_ROUTE_MASK;
        var type = (flag >> 1) & MSG_TYPE_MASK;
        var compressGzip = (flag >> 4) & MSG_COMPRESS_GZIP_MASK;
        // parse id
        if (msgHasId(type)) {
            var m = 0;
            var i = 0;
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
            }
            else {
                var routeLen = bytes[offset++];
                if (routeLen) {
                    route = Buffer.alloc(routeLen);
                    copyArray(route, 0, bytes, offset, routeLen);
                    route = exports.strdecode(route);
                }
                else {
                    route = '';
                }
                offset += routeLen;
            }
        }
        // parse body
        var bodyLen = bytesLen - offset;
        var body = Buffer.alloc(bodyLen);
        copyArray(body, 0, bytes, offset, bodyLen);
        return {
            'id': id, 'type': type, 'compressRoute': compressRoute,
            'route': route, 'body': body, 'compressGzip': compressGzip
        };
    };
    return Message;
}());
exports.Message = Message;
var strencode = function (str) {
    return Buffer.from(str);
};
exports.strencode = strencode;
var strdecode = function (buffer) {
    return buffer.toString();
};
exports.strdecode = strdecode;
var copyArray = function (dest, doffset, src, soffset, length) {
    if (typeof src.copy === 'function') {
        // Buffer
        src.copy(dest, doffset, soffset, soffset + length);
    }
    else {
        // Uint8Array
        for (var index = 0; index < length; index++) {
            dest[doffset++] = src[soffset++];
        }
    }
};
var msgHasId = function (type) {
    return type === MessageType.TYPE_REQUEST || type === MessageType.TYPE_RESPONSE;
};
var msgHasRoute = function (type) {
    return type === MessageType.TYPE_REQUEST || type === MessageType.TYPE_NOTIFY ||
        type === MessageType.TYPE_PUSH;
};
var caculateMsgIdBytes = function (id) {
    var len = 0;
    do {
        len += 1;
        id >>= 7;
    } while (id > 0);
    return len;
};
var encodeMsgFlag = function (type, compressRoute, buffer, offset, compressGzip) {
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
var encodeMsgId = function (id, buffer, offset) {
    do {
        var tmp = id % 128;
        var next = Math.floor(id / 128);
        if (next !== 0) {
            tmp = tmp + 128;
        }
        buffer[offset++] = tmp;
        id = next;
    } while (id !== 0);
    return offset;
};
var encodeMsgRoute = function (compressRoute, _route, buffer, offset) {
    if (compressRoute) {
        var route = _route;
        if (route > MSG_ROUTE_CODE_MAX) {
            throw new Error('route number is overflow');
        }
        buffer[offset++] = (route >> 8) & 0xff;
        buffer[offset++] = route & 0xff;
    }
    else {
        var route = _route;
        if (route) {
            buffer[offset++] = route.length & 0xff;
            copyArray(buffer, offset, route, 0, route.length);
            offset += route.length;
        }
        else {
            buffer[offset++] = 0;
        }
    }
    return offset;
};
var encodeMsgBody = function (msg, buffer, offset) {
    copyArray(buffer, offset, msg, 0, msg.length);
    return offset + msg.length;
};
//# sourceMappingURL=protocol.js.map