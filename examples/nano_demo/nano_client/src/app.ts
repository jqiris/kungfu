import * as Protocol from './protocol';
import path from "path";
// let Package = Protocol.Package;
// let Message = Protocol.Message;
let pomelo = new Protocol.Pomelo();



pomelo.init({
    host: "127.0.0.1",
    port: "8288",
    useProtos: true,
    protoPath: path.join(__dirname, "../src/treaty.proto") // test
}, () => {
    console.log("cb callback")
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
        connector:{
            server_id: "connector_2001",
            server_type: "connector",
            server_name: "kungfu connector",
            server_ip: "127.0.0.1",
            client_port: 8288
          }
    }, (data) => {
        let testInt: number = data.test_int
        console.log("login resp:", data, testInt)
        if (data.code === 0) {
            // 登录成功
            pomelo.request("UserConnector.ChannelMsg", {
                uid: 1001,
                msg_data: "hello chat"
            }, (resp) => {
                console.log("channel msg resp:", resp)
            })

        } else {
            console.log("登录失败:", data.msg);
        }
    });
});
