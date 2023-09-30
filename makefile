PROJECT ?= 'note-service'
MAIN_SERVICE ?= 'ns-app'

## make logging: команда для просмотра логов
logging:
	docker-compose logs -f --tail="25" ${MAIN_SERVICE}

## make ps: команда для вывода списка всех запушенных контейнеров
ps:
	docker-compose ps

## make restart: команда для перезапуска заданного контейнера
restart:
	docker-compose restart ${MAIN_SERVICE}

## make rebuild: команда для "перекомплиляции" композиции
rebuild:
	docker-compose up -d --build

# make up: команда для поднятия контейнера с определенным именем
up:
	docker-compose -p ${PROJECT} up -d

## make down: команда для остановки всех контейнеров
down:
	docker-compose down

## make start: команда запуска заданного контейнера
start:
	docker-compose start ${MAIN_SERVICE}

## make stop: Команда остановки заданного контейнера
stop:
	docker-compose stop ${MAIN_SERVICE}

## make build: команда для создания заданного контейнера проекта
build:
	docker-compose -p ${PROJECT} build