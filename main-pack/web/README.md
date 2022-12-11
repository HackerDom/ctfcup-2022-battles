# CTFCup 2022 | Battles / main | web

## Название

> amazing_notes

## Описание

> Simple Notes Service (SNS). Закрытая экстра найтли сборка нового сервиса для aws. Безос два года ждал этот сервис, чтобы растоптать конкурентов этим прекрасным продуктом. Но мы его украли(Джефф забыл флэшку с исходниками в баре). Правда говорят, что аудит безопасности он пока не прошел...
>
> `Use http Luke. 8080 waiting for you`

## Раздатка

Отдаем все, что лежит в [public/](public/amazing_notes.tar.gz)

## Деплой

```
docker-compose up --build -d
```

## Решение
Флаг лежит в корне бакета notes, но ни один пользователь прошедший авторизацию и имеющий запись в бакете users не может его прочитать, т.к всем пользователям назначается домашняя папка из которой они могут читать.

Значит нужно каким-то образом получить пользователя который сможет пройти авторизацию, но не будет иметь userInfo.

Посмотрим на то, как создаются новые пользователи

```
func (service *NotesServiceImpl) CreateUser(credentials models.Credentials, ctx context.Context) error {
     ...

	msg := models.CreateUser{Credentials: credentials}.ToMessage()

// запрос через SQS отправляется воркеру
	msgId, err := service.queue.Push(msg.ToString(), ctx)
	if err != nil {
		glog.Errorf("Send failed: %s", err)
		return err
	}

// ждем минимум 5 секунд
	for range time.Tick(5 * time.Second) {
		inQueue, err := service.queue.InQueue(msgId)
		if err != nil {
			return err
		}

//если запрос пропал из очереди добавляем информацию об авторизации в память
//Причем клиент реализован так, что здесь достаточно, чтобы воркер взял задачу из очереди, удалять ее не обязательно
//если бы мы не ждали 5 секунд можно было бы поспамить запросами получить флаг
		if !inQueue {
			service.users[credentials.Login] = credentials
			break
		}
	}

	glog.Infof("Successfully add: %+v", credentials)

	return nil
}
```

Очевидно, что в сервисе есть гонка, но как ее эксплуатировать? Посмотрим, как устроено чтение из sqs


```
func (worker *daemonImpl) pull(ctx context.Context) error {
	glog.Info("Start pulling messages")

	var batch []aws.Message
	//Читаем батчи batchSize = 50
	for i := 0; i < batchSize; i++ {
		msgs, err := worker.queue.Pull(ctx)
		if err != nil {
			glog.Errorf("can't pull: %s", err)
			return err
		}

		batch = append(batch, msgs...)

		if len(msgs) == 0 {
			break
		}
	}

// сюда мы попадаем при двух условиях, 
// если прочитали больше 50 сообщений или 
// если не получили ни одного сообщения на последний запрос
	for _, msg := range batch {
		body := msg.Body
		worker.handleMessage(ctx, body)
		_ = worker.queue.Delete(msg)
	}

	return nil
}
```
В sqs есть параметер WaitTimeSeconds, который определяет сколько клиент будет ждать новых сообщений, причем ждать мы будем независимо от того прочитали лы мы MaxNumberOfMessages или нет
```
func (s *SqsClient) Pull(ctx context.Context) ([]Message, error) {
	res, err := s.client.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.queueUrl),
		MaxNumberOfMessages: aws.Int64(1), //читаем по одному сообщению
		WaitTimeSeconds:     aws.Int64(1), //ждем 1 секунду
	})
	if err != nil {
		glog.Errorf("receive: %s", err)
		return nil, err
	}


n
```

Таким образом мы можем остановить запись в s3 на промежуток времени равный WaitTimeSeconds * batchSize, чего нам как раз хватит, чтобы информацию об авторизации проросла в api.

т.е. для того, чтобы похакать сервис на нужно:
1. Начать раз в секунду отправлять любое запрос который пушит сообщение в sqs
2. Создать пользователя
3. Через 5 прочитать флаг

**Пример решения**: [sploit](solution/sploit/main.go)

## Флаг

```
Cup{0h_senpa1_mY_quEeEuEeE_fu11_0f_messages}
```
