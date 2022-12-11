#!/usr/bin/env python3

import os
import time
import random
import collections
from typing import Tuple

from Crypto.Util.number import isPrime  # pycryptodome


PublicKey = collections.namedtuple('PublicKey', ['g', 'p', 'h'])
PrivateKey = collections.namedtuple('PrivateKey', ['g', 'p', 'x'])

Plaintext = collections.namedtuple('Plaintext', ['m'])
Ciphertext = collections.namedtuple('Ciphertext', ['c1', 'c2'])


FLAG = os.getenv('FLAG', 'Cup{***************************}').encode()
assert len(FLAG) == 32


def generate_key(bits: int) -> Tuple[PublicKey, PrivateKey]:
    while True:
        q = random.getrandbits(bits - 1)
        p = 2 * q + 1

        if isPrime(p):
            break

    for g in range(2, p - 1):
        if pow(g, q, p) != 1:
            break

    x = random.randrange(2, q)
    h = pow(g, x, p)

    return PublicKey(g, p, h), PrivateKey(g, p, x)


def encrypt(plaintext: Plaintext, public_key: PublicKey) -> Ciphertext:
    nonce = time.process_time_ns()
    shared = pow(public_key.h, nonce, public_key.p)

    c1 = pow(public_key.g, nonce, public_key.p)
    c2 = plaintext.m * shared % public_key.p

    return Ciphertext(c1, c2)


def decrypt(ciphertext: Ciphertext, private_key: PrivateKey) -> Plaintext:
    shared = pow(ciphertext.c1, private_key.x, private_key.p)
    shared_inv = pow(shared, -1, private_key.p)

    m = ciphertext.c2 * shared_inv % private_key.p

    return Plaintext(m)


def bytes_to_int(data: bytes) -> int:
    return int.from_bytes(data, 'big')


def int_to_bytes(number: int) -> bytes:
    length = (number.bit_length() + 7) // 8

    return number.to_bytes(length, 'big')


def main():
    bits = 1024
    public_key, private_key = generate_key(bits)

    message = bytes_to_int(FLAG)
    plaintext = Plaintext(message)

    ciphertext = encrypt(plaintext, public_key)
    assert decrypt(ciphertext, private_key) == plaintext

    print(public_key)
    print(ciphertext)


if __name__ == '__main__':
    main()
