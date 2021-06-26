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
var Protocol = __importStar(require("./protocol"));
var path_1 = __importDefault(require("path"));
// let Package = Protocol.Package;
// let Message = Protocol.Message;
var pomelo = new Protocol.Pomelo();
pomelo.init({
    host: "127.0.0.1",
    port: "8288",
    useProtos: true,
    protoPath: path_1.default.join(__dirname, "../src/treaty.proto") // test
}, function () {
    console.log("cb callback");
    pomelo.request("UserConnector.Login", {
        uid: 1001,
        nickname: "jason",
        token: "ce0da27df7150196625e48c843deb1f9",
        backend: {
            server_id: "backend_3001",
            server_type: "backend",
            server_name: "kungfu backend",
            server_ip: "127.0.0.1",
            client_port: 8388
        },
        connector: {
            server_id: "connector_2001",
            server_type: "connector",
            server_name: "kungfu connector",
            server_ip: "127.0.0.1",
            client_port: 8288
        }
    }, function (data) {
        var testInt = data.test_int;
        console.log("login resp:", data, testInt);
        if (data.code === 0) {
            // 登录成功
            pomelo.request("UserConnector.ChannelMsg", {
                uid: 1001,
                msg_data: "hello chat"
            }, function (resp) {
                console.log("channel msg resp:", resp);
            });
        }
        else {
            console.log("登录失败:", data.msg);
        }
    });
});
//# sourceMappingURL=app.js.map