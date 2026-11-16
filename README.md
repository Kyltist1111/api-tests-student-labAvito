# API-тесты для проекта https://testboard.avito.com/

## Полезные ссылки

1. [Swagger-спецификация](https://testboard.avito.com/swagger/index.html)
2. [Техническое задание](https://docs.google.com/document/d/1IZ9SxEcePoQT6lyEmukycp1qNok4HZtok73UQnPwKSA/edit?usp=sharing)
3. [Библиотека для написания API-тестов](https://github.com/steinfletcher/apitest)
4. [Мощная библиотека анализа JSON для Go](https://github.com/tidwall/gjson)

---

## Реализованные сценарии (Лабораторная работа №3)

### Часть 1: Добавление объявления (`tests/scenarios/advertisement/`)
- **Позитивный сценарий:** Создание объявления со всеми параметрами (включая загрузку фото через `multipart/form-data`), проверка получения по ID, поиск по полному вхождению названия и доступность фотографий.
- **Негативные сценарии:** Валидация обязательных полей (пустой `title`, отрицательная `price`), проверка корректного ответа при отсутствии или невалидности токена авторизации (401 Unauthorized).

### Часть 2: Работа с избранным (`tests/scenarios/favorites/`)
- **Позитивный сценарий:** Полный цикл добавления чужого объявления в избранное, проверка его наличия в списке, удаление из избранного. Реализована устойчивая проверка удаления самого объявления с учетом бизнес-ограничений прав доступа.
- **Негативные сценарии:** Попытка добавить несуществующее объявление (404/400), попытка добавить своё собственное объявление (403 Forbidden), повторное добавление того же объявления (409 Conflict), запросы без авторизации.

---

## Запуск тестов

1. Установите любое IDE для разработки: [Jetbrains GoLand](https://www.jetbrains.com/ru-ru/go/) (Платное или по студенческой лицензии), [VS Code](https://code.visualstudio.com/) (бесплатное).
2. Установите [Git](https://git-scm.com/install/) на компьютер
3. [Установите](https://go.dev/doc/install) Go версии 1.25.1
4. [Настройте](https://docs.github.com/ru/authentication/connecting-to-github-with-ssh) доступ к своему Github-аккаунту
   по SSH-ключу.
5. Склонируйте себе репозиторий: `git clone <путь к вашему репозиторию>`
6. Подтяните необходимые зависимости через команду в корне проекта:

```bash
   go mod tidy
```

7. В файле `.env.override` (создайте его в корне, если нет) выставьте правильные значения для переменных окружения:

```
TEST_LOGIN=your_login
TEST_PASSWORD=your_password
```

8. Для запуска всех тестов достаточно выполнить команду в корне проекта:

```bash
  go test -v ./...
```

9. Отдельный пакет с тестами можно запустить командой:

```bash
  go test -v ./tests/scenarios/myAdvertisement/
```

10. Запуск тестов лабораторных:

Все тесты проекта:
```bash
   go test -v -count=1 ./...
```

Только тесты объявлений (Часть 1):
```bash
  go test -v -count=1 ./tests/scenarios/advertisement/...
```

Только тесты избранного (Часть 2):
```bash
  go test -v -count=1 ./tests/scenarios/favorites/...
```


## Логирование

Для включения логирования request/response в консоль: выставьте переменную окружения `DEBUG=true` в
`.env.override`

## Запуск линтера

1. Установить приложение [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/)
2. Запуск линтера (из корня репозитория):

```bash
  golangci-lint run
```
