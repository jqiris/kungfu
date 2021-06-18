import { Server } from './Server';
import { CodeType } from './CodeType';
export interface UserConnector_Login_Req {
    uid: number,
    nickname: string,
    token: string,
    backend: Server
}

export interface UserConnector_Login_Res {
    code: CodeType,
    msg: string,
    backend: Server
}