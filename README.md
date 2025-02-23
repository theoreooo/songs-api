# Songs API
Songs API – это онлайн-библиотека песен, реализованная на Go. Проект демонстрирует использование REST API, Redis для кеширования наиболее часто запрашиваемых данных, интеграцию с PostgreSQL через GORM, работу с внешним API для обогащения данных, Swagger-документацию и логирование (с использованием logrus). Кроме того, проект снабжён Dockerfile и docker-compose.yml для простого развёртывания.

## Особенности проекта
**REST API с эндпоинтами для:**
- Получения списка песен с расширенной фильтрацией (по группе, названию песни, дате релиза, тексту и ссылке).
- Получения детальной информации о песне по ID.
- Получения текста песни с пагинацией по куплетам.
- Добавления новой песни (с обогащением данных через внешний API).
- Частичного обновления песни (PATCH).
- Удаления песни.
- Нормализованная база данных:
- Данные о песнях разделены на две модели – Song и Artist (группа/исполнитель).
**Логирование:**
- Используется logrus для логирования на уровнях debug, info, error.
**Swagger-документация:**
- Документация API генерируется с помощью swaggo и доступна через Swagger UI.
**Docker и Docker Compose:**
Dockerfile для сборки образа и docker-compose.yml для развёртывания приложения вместе с PostgreSQL.
**Wait-for-postgres:**
Используется скрипт для ожидания готовности базы данных перед запуском приложения.

## Установка и запуск
**Предварительные требования**
Go (рекомендуется 1.23.5 или выше)
Docker и Docker Compose (если планируете запускать в контейнерах)
PostgreSQL (если запускаете локально без Docker)

**Запуск через Docker Compose**
Запустите Docker Compose:
```bash
docker-compose up --build
```
Это запустит два сервиса:

- db: PostgreSQL с базой данных music_db
- app: Ваше Go-приложение (Songs API)
**Запуск без Docker**
1. Установите зависимости:
```bash
go mod download
```
2. Настройте файл конфигурации .env
3. Если вы запускаете Redis локально, просто установите его (через пакетный менеджер или скачайте образ), запустите Redis‑сервер на порту 6379, а в переменной окружения REDIS_ADDR пропишите localhost:6379
3. Убедитесь, что PostgreSQL запущен и настроен согласно переменной DATABASE_URL в файле .env:
```bash
DATABASE_URL=host=db user=postgres password=0845 dbname=music_db port=5432 sslmode=disable TimeZone=Europe/Moscow
```
3. Запустите приложение:
```bash
go run ./cmd/songs/main.go
```
## Swagger-документация
Swagger-документация
После запуска приложения откройте в браузере:
```bash
http://localhost:8080/swagger/index.html
```
Здесь доступна полная документация API, описание всех эндпоинтов и примеры запросов/ответов.
## Внимание
В файле .env и в docker-compose.yml нужно указать MUSIC_API_URL API для получения данных о песне. Иначе метод POST для добавления песни не будет работать.
Для тестирования вы можете запустить этот API https://github.com/theoreooo/external-api-songs(Тогда ничего не надо будет менять в файлах env и docker-compose.yml, будет использоваться порт 8081).
Описание внешнего API:
```bash
paths:
  /info:
    get:
      parameters:
        - name: group
        in: query
        required: true
        schema:
          type: string
        - name: song
        in: query
        required: true
        schema:
          type: string

      responses:
        '200':
          description: Ok
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/SongDetail'

        '400':
          description: Bad request
        '500':
          description: Internal server error

components:
  schemas:
    SongDetail:
      required:
        - releaseDate
        - text
        - link
      type: object
      properties:
        releaseDate:
          type: string
          example: 16.07.2006
        text:
          type: string
          example: Ooh baby, don't you know I suffer?\nOoh baby, can
you hear me moan?\nYou caught me under false pretenses\nHow long
before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set
my soul alight
        link:
          type: string
          example: https://www.youtube.com/watch?v=Xsp3_a-PMTw
```
