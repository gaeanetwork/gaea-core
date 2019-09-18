// https://nodejs.org/api/crypto.html
const crypto = require('crypto');
const algorithm = 'aes-256-cbc';
const keyLength = 32
const ivLength = 16

function encrypt(plaintext, aesKey) {
    let key = paddingAESKey(aesKey);
    let ivBuffer = Buffer.alloc(ivLength).fill(key);

    const cipher = crypto.createCipheriv(algorithm, key, ivBuffer);
    let encrypted = cipher.update(plaintext, 'utf8', 'hex');
    encrypted += cipher.final('hex');
    return encrypted
}

function paddingAESKey(aesKey) {
    let padding = Buffer.alloc(keyLength);
    
    // 32
    padding.write(aesKey);
    return padding.toString()
}

module.exports = {
    paddingAESKey, encrypt
};
