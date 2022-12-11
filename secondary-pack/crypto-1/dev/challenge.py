#!/usr/bin/env python3

import os


FLAG = os.getenv('FLAG', 'Cup{************************}').encode()
assert len(FLAG) == 29 and FLAG.startswith(b'Cup{') and FLAG.endswith(b'}')


def encrypt(plaintext: bytes, key: bytes) -> bytes:
    data = [i - byte for i, byte in enumerate(plaintext)]

    for i, byte in enumerate(data):
        idx = 2 * i + 1
        byte = sum(data) - idx * byte
        data[i] = (byte + idx * key[0]) % 256

    return bytes(data)


def main():
    key = os.urandom(len(FLAG))
    ciphertext = encrypt(FLAG, key)

    print(ciphertext.hex())


if __name__ == '__main__':
    main()
