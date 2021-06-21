import * as Protocol from './protocol';

// let Package = Protocol.Package;
// let Message = Protocol.Message;
let pomelo = new Protocol.Pomelo();



pomelo.init({
    host: "127.0.0.1",
    port: "8288"
}, (res) => {

});
