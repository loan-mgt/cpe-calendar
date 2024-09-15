/*
The unwrapped signing key.
*/
let encryptionKey;

/*
Convert a string into an ArrayBuffer
from https://developers.google.com/web/updates/2012/06/How-to-convert-ArrayBuffer-to-and-from-String
*/
function str2ab(str) {
    const buf = new ArrayBuffer(str.length);
    const bufView = new Uint8Array(buf);
    for (let i = 0, strLen = str.length; i < strLen; i++) {
        bufView[i] = str.charCodeAt(i);
    }
    return buf;
}

/*
Import a PEM encoded RSA public key, to use for RSA-OAEP encryption.
Takes a string containing the PEM encoded key, and returns a Promise
that will resolve to a CryptoKey representing the public key.
*/


function importPublicKey(pem) {
    // fetch the part of the PEM string between header and footer
    const pemHeader = "-----BEGIN PUBLIC KEY-----";
    const pemFooter = "-----END PUBLIC KEY-----";
    const pemContents = pem.substring(pemHeader.length, pem.length - pemFooter.length);

    // Remove newlines and any spaces from the base64-encoded content
    const pemCleaned = pemContents.replace(/\r?\n|\r/g, "");

    // base64 decode the string to get the binary data
    const binaryDerString = window.atob(pemCleaned);

    // convert from a binary string to an ArrayBuffer
    const binaryDer = str2ab(binaryDerString);

    return window.crypto.subtle.importKey(
        "spki",
        binaryDer,
        {
            name: "RSA-OAEP",
            hash: "SHA-256"
        },
        true,
        ["encrypt"]
    );
}

/*
Fetch the contents of the "message" textbox, and encode it
in a form we can use for the encrypt operation.
*/
function getMessageEncoding(message) {
    const enc = new TextEncoder();
    return enc.encode(message);
}

/*
Get the encoded message, encrypt it and display a representation
of the ciphertext in the "Ciphertext" element.
*/
function arrayBufferToBase64(buffer) {
    let binary = '';
    const bytes = new Uint8Array(buffer);
    const len = bytes.byteLength;
    for (let i = 0; i < len; i++) {
        binary += String.fromCharCode(bytes[i]);
    }
    return window.btoa(binary);
}

/*
  Get the encoded message, encrypt it, and return the Base64 encoded ciphertext.
*/
async function encryptMessage(message) {

    const encoded = getMessageEncoding(message);
    const ciphertext = await window.crypto.subtle.encrypt(
        {
            name: "RSA-OAEP"
        },
        encryptionKey,
        encoded
    );

    // Convert the ciphertext to Base64
    return arrayBufferToBase64(ciphertext);
}
