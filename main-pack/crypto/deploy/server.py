#!/usr/bin/env python3

import os
import cmd
import sys
import random
import secrets
import collections

import gmpy2


FLAG = os.getenv('FLAG', 'Cup{*********************************************}').encode()
assert len(FLAG) == 50


RSA = collections.namedtuple('RSA', ['e', 'n', 'p', 'q'])


class Dummy:
    def __getattribute__(self, _):
        raise Exception('key is not generated')


class Challenge(cmd.Cmd):
    intro = 'Welcome to babyrsa-2022 challenge. Type help or ? to list commands.\n'
    prompt = '> '

    key = Dummy()

    def onecmd(self, line):
        line = line.lower()

        if line == 'eof':
            sys.exit(0)

        try:
            return super().onecmd(line)
        except Exception as e:
            return print(f'ERROR: {e}')

    def do_generate(self, _):
        'Generate new RSA key:  GENERATE'
        self.key = self._generate_key(512)

    def do_key(self, _):
        'Print the RSA public key:  KEY'
        print(self.key.e)
        print(self.key.n)

    def do_hint(self, _):
        'Print a hint for RSA:  HINT'
        hint = self._get_hint()
        print(hint)

    def do_flag(self, _):
        'Print the encrypted flag:  FLAG'
        secret = secrets.token_bytes(8)
        plaintext = self._to_number(FLAG + b'\x00' + secret)
        ciphertext = self._encrypt(plaintext)
        print(ciphertext)

    def do_encrypt(self, arg):
        'Encrypt a number using RSA key:  ENCRYPT 12345'
        plaintext = int(arg)
        ciphertext = self._encrypt(plaintext)
        print(ciphertext)

    def do_super(self, arg):
        'Super encrypt a number using RSA key:  SUPER 12345'
        plaintext = int(arg)
        ciphertext = self._super_encrypt(plaintext)
        print(ciphertext)

    def do_exit(self, _):
        'Exit from challenge:  EXIT'
        sys.exit(0)

    def _random_prime(self, bits):
        rnd = random.getrandbits(bits)
        return int(gmpy2.next_prime(rnd))

    def _generate_key(self, bits):
        e = 65537
        p = self._random_prime(bits)
        q = self._random_prime(bits)
        n = p * q
        return RSA(e, n, p, q)

    def _get_hint(self):
        return self.key.p - self._to_number(FLAG)

    def _encrypt(self, plaintext):
        return pow(plaintext, self.key.e, self.key.n)

    def _super_encrypt(self, plaintext):
        exponent = self._to_number(FLAG)
        return pow(plaintext, exponent, self.key.n)

    def _to_number(self, data):
        return int.from_bytes(data, 'big')


def main():
    challenge = Challenge()
    challenge.cmdloop()


if __name__ == '__main__':
    main()
