# WS Backend

# Инструкция по локальному запуску приложения

## Зависимости
1. CLI для компиляции Swagger документации:
    ```shell
    $ go install github.com/swaggo/swag/cmd/swag@latest
    ```

    Пример использования:
    ```shell
    $ swag init -o web/static/apidocs --ot json -q -g cmd/app/main.go
    ```

2. CLI для создания миграций базы данных:

    [Инструкция по установке для разных ОС](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
    
    Пример использования:
    ```shell
    $ migrate create -dir migrations -seq -ext sql <migration_name> # Создает два файлика в папке migrations
    ```

   Также можно запустить скрипт для применения миграций:
   ```shell
   $ go run cmd/migrate/migrate.go
   ```

3. База данных:

    Мы используем PostgreSQL, поэтому для работы приложения необходима запущенная бд.

    (docker run --name ws -e POSTGRESQL_USERNAME=webstudio -e POSTGRESQL_PASSWORD=webstudio -e POSTGRESQL_DATABASE=ws -p 5432:5432 bitnami/postgresql:latest)


4. (Опционально) Make для использования скриптов из `Makefile`.


5. Mocks

   Моки используются для тестирования путем скрытия имплементации зависимостей.
   
   ```shell
   $ go install go.uber.org/mock/mockgen@latest
   ```

### Немного о криптографии приложения
Креды к бд находятся в конфиге, и их необходимо указывать зашифрованными.
Чтобы их сгенерировать, нужно выполнить скрипт(при необходимости указать там свои значения):
```shell
$ go run cmd/scripts/generate_db_creds.go
```

Далее указанные креды нужно вставить в свой конфиг.

Значения в `config.default.yml` используют креды `webstudio:webstudio`.

## Запуск приложения
```shell
$ make run
```

Посмотреть флаги приложения:
```shell
$ make run -help
```