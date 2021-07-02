import * as crypto from 'crypto';

/**
 * 加密方法
 * @param key 加密key
 * @param iv       向量
 * @param data     需要加密的数据
 * @returns string
 */
let encrypt = function (key, iv, data) {
    let cipher = crypto.createCipheriv('aes-128-cbc', key, iv);
    let crypted = cipher.update(data, 'utf8', 'binary');
    crypted += cipher.final('binary');
    crypted = Buffer.from(crypted, 'binary').toString('base64');
    return crypted;
};

/**
 * 解密方法
 * @param key      解密的key
 * @param iv       向量
 * @param crypted  密文
 * @returns string
 */
let decrypt = function (key, iv, crypted) {
    let data = Buffer.from(crypted, 'base64').toString('binary');
    let decipher = crypto.createDecipheriv('aes-128-cbc', key, iv);
    let decoded = decipher.update(data, 'binary', 'utf8');
    decoded += decipher.final('utf8');
    return decoded;
};


let k = 'SwGqJKXUgpDfKza8';
console.log('加密的key:', k.toString());
let v = '8u8LP27JDTYtX2Wl';
console.log('加密的iv:', v);
let msg = "hellonodejsaesd128-cbc加密和解密";
console.log("需要加密的数据:", msg);
let crypted = encrypt(k, v, msg);
console.log("数据加密后:", crypted);
let dec = decrypt(k, v, crypted);
console.log("数据解密后:", dec);
