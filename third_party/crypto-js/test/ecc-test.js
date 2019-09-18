const EC = require('elliptic').ec;
const CryptoJS = require("crypto-js");

QUnit.test("hello test", function (assert) {
    // Create and initialize EC context
    // (better do it once and reuse it)
    var ec = new EC('p256');

    // Generate keys
    var key = ec.genKeyPair();
    var priv = key.getPrivate().toString(16)
    var pub = key.getPublic().encode('hex')
    console.log('priv: ', priv)
    console.log('pub: ', pub)

    // Recover key pair
    var key1 = ec.keyFromPrivate(priv)
    var key2 = ec.keyFromPublic(pub, 'hex');
    
    assert.ok(key.getPublic().toString() == key1.getPublic().toString(), "Recover Private Key Passed!");
    assert.ok(key.getPublic().toString() == key2.getPublic().toString(), "Recover Public Key Passed!");
    
    // Private Key Sign
    priv = "30435376894abc8771edb4452be55307cb95e7e61e4709dcb37f31b3e3156d60"
    pub = "0494c5036edb28489f1fba70b9a855aea31e41678a9ee6d5148d222df845e179ecf40371da34e7df3bedbf86a68bd652e35d1e1f33de498bdae657b90a755fcc22"
    var key = ec.keyFromPrivate(priv)
    var msgHash = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
    var signature = key.sign(msgHash);
    
    // Public Key Verify
    key = ec.keyFromPublic(pub, 'hex');
    assert.ok(key.verify(msgHash, signature) == true, "Sign and Verify");

    // // Sign the message's hash (input must be an array, or a hex-string)
    // var msgHash = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
    // var signature = key.sign(msgHash);
    // console.log("signature.r: ", signature.r.toString())
    // console.log("signature.s: ", signature.s.toString())

    // // Export DER encoded signature in Array
    // var derSign = signature.toDER();
    // console.log("signature: ", derSign)

    // // Verify signature
    // console.log(key.verify(msgHash, derSign));

    // // CHECK WITH NO PRIVATE KEY

    // var pubPoint = key.getPublic();
    // var x = pubPoint.getX();
    // var y = pubPoint.getY();

    // // Public Key MUST be either:
    // // 1) '04' + hex string of x + hex string of y; or
    // // 2) object with two hex string properties (x and y); or
    // // 3) object with two buffer properties (x and y)
    // var pub = pubPoint.encode('hex');                                 // case 1
    // var pub = { x: x.toString('hex'), y: y.toString('hex') };         // case 2
    // var pub = { x: x.toBuffer(), y: y.toBuffer() };                   // case 3
    // var pub = { x: x.toArrayLike(Buffer), y: y.toArrayLike(Buffer) }; // case 3

    // // Import public key
    // var key = ec.keyFromPublic(pub, 'hex');

    // // Signature MUST be either:
    // // 1) DER-encoded signature as hex-string; or
    // // 2) DER-encoded signature as buffer; or
    // // 3) object with two hex-string properties (r and s); or
    // // 4) object with two buffer properties (r and s)

    // var signature = '3046022100...'; // case 1
    // var signature = new Buffer('...'); // case 2
    // var signature = { r: 'b1fc...', s: '9c42...' }; // case 3

    // console.log("pub: ", key.getPublic().encode('hex'))
    // console.log("msgHash: ", msgHash)
    // // Verify signature
    // console.log(key.verify(msgHash, signature));
});

QUnit.test("hello test", function (assert) {
    var ec = new EC('p256');

    // Generate keys
    // var key1 = ec.genKeyPair();
    var key1 = ec.keyFromPrivate('d233a716bf371afc597636a9b00342603759ab9f39ab5954e6d51a996cd2bfdd')
    var serverPub = ec.keyFromPublic('048f03f8321b00a4466f4bf4be51c91898cd50d8cc64c6ecf53e73443e348d5925a16f88c8952b78ebac2dc277a2cc54c77b4c3c07830f49629b689edf63086293', 'hex')

    var shared1 = key1.derive(serverPub.getPublic());
    // var shared2 = serverPub.getPublic().mul(key1.getPrivate()).getX()

    // console.log('Both shared secrets are BN instances');
    // console.log('priv: ',key1.getPrivate().toString(16))

    // console.log('pub: ',key1.getPublic().encode('hex'))
    // console.log('pub.X: ',key1.getPublic().getX().toString())
    // console.log('pub.Y: ',key1.getPublic().getY().toString())
    console.log("shared1:", shared1.toString(10));
    assert.ok(shared1.toString(10) == "99629494961789099446600506511571181974131020151128428048581066925321839516601", "Passed!");
    // console.log("shared2:", shared2.toString(10));
    var message = '99629494961789099446600506511571181974131020151128428048581066925321839516601'
    // cd3ae50c26656fa7c927c84c1dcbb736cd73c77d2b5b11f1d20b268bd3249fa5
    console.log("sha256:", CryptoJS.SHA256(message).toString());
    assert.ok(1 == "1", "Passed!");
});