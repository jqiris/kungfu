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
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.strdecode = exports.strencode = exports.Message = exports.MessageType = exports.Package = exports.PacketType = exports.Pomelo = void 0;
/* eslint-disable guard-for-in */
/* eslint-disable @typescript-eslint/prefer-optional-chain */
/* eslint-disable no-param-reassign */
/* eslint-disable max-params */
var net = __importStar(require("net"));
var protocol_buffers_1 = __importDefault(require("protocol-buffers"));
var fs_1 = __importDefault(require("fs"));
var mitt_1 = __importDefault(require("mitt"));
var PKG_HEAD_BYTES = 4;
var MSG_FLAG_BYTES = 1;
var MSG_ROUTE_CODE_BYTES = 2;
var MSG_ROUTE_LEN_BYTES = 1;
var MSG_ROUTE_CODE_MAX = 0xffff;
var MSG_COMPRESS_ROUTE_MASK = 0x1;
var MSG_COMPRESS_GZIP_MASK = 0x1;
var MSG_COMPRESS_GZIP_ENCODE_MASK = 1 << 4;
var MSG_TYPE_MASK = 0x7;
var RES_OK = 200;
var RES_FAIL = 500;
var RES_OLD_CLIENT = 501;
var JS_CLIENT_TYPE = 'js-socket';
var JS_CLIENT_VERSION = '0.0.1';
var Pomelo = /** @class */ (function () {
    function Pomelo() {
        this.reqId = 1;
        this.routeMap = {};
        this.callbacks = {};
        this.packet = new Package();
        this.message = new Message();
        this.protos = null;
        this.dict = null;
        this.handlers = {};
        this.emitter = mitt_1.default();
        this.heartbeatInterval = 0;
        this.heartbeatTimeout = 0;
        this.nextHeartbeatTimeout = 0;
        this.gapThreshold = 100; // heartbeat gap threashold
        this.heartbeatId = null;
        this.heartbeatTimeoutId = null;
        this.handshakeBuffer = {
            'sys': {
                type: JS_CLIENT_TYPE,
                version: JS_CLIENT_VERSION
            },
            'user': {}
        };
        this.handshakeCallback = null;
        this.initCallback = null;
    }
    Pomelo.prototype.init = function (params, cb) {
        var _this = this;
        this.initCallback = cb;
        this.handshakeBuffer.user = params.user;
        this.handshakeCallback = params.handshakeCallback;
        this.emitter.on("*", function (t, e) { _this["on" + t](e); });
        this.params = params;
        if (params.useProtos && params.protoPath.length > 0) {
            this.protos = protocol_buffers_1.default(fs_1.default.readFileSync(params.protoPath));
        }
        var host = params.host;
        var port = params.port;
        var url = host + ":" + port;
        this.socket = new net.Socket();
        this.socket.setNoDelay(true);
        this.socket.connect(port, host);
        this.socket.on('connect', function () {
            console.log('connect to url:', url);
        });
        this.socket.on('ready', function () {
            console.log('connect ready');
            var obj = _this.packet.encode(PacketType.TYPE_HANDSHAKE, exports.strencode(JSON.stringify(_this.handshakeBuffer)));
            _this.send(obj);
        });
        this.socket.on('data', function (data) {
            // console.log("socket data:", data);
            _this.processPackage(_this.packet.decode(data));
            // new package arrived, update the heartbeat timeout
            if (_this.heartbeatTimeout) {
                _this.nextHeartbeatTimeout = Date.now() + _this.heartbeatTimeout;
            }
        });
        this.socket.on('error', function (error) {
            _this.emitter.emit('Error', error);
        });
        this.socket.on("close", function () {
            _this.disconnect();
        });
        // 注册事件
        var handshake = function (data) {
            data = JSON.parse(exports.strdecode(data));
            if (data.code === RES_OLD_CLIENT) {
                _this.emitter.emit('Error', 'client version not fullfill');
                return;
            }
            if (data.code !== RES_OK) {
                _this.emitter.emit('Error', 'handshake fail');
                return;
            }
            console.log('handshake callback:', data);
            _this.handshakeInit(data);
            var obj = _this.packet.encode(PacketType.TYPE_HANDSHAKE_ACK);
            _this.send(obj);
            if (_this.initCallback) {
                _this.initCallback();
                _this.initCallback = null;
            }
        };
        var onData = function (data) {
            var msg = _this.message.decode(data);
            if (msg.id > 0) {
                msg.route = _this.routeMap[msg.id];
                delete _this.routeMap[msg.id];
                if (!msg.route) {
                    return;
                }
            }
            msg.body = _this.deCompose(msg);
            _this.processMessage(msg);
        };
        var onKick = function (data) {
            console.log("onKick:", data);
        };
        var heartbeat = function () {
            console.log("heartbeat");
            if (!_this.heartbeatInterval) {
                // no heartbeat
                return;
            }
            var obj = _this.packet.encode(PacketType.TYPE_HEARTBEAT);
            if (_this.heartbeatTimeoutId) {
                clearTimeout(_this.heartbeatTimeoutId);
                _this.heartbeatTimeoutId = null;
            }
            if (_this.heartbeatId) {
                // already in a heartbeat interval
                return;
            }
            _this.heartbeatId = setTimeout(function () {
                _this.heartbeatId = null;
                _this.send(obj);
                _this.nextHeartbeatTimeout = Date.now() + _this.heartbeatTimeout;
                _this.heartbeatTimeoutId = setTimeout(heartbeatTimeoutCb, _this.heartbeatTimeout);
            }, _this.heartbeatInterval);
        };
        var heartbeatTimeoutCb = function () {
            var gap = _this.nextHeartbeatTimeout - Date.now();
            if (gap > _this.gapThreshold) {
                _this.heartbeatTimeoutId = setTimeout(heartbeatTimeoutCb, gap);
            }
            else {
                console.error('server heartbeat timeout');
                _this.emitter.emit("Error", 'heartbeat timeout');
                _this.disconnect();
            }
        };
        this.handlers[PacketType.TYPE_HANDSHAKE] = handshake;
        this.handlers[PacketType.TYPE_HEARTBEAT] = heartbeat;
        this.handlers[PacketType.TYPE_DATA] = onData;
        this.handlers[PacketType.TYPE_KICK] = onKick;
    };
    Pomelo.prototype.disconnect = function () {
        console.log('disconnect');
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
    };
    Pomelo.prototype.processPackage = function (msg) {
        this.handlers[msg.type](msg.body);
    };
    Pomelo.prototype.processMessage = function (msg) {
        if (!msg.id) {
            // server push message
            this.emitter.emit(msg.route, msg.body);
            return;
        }
        // if have a id then find the callback function with the request
        var cb = this.callbacks[msg.id];
        delete this.callbacks[msg.id];
        if (typeof cb !== 'function') {
            return;
        }
        cb(msg.body);
    };
    Pomelo.prototype.onError = function (data) {
        console.log("received error:", data);
    };
    Pomelo.prototype.deCompose = function (msg) {
        var protos = this.data.protos ? this.data.protos.server : {};
        var abbrs = this.data.abbrs;
        var route = msg.route;
        if (msg.compressRoute) {
            if (!abbrs[route]) {
                return {};
            }
            route = msg.route = abbrs[route];
        }
        var proto = protos[route];
        if (proto && this.protos[proto]) {
            return this.protos[proto].decode(msg.body);
        }
        else {
            return JSON.parse(exports.strdecode(msg.body));
        }
    };
    ;
    Pomelo.prototype.handshakeInit = function (data) {
        if (data.sys && data.sys.heartbeat) {
            this.heartbeatInterval = data.sys.heartbeat * 1000; // heartbeat interval
            this.heartbeatTimeout = this.heartbeatInterval * 2; // max heartbeat timeout
        }
        else {
            this.heartbeatInterval = 0;
            this.heartbeatTimeout = 0;
        }
        this.initData(data);
        if (typeof this.handshakeCallback === 'function') {
            this.handshakeCallback(data.user);
        }
    };
    ;
    Pomelo.prototype.initData = function (data) {
        if (!data || !data.sys) {
            return;
        }
        this.data = this.data || {};
        var dict = data.sys.dict;
        var protos = data.sys.protos;
        // Init compress dict
        if (dict) {
            this.data.dict = dict;
            this.data.abbrs = {};
            for (var route in dict) {
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
    };
    Pomelo.prototype.request = function (route, msg, cb) {
        if (!route) {
            return;
        }
        this.reqId++;
        this.sendMessage(this.reqId, route, msg);
        this.callbacks[this.reqId] = cb;
        this.routeMap[this.reqId] = route;
    };
    Pomelo.prototype.sendMessage = function (reqId, route, msg) {
        var _a;
        var type = reqId ? MessageType.TYPE_REQUEST : MessageType.TYPE_NOTIFY;
        // compress message by protobuf
        var protos = this.data.protos ? this.data.protos.client : {};
        var proto = protos[route];
        if (proto && this.protos[proto]) {
            msg = this.protos[proto].encode(msg);
        }
        else {
            msg = exports.strencode(JSON.stringify(msg));
        }
        var compressRoute = false;
        if ((_a = this.dict) === null || _a === void 0 ? void 0 : _a[route]) {
            route = this.dict[route];
            compressRoute = true;
        }
        msg = this.message.encode(reqId, type, compressRoute, route, msg);
        var packet = this.packet.encode(PacketType.TYPE_DATA, msg);
        this.send(packet);
    };
    Pomelo.prototype.send = function (packet) {
        this.socket.write(packet);
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