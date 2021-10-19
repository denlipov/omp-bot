# Ozon Marketplace Project

![schema](schema.png)

Дальше везде используются **placeholder**-ы:
- `{domain}`,`{Domain}`
- `{subdomain}`,`{Subdomain}`

Например, для поддомена `package` из домена `logistic` значение **placeholder**-ов будет:
- `{domain}`,`{Domain}` = `logistic`,`Logistic`
- `{subdomain}`,`{Subdomain}` = `package`,`Package`
- `{domain}`/`{subdomain}` = `logistic`/`package`
---

### Задание 1

1. Сделать форк **ozonmp/omp-bot** репозитория в свой профиль
2. Запросить у своего тьютора свой домен/поддомен: **{domain}/{subdomain}**
3. Добавить в ветку `feature/task-1` своего форка поддержку следующих команд:
```
/help__communication__request — print list of commands
/get__communication__request — get a entity
/list__communication__request — get a list of your entity (💎: with pagination via telegram keyboard)
/delete__communication__request — delete an existing entity

/new__communication__request — create a new entity // not implemented (💎: implement list fields via arguments)
/edit__communication__request — edit a entity      // not implemented
```
4. Сделать PR из ветки `feature/task-1` своего форка в ветку `master` своего форка
5. Отправить ссылку на PR личным сообщением своему тьютору до конда дедлайна сдачи (см. таблицу прогресса)

#### Рецепт

Для добавления поддержки команд в рамках своего поддомена:

1. Написать структуру `Request` с методом `String()`
2. Написать интерфейс `RequestService` и **dummy** имплементацию
3. Написать интерфейс `RequestCommander` по обработке команд

---

2. Реализовать `RequestService` в **internal/service/communication/request/**

```go
package request

import "github.com/ozonmp/omp-bot/internal/model/communication"

type RequestService interface {
  Describe(requestID uint64) (*communication.Request, error)
  List(cursor uint64, limit uint64) ([]communication.Request, error)
  Create(communication.Request) (uint64, error)
  Update(requestID uint64, request communication.Request) error
  Remove(requestID uint64) (bool, error)
}

type DummyRequestService struct {}

func NewDummyRequestService() *DummyRequestService {
  return &DummyRequestService{}
}

// ...
```

---

3. Реализовать `RequestCommander` по обработке команд в **internal/app/commands/communication/request/**

```go
package request

import (
  model "github.com/ozonmp/omp-bot/internal/model/communication"
  service "github.com/ozonmp/omp-bot/internal/service/communication/request"
)

type RequestCommander interface {
  Help(inputMsg *tgbotapi.Message)
  Get(inputMsg *tgbotapi.Message)
  List(inputMsg *tgbotapi.Message)
  Delete(inputMsg *tgbotapi.Message)

  New(inputMsg *tgbotapi.Message)    // return error not implemented
  Edit(inputMsg *tgbotapi.Message)   // return error not implemented
}

func NewRequestCommander(bot *tgbotapi.BotAPI, service service.RequestService) RequestCommander {
  // ...
}
```
