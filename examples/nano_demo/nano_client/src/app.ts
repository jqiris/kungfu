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
        }
    }, (data) => {
        console.log(data);
    });
});
