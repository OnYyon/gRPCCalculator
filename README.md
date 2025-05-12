# gRPCCalculator
# GoroutineRPNServer

## Описание
Это серевер, кторый обрабатывает математическое выражение(может состоять только из `цифр`, `+`, `-`, `*`, `/`). Так же есть `web-interace`. Вкраце, есть `api`, которое получает от `web-interace` выражение, переводит его в RPN(Обратная польская нотация) и разбивате на "простые" таски.
Контакт тг: @OnYyon
## Навигация
- [Установка и запуск](#установка)
- [Как работает](#как-все-работает)
- [Примеры](#примеры)
## Установка и запуск
Установка `docker`(если не имееться в наличии)
- Идем по ссылке https://www.docker.com/
- Выбираем нужную платмформу и скачиваем

Клонируем репозиторий с `github`
```bash
git clone https://github.com/OnYyon/GoroutineRPNServer.git
```
Переходи в дирикторию проекта
```bash
cd GoroutineRPNServer
```
Запуск
1. Три отдельных окна терминала
```
go run cmd/orchestrator/main.go  
```
```
go run cmd/worker/main.go    
```
```
go run web/main.go 
```
2. Docker-compose
```
docker-compose up --build 
```
>ВАЖНО: Ждем запуска программ
> Будет выведнно:
>>orchestrator-1  | Starting orhcestrator on localhost:8080
>>web-1           | Starting web-interface on localhost:8081
>>agent-1         | Starting gouroutine: 0
>>agent-1         | Starting gouroutine: 1
>> ...
>>agent-1         | Starting gouroutine: 12
>
> Только после этого все запуститься и можно пойти на [Web interface](http://127.0.0.1:8081)

## Как все работает
![](./img/1img.jpg)
На картинке показан процесс поставление задачи клиентами. Все они обращаються к `оркестартору` он парсист выражение в RPN и делит RPN на таски. Количество воркеров задаеться переменнной среды `COMPUTING_POWER`
Посмотрим подробнее как работает раздение на таски.
```
RPN expression: 1 2 3 + + 4 5 + + 6 7 + + 4 3 * 8 / +
таски: (2 3 +), (4 5 +), (6 7 +), (4 3 *)

RPN expression: 1 5 + 9 + 13 + 12 8 / +
таски: (1 5 +), (12 8 /)

RPN expression: 6 9 + 13 + 1.5 +
таски: (6 9 +)

RPN expression: 15 13 + 1.5 +
таски: (15 13 +)

RPN expression: 28 1.5 +
таски: (28 1.5 +)

RPN expression: 29.5 - итоговый результат
```
Програма ищет тройки, кторые может выполнить сразу то есть все операторы однозначно определенны. И так пока не получит результат

## Примеры
Регистрация пользователя
```
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"login": "test_user", "password": "strong_password"}'
```
Неправильные запросы:
```
# Без пароля
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"login": "test_user"}'

# Без логина
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"password": "strong_password"}'

# Пустые данные
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{}'
  ```
Авторизация
```
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"login": "test_user", "password": "strong_password"}'
```
Неправильные запросы:
```
 Неверный пароль
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"login": "test_user", "password": "wrong_password"}'

# Несуществующий пользователь
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"login": "unknown_user", "password": "strong_password"}'
  ```
Добавление выражения (требует аутентификации)
```
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"login": "test_user", "password": "strong_password"}')
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
```
```curl -X POST http://localhost:8080/api/expression \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"expression": "2 + 2 * 2"}'
  ```
Получение выражения по ID
  ```# Предположим, что ID выражения "123"
curl -X GET "http://localhost:8080/api/expression/123" \
  -H "Authorization: Bearer $TOKEN"
  ```
Получение списка выражений
```
curl -X GET http://localhost:8080/api/expressions \
  -H "Authorization: Bearer $TOKEN"
  ```