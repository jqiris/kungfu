"use strict";
exports.__esModule = true;
var CryptoJS = require("./crypto-js");
var AesCbc = /** @class */ (function () {
    function AesCbc(key, iv) {
        this.key = CryptoJS.enc.Utf8.parse(key);
        this.iv = CryptoJS.enc.Utf8.parse(iv);
    }
    AesCbc.prototype.Encrypt = function (word) {
        var srcs = CryptoJS.enc.Utf8.parse(word);
        // 加密模式为CBC，补码方式为PKCS5Padding（也就是PKCS7）
        var encrypted = CryptoJS.AES.encrypt(srcs, this.key, {
            iv: this.iv,
            mode: CryptoJS.mode.CBC,
            padding: CryptoJS.pad.Pkcs7
        });
        // 返回base64
        return CryptoJS.enc.Base64.stringify(encrypted.ciphertext);
    };
    AesCbc.prototype.Decrypt = function (word) {
        var base64 = CryptoJS.enc.Base64.parse(word);
        var src = CryptoJS.enc.Base64.stringify(base64);
        // 解密模式为CBC，补码方式为PKCS5Padding（也就是PKCS7）
        var decrypt = CryptoJS.AES.decrypt(src, this.key, {
            iv: this.iv,
            mode: CryptoJS.mode.CBC,
            padding: CryptoJS.pad.Pkcs7
        });
        var decryptedStr = decrypt.toString(CryptoJS.enc.Utf8);
        return decryptedStr.toString();
    };
    return AesCbc;
}());
var k = 'MPbcthxqyCT0pr1Z';
console.log('加密的key:', k);
var v = 'KUdkFuunmQ0hndvH';
console.log('加密的iv:', v);
var msg = "hellonodejsaesd128-cbc加密和解密";
console.log("需要加密的数据:", msg);
var aescbc = new AesCbc(k, v);
var a = aescbc.Encrypt(msg);
console.log("加密后的数据为:", a);
var b = aescbc.Decrypt(a);
console.log("解密后的数据为:", b);
