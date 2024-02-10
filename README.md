# Вычислитель арифметических выражений(AEC)
## Сервер агент(воркер) каждые 30сек посылает ping на оркестратор о том что он жив 
## Иногда есть проблемы с большими числами при генерации ответа
## Если хотите отправлять запросы не через swagger то вот url - http://localhost:9999/
### tg - https://t.me/GusGus153

## Запуск
На ПК должен быть установлен <b>docker</b>(docker демон должен быть запушен) и <b>docker-compose>=2</b>

    git clone https://github.com/Filin153/AEC.git

    cd AEC

    docker-compose build

    docker-compose up -d

# Swagger схема
    http://localhost:8080/swagger/index.html

### Task
    Запрос post - / (добавляет заданние(каждый раз будет рандомный UserId его можно вставить в запрос чтобы отправлять их как 1 пользователь или оставить string))
    Запрос get - /task/{id} (принемает ID задания и выводит инфо о нем)

### Server
    Запрос post - /server/add/{id}/{add} (принемает ID сервера(воркера) и количество воркеров которое нужно добавить)
    Запрос get - /server/all (выводит все сервера(воркеры аодключенные на данный момент))
    Запрос delete - /server/del/{id} (принемает ID сервера(воркера) и удаляет его)

### User
    Запрос get - /user/{id}/ (принемает ID пользователя и возврашяет иноф о всех его заданиях)

# БД
Чтобы подключится нужно перейти на http://localhost:9009/

    Движок - PostgreSQL
    Сервер - db
    Имя пользователя - postgres
    Пароль - gus
    База данных - AEC

# Схема работы
https://miro.com/app/board/uXjVNvvRUYI=/?share_link_id=807062704044

![alt text](https://github.com/Filin153/AEC/blob/[branch]/img/img.png?raw=true)
    