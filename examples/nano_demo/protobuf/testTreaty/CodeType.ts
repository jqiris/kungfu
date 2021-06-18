export enum CodeType {
    CodeSuccess, //成功
    CodeFailed, //失败
    CodeChooseBackendLogin, //选择后端服务器进行登录
    CodeLoginReconnect,//进行重连
    CodeCannotFindBackend,//找不到后端服务器
    CodeUndefinedDealMsg,//未定义处理消息
    CodeNotLogin,//未登录
    CodeNotRightConnector, //请登录绑定connector
    CodeNotLoginBackend,//没有登录后端服务器
}