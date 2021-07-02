import * as CryptoJS from "./crypto-js";

class AesCbc {
    private key;
    private iv;
    public constructor(key, iv) {
        this.key = CryptoJS.enc.Utf8.parse(key);
        this.iv = CryptoJS.enc.Utf8.parse(iv);
    }
    public Encrypt(word) {
        let srcs = CryptoJS.enc.Utf8.parse(word);
        // 加密模式为CBC，补码方式为PKCS5Padding（也就是PKCS7）
        let encrypted = CryptoJS.AES.encrypt(srcs, this.key, {
            iv: this.iv,
            mode: CryptoJS.mode.CBC,
            padding: CryptoJS.pad.Pkcs7
        });
        // 返回base64
        return CryptoJS.enc.Base64.stringify(encrypted.ciphertext);

    }

    public Decrypt(word) {
        let base64 = CryptoJS.enc.Base64.parse(word);

        let src = CryptoJS.enc.Base64.stringify(base64);

        // 解密模式为CBC，补码方式为PKCS5Padding（也就是PKCS7）
        let decrypt = CryptoJS.AES.decrypt(src, this.key, {
            iv: this.iv,
            mode: CryptoJS.mode.CBC,
            padding: CryptoJS.pad.Pkcs7
        });

        let decryptedStr = decrypt.toString(CryptoJS.enc.Utf8);
        return decryptedStr.toString();
    }
}



let k = 'MPbcthxqyCT0pr1Z';
console.log('加密的key:', k);
let v = 'KUdkFuunmQ0hndvH';
console.log('加密的iv:', v);
let msg = "hellonodejsaesd128-cbc加密和解密";
console.log("需要加密的数据:", msg);

let aescbc = new AesCbc(k, v);
let a = aescbc.Encrypt(msg);
console.log("加密后的数据为:", a);

let b = aescbc.Decrypt(a);
console.log("解密后的数据为:", b);

