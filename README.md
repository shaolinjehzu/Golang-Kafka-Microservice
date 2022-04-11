## Golang Balance Microservice Example ###

### Пояснительная записка

#### Что было использовано
* [Postgresql](https://www.postgresql.org/) as DB<br/>
* [Kafdrop](https://github.com/obsidiandynamics/kafdrop) web UI for kafka<br/>
* [Kafka](https://github.com/segmentio/kafka-go) as messages broker<br/>
* [sqlx](https://github.com/jmoiron/sqlx) - Extensions to database/sql.
* [sql/pq](https://github.com/lib/pq) - PostgreSQL driver for Go
* [viper](https://github.com/spf13/viper) - Go configuration with fangs
* [zap](https://github.com/uber-go/zap) - Logger
* [validator](https://github.com/go-playground/validator) - Go Struct and Field validation
* [migrate](https://github.com/golang-migrate/migrate) - Database migrations. CLI and Golang library.
* [CompileDaemon](https://github.com/githubnemo/CompileDaemon) - Compile daemon for Go
* [Docker](https://www.docker.com/) - Docker


#### Основные структуры данных:
    Balance         - структура модели данных в БД
    Task            - структура получаемых заданий для воркеров
    ErrorMessage    - структура сообщения отправляемых при возникновении ошибки обновления 
    SuccessMessage  - структура сообщения отправляемых при успешном обновлении 

#### Очереди:
    update-balance          - топик для задач
    dead-letter-queue       - топик для сообщений об ошибках
    success-letter-queue    - топик для сообщений о выполненных заданиях

#### Механизм работы:
    При запуске приложения, создаётся экземпляр BalancesConsumerGroup, после чего происходит 
    запуск Consumer'а, который получает задания из топика и запускает воркеров.
    
    Воркер запускается в отдельной горутине и обрабатывает задачу. 
    В случае успешной обработки, сообщение отправляется в топик с выполненыыми заданиями
    В случае ошибки, сообщение об этом отправляется в топик с ошибками

#### 🚀 Docker-compose файлы:
    docker-compose.yml - запускает postgresql, kafka, zookeper, kafdrop and migrate контейнеры

#### Команды для развертывания:
    make local
    make create_topics
    make run

### Kafdrop UI:

http://localhost:9000
