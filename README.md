# Сервис бронирования столиков в ресторанах

Тестовое задание на стажировку в [Aero](https://aeroidea.ru/).

## Содержание

- [Задание](#Задание)

## Задание

Представь знакомую всем ситуацию. Ты собираешься с друзьями в ресторан и хочешь забронировать столик. В вашем городе открыты 3 ресторана: «Каравелла», «Молодость», «Мясо и Салат» (любые совпадения случайны).

**У ресторанов свои особенности:**

1. «Каравелла»:

   — 10 столиков: 6 столиков вмещает до 4 человек. 2 столика вмещают до 3 человек. 2 столика вмещает до 2 человек;

   — Среднее время ожидания блюда: 30 минут;

   — Средний чек: 2000 рублей.

2. «Молодость»

   — 3 столика, каждый из которых вмещает 3 человека;

   — Среднее время ожидания блюда: 15 минут;

   — Средний чек: 1000 рублей.

3. «Мясо и Салат»

   — 6 столиков: 2 столика вмещает до 8 человек. 4 столика вмещают до 3 человек;

   — Среднее время ожидания блюда: 60 минут;

   — Средний чек: 1500 рублей.


**Необходимо создать небольшую систему бронирования столиков.** Предусматриваются следующие шаги:

1. Пользователь указывает количество человек и желаемое время посещения ресторана.
2. Время брони - 2 часа. Рестораны работают с 9:00 до 23:00 (Последнюю бронь можно создать на 21:00).
3. Сдвигать столики - можно.
4. Система предлагает доступные варианты (рестораны).
5. Необходимо указывать актуальное количество свободных мест.
6. Необходимо отсортировать подходящие варианты по возрастанию среднего времени ожидания и среднего чека.
7. Необходимо скрывать недоступные варианты.
8. Пользователь указывает имя и номер телефона и завершает процесс бронирования.

**В качестве результата ожидается** ссылка на github-репозиторий, в котором находятся:

1. http-сервер, который общается с миром на языке REST API.
2. readme.md-файл, в котором подробно описана инструкция по установке системы

   — Все действия по установке должны быть автоматизированы (= вызов команд)

3. sql-дамп базы данных.

**Нужно не забыть:**

1. Проверить входные параметры. Сделать защиту от дурака.
2. Написать много комментариев к коду.
3. Проверить полученный результат дважды (а лучше трижды).

**Будет круто, если** (но совсем не критично, если не получится):

1. В качестве БД ты выберешь postgres.
2. Сделаешь визуальный интерфейс. Можно сверстать самому или использовать готовые решения.
3. Напишешь код на php или go (вообще идеально).
4. Используешь docker-compose.
5. Вынесешь все настройки подключения к базе в переменные окружения (env-параметры).
