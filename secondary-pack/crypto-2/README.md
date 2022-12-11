# CTFCup 2022 | Battles / Secondary | Crypto-2

## Название

> not-rsa

## Описание

> Мы пытались реализовать RSA, но что-то пошло не так. Сможете разобраться?

## Раздатка

Участникам нужно выдать файлы:

* [public/not-rsa.tar.gz](public/not-rsa.tar.gz)

## Деплой

Не требуется.

## Решение

В решении реализовано [шифрование Эль-Гамаля](https://en.wikipedia.org/wiki/ElGamal_encryption). Участникам выдаётся публичный ключ криптосистемы и шифртекст.

В реализации алгоритма нет ошибок, но есть уязвимость в генерации случайного непредсказуемого числа в процессе шифрования (nonce):

```python
def encrypt(plaintext: Plaintext, public_key: PublicKey) -> Ciphertext:
    nonce = time.process_time_ns()
    shared = pow(public_key.h, nonce, public_key.p)

    c1 = pow(public_key.g, nonce, public_key.p)
    c2 = plaintext.m * shared % public_key.p

    return Ciphertext(c1, c2)
```

Число nonce должно быть большим (порядка p), чтобы нельзя было посчитать дискретный логарифм. Но в задании в качестве nonce используется значение функции `time.process_time_ns()`, которое не превышает `1_000_000_000_000` (грубая оценка сверху). Можно воспользоваться [алгоритмом больших и малых шагов](https://en.wikipedia.org/wiki/Baby-step_giant-step), чтобы быстро найти nonce:

```
nonce = discrete_log(g, c2, p)
```

После нахождения nonce можно применить следующую атаку:

```
shared = h ^ nonce (mod p)
shared_inv = 1 / shared (mod p)
m = c2 * shared_inv (mod p)
```

Осталось перевести m из числа в строку и получить флаг.

**Пример решения**: [solution/solver.py](solution/solver.py)

## Флаг

```
Cup{small3st_n0nce_in_the_w0rld}
```
