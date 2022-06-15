CREATE TABLE IF NOT EXISTS restaurants
(
    id                   SERIAL PRIMARY KEY,
    name                 VARCHAR(255)   NOT NULL,
    average_waiting_time INTEGER        NOT NULL,
    average_check        DECIMAL(27, 4) NOT NULL
);

CREATE TABLE IF NOT EXISTS tables
(
    id            SERIAL PRIMARY KEY,
    restaurant_id INTEGER NOT NULL,
    seats_number  INTEGER NOT NULL,
    CONSTRAINT fk_tables_restaurants FOREIGN KEY (restaurant_id) REFERENCES restaurants (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS bookings
(
    id               SERIAL PRIMARY KEY,
    client_name      VARCHAR(255) NOT NULL,
    client_phone     VARCHAR(11)  NOT NULL,
    booked_date      DATE         NOT NULL,
    booked_time_from TIME         NOT NULL,
    booked_time_to   TIME         NOT NULL
);

CREATE TABLE IF NOT EXISTS bookings_tables
(
    id         SERIAL PRIMARY KEY,
    booking_id INTEGER NOT NULL,
    table_id   INTEGER NOT NULL,
    CONSTRAINT fk_bookings_tables_bookings FOREIGN KEY (booking_id) REFERENCES bookings (id) ON DELETE CASCADE,
    CONSTRAINT fk_bookings_tables_tables FOREIGN KEY (table_id) REFERENCES tables (id) ON DELETE CASCADE
);

/*
 Функция get_available_tables возвращает таблицу вида tables с информацией о столиках, свободных для бронирования.
 */
CREATE OR REPLACE FUNCTION get_available_tables(
    desired_booking_date date, -- желаемая дата брони
    desired_booking_time time -- желаемое время брони (столик бронируется на 2 часа с этого момента времени)
)
    RETURNS TABLE
            (
                id            INTEGER,
                restaurant_id INTEGER,
                seats_number  INTEGER
            )
AS
$$
BEGIN
    -- столики которые ни разу не бронировались
    RETURN QUERY
        SELECT tables.id, tables.restaurant_id, tables.seats_number
        FROM tables
        WHERE tables.id NOT IN (SELECT bookings_tables.table_id FROM bookings_tables)
        UNION
        -- столики которые хотя бы раз бронировались
        SELECT tables.id, tables.restaurant_id, tables.seats_number
        FROM tables
                 JOIN bookings_tables bt on tables.id = bt.table_id
                 JOIN bookings b on b.id = bt.booking_id
        WHERE is_table_available(bt.table_id, desired_booking_date, desired_booking_time);
END;
$$ LANGUAGE plpgsql;

/*
    Функция is_table_available проверяет, можно ли забронировать столик на 2 часа в желаемые дату и время.
    Алгоритм:
        1) Среди столиков, которые хотя бы 1 раз бронировались в желаемую дату, ищется выбранный столик.
        Если он не находится, значит, он свободен в этот день, и его точно можно забронировать.
        Если записи были найдены, значит, столик в какой-то временной промежуток занят, и нужно проверить,
        будет ли наложение желаемой брони на уже зарегистрированную в системе. Переходим к п. 2
        2) Получив таблицу временных промежутков и добавив туда желаемый, объединяем временные промежутки,
        накладывающиеся друг на друга, в обшие, более длинные временные промежутки.
        Таким образом, если временной промежуток желаемой брони будет слит с уже зарегистрированными бронями,
        следовательно, начальное количество интервалов и конечное (после слияния) будут различаться,
        то столик не будет доступен для брони в желаемое время.
        Если же количество интервалов до и после слияния равны, то столик можно забронировать.
 */
CREATE OR REPLACE FUNCTION is_table_available(
    checked_table_id int, -- ID столика
    desired_booking_date date, -- желаемая дата брони
    desired_booking_time time -- желаемое время брони (столик бронируется на 2 часа с этого момента времени)
)
    RETURNS BOOLEAN -- возвращает TRUE, если забронировать можно, иначе - FALSE
AS
$$
DECLARE
    rows_before_merge INTEGER;
    rows_after_merge  INTEGER;
BEGIN
    -- ищем временные промежутки броней столиков, которые хотя бы раз бронировались в выбранный день
    CREATE TEMP TABLE booking_intervals
    AS
    SELECT booked_time_from, booked_time_to
    FROM tables
             JOIN bookings_tables bt on tables.id = bt.table_id
             JOIN bookings b on b.id = bt.booking_id
    WHERE booked_date = desired_booking_date
      AND table_id = checked_table_id
    UNION
    -- добавляем временной промежуток желаемой брони к полученным
    VALUES (desired_booking_time, desired_booking_time + interval '2 hours');

    /*
    если остался только один временной промежуток (время желаемой брони), значит, столик вообще не бронировался
    в выбранную дату, и его можно забронировать
     */
    rows_before_merge := (SELECT COUNT(*) FROM booking_intervals);
    IF rows_before_merge = 1 THEN
        DROP TABLE booking_intervals;
        RETURN TRUE;
    END IF;

    rows_after_merge := (
        SELECT COUNT(*)
        FROM (
                 WITH rng(s, e) AS (
                     SELECT *
                     FROM booking_intervals
                 )
                 SELECT -- min/max по группе
                        min(s) s,
                        max(e) e
                 FROM (
                          SELECT *,
                                 sum(ns::integer) OVER (ORDER BY s, e) grp -- определение групп
                          FROM (
                                   SELECT *,
                                          coalesce(s > max(e)
                                                       OVER (ORDER BY s, e ROWS BETWEEN UNBOUNDED PRECEDING AND 1 PRECEDING),
                                                   TRUE) ns -- начало правее самого правого из предыдущих концов == разрыв
                                   FROM rng
                               ) t
                      ) t
                 GROUP BY grp
             ) merged_intervals
    );

    DROP TABLE booking_intervals;
    RETURN rows_before_merge = rows_after_merge;
END;
$$ LANGUAGE plpgsql;