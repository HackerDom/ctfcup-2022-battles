#!/usr/bin/env python3

def main():
    output = 'b2ea3c32a5e06d57be3a1f050bfba7692574623d67be0bcaea3bb9118b'
    ciphertext = bytes.fromhex(output)

    for key_byte in range(256):
        for initial_sum in range(256):
            data = []

            for i in range(len(ciphertext)):
                idx = 2 * i + 1
                idx_inv = pow(idx, -1, 256)

                byte = idx_inv * (initial_sum - (ciphertext[i] - key_byte)) % 256
                data.append(byte)

                initial_sum = initial_sum - byte + (initial_sum - (idx * byte))

            plaintext = []

            for i, byte in enumerate(data):
                plaintext.append((i - byte) % 256)

            flag = bytes(plaintext)

            if flag.startswith(b'Cup{') and flag.endswith(b'}'):
                print(flag)
                return


if __name__ == '__main__':
    main()
