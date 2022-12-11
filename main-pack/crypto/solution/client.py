#!/usr/bin/env python3

import socket


class Client:
    def __init__(self, hostname, port):
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.sock.settimeout(3)
        self.sock.connect((hostname, port))
        self.file = self.sock.makefile('rwb')

        for _ in range(2):
            self.file.readline()

    def __del__(self):
        self.sock.close()

    def _sendcmd(self, cmd, arg = ''):
        line = cmd + ' ' + str(arg)
        self.file.read(2)
        self.file.write(line.encode() + b'\n')
        self.file.flush()

    def _readint(self):
        return int(self.file.readline().decode().strip())

    def generate(self):
        self._sendcmd('generate')

    def key(self):
        self._sendcmd('key')
        e = self._readint()
        n = self._readint()
        return (e, n)

    def hint(self):
        self._sendcmd('hint')
        return self._readint()

    def flag(self):
        self._sendcmd('flag')
        return self._readint()

    def encrypt(self, n):
        self._sendcmd('encrypt', n)
        return self._readint()

    def super(self, n):
        self._sendcmd('super', n)
        return self._readint()

    def exit(self):
        self._sendcmd('exit')


if __name__ == '__main__':
    client = Client('0.0.0.0', 17171)

    try:
        client.generate()
        e, n = client.key()

        plaintext = 12345
        ciphertext = client.encrypt(plaintext)
        assert pow(plaintext, e, n) == ciphertext
    finally:
        client.exit()
