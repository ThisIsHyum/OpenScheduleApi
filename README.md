# OpenScheduleApi

⚠️ OpenScheduleApi находится в стадии активной разработки. Архитектура и API могут измениться.

**OpenScheduleApi** - REST API сервер для управления расписаниями колледжей

## Установка и запуск
### Бинарный файл
1. перейдите на [страницу релизов](https://github.com/ThisIsHyum/OpenScheduleApi/releases)
2. выберите бинарный файл под вашу ОС и архитектуру
3. установите [MySQL Server](https://www.mysql.com/)
4. поместите бинарный файл рядом с `.env`
5. запустите бинарный файл
```sh
# как программа
./OpenScheduleApi-linux-x86-64

# как фоновый процесс
./OpenScheduleApi-linux-x86-64 &
```
### Docker
1. установите docker
2. загрузите образ  
```sh
# последняя версия
docker pull ghcr.io/thisishyum/openscheduleapi:latest

# последняя бета-версия
docker pull ghcr.io/thisishyum/openscheduleapi:edge
```
3. запустите контейнер с `.env`  
```sh
#последняя версия
docker run -p 3530:3530 --env-file .env ghcr.io/thisishyum/openscheduleapi:latest

#последняя бета-версия
docker run -p 3530:3530 --env-file .env ghcr.io/thisishyum/openscheduleapi:edge
```

## Конфигурация
Конфигурация приложения задается через переменные окружения. Пример находится в файле [.env.example](https://github.com/ThisIsHyum/OpenScheduleApi/blob/main/.env.example)
1. **OSA_SERVER_HOST**  
IP-адрес или домен сервера  
По умолчанию: _localhost_
2. **OSA_SERVER_PORT**   
Порт сервера  
По умолчанию: _3530_
3. **OSA_ADMINTOKEN**  
Токен администратора  
_Обязательно_
4. **OSA_DB_HOST**  
IP-адрес или домен MySQL сервера  
_Обязательно_
5. **OSA_DB_PORT**  
Порт MySQL сервера  
_Обязательно_
6. **OSA_DB_USER**  
Имя пользователя MySQL сервера  
_Обязательно_
7. **OSA_DB_PASSWORD**  
Пароль пользователя MySQL сервера  
_Обязательно_
8. **OSA_DB_NAME**  
  Имя базы данных  
  _Обязательно_
## Принцип работы
1. Администратор создаёт колледж, указывая:
	- название колледжа
	- список кампусов
2. В ответ API возвращает токен парсера
3. Парсер использует этот токен для:
	- обновления учебных занятий
	- групп студентов
	- расписания звонков
4. Клиенты получают данные через публичные эндпоинты

## Документация API
- Все эндпоинты:  
  https://github.com/ThisIsHyum/OpenScheduleApi/blob/main/api/openapi.yml
- Модели:  
  https://github.com/ThisIsHyum/OpenScheduleApi/blob/main/api/models.yml
- Параметры:  
  https://github.com/ThisIsHyum/OpenScheduleApi/blob/main/api/parameters.yml
- Ошибки:  
  https://github.com/ThisIsHyum/OpenScheduleApi/blob/main/api/errors.yml
## Примеры запросов
1. Создание колледжа  
```sh
curl -X POST http://localhost:3530/admin/parser \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "collegeName": "test_college",
    "campusNames": ["campus1", "campus2", "campus3"]
  }'
```
Ответ:
```json
{
	"token": "PARSER_TOKEN"
}
```
2. Обновление групп кампуса  
```sh
curl -X POST http://localhost:3530/parser/groups \
  -H "Authorization: Bearer PARSER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
		"campusId": 1,
		"studentGroupNames": ["group1", "group2", "group3"]
  }'
```
3. Получение групп кампуса  
```sh
curl http://localhost:3530/campuses/1/groups
```
Ответ:
```json
[
	{
		"studentGroupId": 1,
		"name": "group1",
		"campusId": 1
	},
	{
		"studentGroupId": 2,
		"name": "group2",
		"campusId": 1
	},
	{
		"studentGroupId": 3,
		"name": "group3",
		"campusId": 1
	}
]
```

## Клиенты
- [osago](https://github.com/ThisIsHyum/osago)  
  Клиентская библиотека для Go
- [osa-cli](https://github.com/ThisIsHyum/osa-cli)  
  CLI-клиент, написанный на Go
- [osars](https://github.com/Bircoder432/osars)  
  Клиентская библиотека для Rust
- [osatui](https://github.com/Bircoder432/osatui)  
  TUI-клиент, написанный на Rust
## Парсеры
- [tkpst_parser](https://github.com/ThisIsHyum/tkpst_parser)  
  парсер расписания колледжа ТКПСТ, написанный на Go
