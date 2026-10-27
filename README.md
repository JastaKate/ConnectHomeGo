# ConnectHome

REST API для управления домами. Переношу свой пет Java-проект (https://github.com/JastaKate/home-service/tree/master) на Go,
чтобы разобраться с языком и разработкой backend-сервисов.

Пока реализованы создание дома и получение списка. Данные хранятся в PostgreSQL.

## Стек

Go, net/http, PostgreSQL, pgx, Goose, envconfig.

## Запуск

Нужны Go 1.27 и PostgreSQL. Подключение:

```sh
export DATABASE_URL='postgres://user:password@localhost:5432/connect_home?sslmode=disable'
```

Сервер запускается на `http://localhost:8080`.
Настройки также можно передать отдельными переменными из `.env.example`.
Приложение не загружает `.env` автоматически.
Локальный `compose.yaml` не хранится в репозитории.

## API

- `GET /health` - проверка работы сервера.
- `POST /v2/api/homes` - создание дома.
- `GET /v2/api/homes` - список домов.

Пример создания:

```sh
curl -X POST http://localhost:8080/v2/api/homes \
  -H 'Content-Type: application/json' \
  -d '{"name":"Мой дом","address":"Москва"}'
```

Код разделён на controller, service, repository, model и request.
