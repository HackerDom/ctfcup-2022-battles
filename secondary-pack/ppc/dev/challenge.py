#!/usr/bin/env python3

import collections


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


def main():
    value = func(1337, 10 ** 1000, 13 ** 37)
    print('Cup{' + str(value) + '}')


if __name__ == '__main__':
    main()
