# CTFCup 2022 | Battles / Secondary | Crypto-1

## Название

> simple-cipher

## Описание

> Взлом этого незамысловатого шифра вряд ли займёт много времени.

## Раздатка

Участникам нужно выдать файлы:

* [public/simple-cipher.tar.gz](public/simple-cipher.tar.gz)

## Деплой

Не требуется.

## Решение

В задании реализован простой самописный шифр:

```python
def encrypt(plaintext: bytes, key: bytes) -> bytes:
    data = [i - byte for i, byte in enumerate(plaintext)]

    for i, byte in enumerate(data):
        idx = 2 * i + 1
        byte = sum(data) - idx * byte
        data[i] = (byte + idx * key[0]) % 256

    return bytes(data)
```

Нужно заметить, что:

- от ключа используется всего один байт (`key[0]`), его можно перебрать
- изначальную сумму символов флага (`sum(data)`) также можно перебрать

Итоговая сложность перебора — 256^2. На каждой итерации нужно аккуратно развернуть операции в шифре, затем проверить полученный текст на совпадение с форматом флага.

**Пример решения**: [solution/solver.py](solution/solver.py)

## Флаг

```
Cup{s1mple_rev3rs1ble_c1pher}
```
