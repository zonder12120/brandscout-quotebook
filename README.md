# brandscout-quotebook

Микросервис для управления цитатами, тестовое задание для Brand Scout, надеюсь, вы меня позовёте на следующий этап <3

## Содержание

- [Краткое описание](#краткое-описание)
- [Запуск](#запуск)
    - [Требования](#требования)
    - [Установка зависимостей](#установка-зависимостей)
        - [Установка всего необходимого на Windows](#установка-всего-необходимого-на-windows)
        - [Установка всего необходимого на Mac OS](#установка-всего-необходимого-на-mac-os)
        - [Установка make на Linux](#установка-make-на-linux)
    - [Настройка окружения](#настройка-окружения)
    - [Команды Makefile](#команды-makefile)
- [API Endpoints](#api-endpoints)
- [Примеры запросов](#примеры-запросов)

### Краткое описание

REST API сервис для управления цитатами с возможностью:
- Добавления новых цитат
- Получения всех цитат
- Получения случайной цитаты
- Фильтрации по автору
- Удаления цитат по ID

В сервисе реализовано логирование, данные хранятся в памяти (in-memory cache).  
Добавлена автоматическая очистка старых цитат при достижении лимита.  
Хэндлер, сервсис и in-memory cache покрыты тестами.  
Структуру проекта реализовывал опираясь на https://github.com/golang-standards/project-layout/

Для упрощения проверки я написал простой autofill: https://github.com/zonder12120/quotebook-autofill

### Запуск

#### Требования
- Go 1.24+
- Docker
- Make 

#### Установка зависимостей
##### Установка всего необходимого на Windows

Для запуска на Windows потребуется установить Chocolatey, инструкция: [https://docs.chocolatey.org/en-us/choco/setup/#install-with-cmdexe](https://docs.chocolatey.org/en-us/choco/setup/#install-with-cmdexe)

После установки Chocolatey перезапускаем Powershell/cmd

Затем в PowerShell/cmd запускаем команду:
```text
choco install make
```
После этого проверяем успех установки:
```text
make -v
```
Если получили что-то похожее на текст ниже, значит всё ок:
```text
GNU Make 4.4.1
Built for Windows32
Copyright (C) 1988-2023 Free Software Foundation, Inc.
License GPLv3+: GNU GPL version 3 or later https://gnu.org/licenses/gpl.html
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
```
Качаем, устанавливаем и запускаем Docker desktop: [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/)

##### Установка всего необходимого на Mac OS

Устанавливаем Homebrew: [https://brew.sh/](https://brew.sh/)

Там в разделе **Install Homebrew** просто копируем команду и исполняем в терминале.

После этого проверяем успех установки:
```text
brew -v
```
Если получили что-то похожее на текст ниже, значит всё ок:
```text
Homebrew 4.4.8
```
Затем выполняем команду:
```text
brew install make
```
После этого проверяем успех установки:
```text
make -v
```
Если получили что-то похожее на текст ниже, значит всё ок:
```text
GNU Make 3.81
Copyright (C) 2006 Free Software Foundation, Inc.
This is free software; see the source for copying conditions.
There is NO warranty; not even for MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE.
```
Качаем, устанавливаем и запускаем Docker desktop: [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/)

##### Установка make на Linux

Серьёзно? Если вы работаете на Линуксе, то и без меня знаете, что нужно делать.

#### Настройка окружения

1. Отредактируйте config/.env, опирайтесь на .env.dist:
```text
PORT=8080
QUOTES_LIMIT=1000
LOG_LEVEL=debug
```

**PORT -** порт для запуска сервера

**QUOTES_LIMIT -** максимальное количество хранимых цитат (при превышении старые удаляются)

**LOG_LEVEL -** уровень логирования (реализованы: debug, info, warn, error)

#### Команды Makefile
```text
# Сборка образа
make build

# Запуск контейнера
make run

# Остановка контейнера
make stop

# Перезапуск контейнера
make restart

# Просмотр логов
make logs

# Очистка образа
make clean
```

### API Endpoints
| Метод  | Путь                         | Описание                       |
|--------|------------------------------|--------------------------------|
| POST   | /quotes                      | Добавить новую цитату          |
| GET    | /quotes                      | Получить все цитаты            |
| GET    | /quotes/random               | Получить случайную цитату      |
| GET    | /quotes?author={name}        | Фильтр по автору               |
| DELETE | /quotes/{id}                 | Удалить цитату по ID           |

### Примеры запросов
Добавление цитаты:
```text
curl -X POST http://localhost:8080/quotes \
  -H "Content-Type: application/json" \
  -d '{"author":"Confucius", "quote":"Life is simple, but we insist on making it complicated."}'
```

Получение всех цитат:
```text
curl http://localhost:8080/quotes
```

Получение случайной цитаты:
```text
curl http://localhost:8080/quotes/random
```

Фильтрация по автору:
```text
curl http://localhost:8080/quotes?author=Confucius
```

Удаление цитаты:
```text
curl -X DELETE http://localhost:8080/quotes/1
```

