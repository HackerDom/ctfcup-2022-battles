#!/usr/bin/env python3

import random
import collections

import numpy as np


def func(init: int, k: int, mod: int) -> int:
    queue = collections.deque(
        [3, 1, 3, 3, 7, init],
    )

    for _ in range(k):
        value = sum(
            x * y for x, y in enumerate(queue)
        ) % mod

        queue.append(value)
        queue.popleft()

    return queue.pop()


def fast_func(init: int, k: int, mod: int) -> int:
    vec = np.matrix([3, 1, 3, 3, 7, init])
    mat = np.matrix([
        [0, 0, 0, 0, 0, 0],
        [1, 0, 0, 0, 0, 1],
        [0, 1, 0, 0, 0, 2],
        [0, 0, 1, 0, 0, 3],
        [0, 0, 0, 1, 0, 4],
        [0, 0, 0, 0, 1, 5],
    ], dtype = object)

    pow = np.identity(6, dtype = object)
    cur = mat

    while k > 0:
        if k & 1 == 1:
            pow = (pow * cur) % mod

        cur = (cur * cur) % mod
        k >>= 1

    res = vec * pow % mod

    return res.tolist().pop().pop()


def main():
    for _ in range(100):
        init = random.randint(1, 10_000)
        k = random.randint(1, 10_000)
        mod = random.randint(1, 1_000_000)

        assert func(init, k, mod) == fast_func(init, k, mod)

    value = fast_func(1337, 10 ** 1000, 13 ** 37)
    print('Cup{' + str(value) + '}')


if __name__ == '__main__':
    main()
