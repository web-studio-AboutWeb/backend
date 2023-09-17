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


4. (Опционально) Make для использования скриптов из `Makefile`.


5. TODO: mock

### Немного о криптографии приложения
Мы используем ключи для шифрования важных данных. Для того, чтобы шифрование работало, необходимо указать ключ в конфиге.
Чтобы его сгенерировать, нужно выполнить скрипт:
```shell
$ go run cmd/scripts/key_and_db_creds.go
```

Далее значения(ключ и зашифрованные креды бд) нужно вставить в свой конфиг.

Значения в `config.default.yml` используют креды `webstudio:webstudio`.

## Запуск приложения
```shell
$ make run
```

Посмотреть флаги приложения:
```shell
$ make run -help
```