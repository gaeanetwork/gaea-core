const aes = require("../aes/crypto");
const CryptoJS = require("crypto-js");

// paddingAESKey
QUnit.test("paddingAESKey - Key length", function (assert) {
    var key = "Adf";
    assert.ok(aes.paddingAESKey(key).length == 32, "Passed! Key length is less than 32");

    var key = "AdfAdfAdfAdfAdfAAdfAdfAdfAdfAdfA";
    assert.ok(aes.paddingAESKey(key).length == 32, "Passed! Key length is equal to 32");

    var key = "AdfAdfAdfAdfAdfAAdfAdfAdfAdfAdfAAdfAdfAdfAdfAdfAAdfAdfAdfAdfAdfA";
    assert.ok(aes.paddingAESKey(key).length == 32, "Passed! Key length is greater than 32");
});

QUnit.test("encrypt - crypto", function (assert) {
    let key = 'AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', data = 'dmyz.org';
    assert.ok(aes.encrypt(data, key) == "1527e127f824de210938f54bc8555cd9", "Passed!");
    assert.ok(CryptoJS.AES.encrypt(data, CryptoJS.enc.Utf8.parse(key), {
        iv: CryptoJS.enc.Utf8.parse(key),
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7
    }).ciphertext.toString() == "1527e127f824de210938f54bc8555cd9", "Passed!");
    
    // Key - smaller length
    key = 'hello';
    assert.ok(aes.encrypt(data, key) == "fdc9769dc539780fdd6d5d3b3f11151a", "Passed!");
    key = aes.paddingAESKey(key)
    assert.ok(CryptoJS.AES.encrypt(data, CryptoJS.enc.Utf8.parse(key), {
        iv: CryptoJS.enc.Utf8.parse(key),
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7
    }).ciphertext.toString() == "fdc9769dc539780fdd6d5d3b3f11151a", "Passed!");

    // Key - larger length
    key = 'helloAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA';
    assert.ok(aes.encrypt(data, key) == "6ece0eab57eb212f4a25941f4c4c9cb3", "Passed!");
    key = aes.paddingAESKey(key)
    assert.ok(CryptoJS.AES.encrypt(data, CryptoJS.enc.Utf8.parse(key), {
        iv: CryptoJS.enc.Utf8.parse(key),
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7
    }).ciphertext.toString() == "6ece0eab57eb212f4a25941f4c4c9cb3", "Passed!");
    
    // data - larger length
    data = 'helloAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAdmyz.org';
    assert.ok(aes.encrypt(data, key) == "241a6575abd0be96ce6598bd3053ad9e7a8937c5c703797ab30ff5e2cb51e006759e5b2543c280155ebb59c5ab085b87", "Passed!");
    assert.ok(CryptoJS.AES.encrypt(data, CryptoJS.enc.Utf8.parse(key), {
        iv: CryptoJS.enc.Utf8.parse(key),
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7
    }).ciphertext.toString() == "241a6575abd0be96ce6598bd3053ad9e7a8937c5c703797ab30ff5e2cb51e006759e5b2543c280155ebb59c5ab085b87", "Passed!");
});
