--
-- PostgreSQL database dump
--

-- Dumped from database version 13.4
-- Dumped by pg_dump version 13.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: get_available_tables(date, time without time zone); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.get_available_tables(desired_booking_date date, desired_booking_time time without time zone) RETURNS TABLE(id integer, restaurant_id integer, seats_number integer)
    LANGUAGE plpgsql
    AS $$
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
$$;


ALTER FUNCTION public.get_available_tables(desired_booking_date date, desired_booking_time time without time zone) OWNER TO postgres;

--
-- Name: get_unavailable_table_ids(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.get_unavailable_table_ids() RETURNS boolean
    LANGUAGE plpgsql
    AS $$
DECLARE
    cnt_before INTEGER;
    cnt_after  INTEGER;
BEGIN
    cnt_before := (SELECT COUNT(*)
                   FROM (
                            SELECT booked_time_from,
                                   booked_time_to -- среди столиков которые хотя бы раз бронировались
                            FROM tables
                                     JOIN bookings_tables bt on tables.id = bt.table_id
                                     JOIN bookings b on b.id = bt.booking_id
                            WHERE booked_date = date '2022-06-13'
                            UNION
                            VALUES (time '12:30', time '12:30' + interval '2 hours')
                        ) all_intervals
    );

    cnt_after := (
        SELECT COUNT(*)
        FROM (
                 WITH rng(s, e) AS (
                     SELECT *
                     FROM (
                              SELECT booked_time_from,
                                     booked_time_to -- среди столиков которые хотя бы раз бронировались
                              FROM tables
                                       JOIN bookings_tables bt on tables.id = bt.table_id
                                       JOIN bookings b on b.id = bt.booking_id
                              WHERE booked_date = date '2022-06-13'
                              UNION
                              VALUES (time '12:30', time '12:30' + interval '2 hours')
                          ) all_intervals
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

    RETURN cnt_before = cnt_after;
END;
$$;


ALTER FUNCTION public.get_unavailable_table_ids() OWNER TO postgres;

--
-- Name: get_unavailable_table_ids(date, time without time zone); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.get_unavailable_table_ids(wish_date date, wish_time time without time zone) RETURNS boolean
    LANGUAGE plpgsql
    AS $$
DECLARE
    cnt_before INTEGER;
    cnt_after  INTEGER;
BEGIN
    cnt_before := (SELECT COUNT(*)
                   FROM (
                            SELECT booked_time_from,
                                   booked_time_to -- среди столиков которые хотя бы раз бронировались
                            FROM tables
                                     JOIN bookings_tables bt on tables.id = bt.table_id
                                     JOIN bookings b on b.id = bt.booking_id
                            WHERE booked_date = wish_date
                            UNION
                            VALUES (wish_time, wish_time + interval '2 hours')
                        ) all_intervals
    );

    cnt_after := (
        SELECT COUNT(*)
        FROM (
                 WITH rng(s, e) AS (
                     SELECT *
                     FROM (
                              SELECT booked_time_from,
                                     booked_time_to -- среди столиков которые хотя бы раз бронировались
                              FROM tables
                                       JOIN bookings_tables bt on tables.id = bt.table_id
                                       JOIN bookings b on b.id = bt.booking_id
                              WHERE booked_date = wish_date
                              UNION
                              VALUES (wish_time, wish_time + interval '2 hours')
                          ) all_intervals
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

    RETURN cnt_before = cnt_after;
END ;
$$;


ALTER FUNCTION public.get_unavailable_table_ids(wish_date date, wish_time time without time zone) OWNER TO postgres;

--
-- Name: is_table_available(integer, date, time without time zone); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.is_table_available(checked_table_id integer, desired_booking_date date, desired_booking_time time without time zone) RETURNS boolean
    LANGUAGE plpgsql
    AS $$
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
$$;


ALTER FUNCTION public.is_table_available(checked_table_id integer, desired_booking_date date, desired_booking_time time without time zone) OWNER TO postgres;

--
-- Name: array_cat_agg(anyarray); Type: AGGREGATE; Schema: public; Owner: postgres
--

CREATE AGGREGATE public.array_cat_agg(anyarray) (
    SFUNC = array_cat,
    STYPE = anyarray,
    INITCOND = '{}'
);


ALTER AGGREGATE public.array_cat_agg(anyarray) OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: bookings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.bookings (
    id integer NOT NULL,
    client_name character varying(255) NOT NULL,
    client_phone character varying(11) NOT NULL,
    booked_date date NOT NULL,
    booked_time_from time without time zone NOT NULL,
    booked_time_to time without time zone NOT NULL
);


ALTER TABLE public.bookings OWNER TO postgres;

--
-- Name: bookings_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.bookings_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.bookings_id_seq OWNER TO postgres;

--
-- Name: bookings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.bookings_id_seq OWNED BY public.bookings.id;


--
-- Name: bookings_tables; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.bookings_tables (
    id integer NOT NULL,
    booking_id integer NOT NULL,
    table_id integer NOT NULL
);


ALTER TABLE public.bookings_tables OWNER TO postgres;

--
-- Name: bookings_tables_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.bookings_tables_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.bookings_tables_id_seq OWNER TO postgres;

--
-- Name: bookings_tables_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.bookings_tables_id_seq OWNED BY public.bookings_tables.id;


--
-- Name: restaurants; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.restaurants (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    average_waiting_time integer NOT NULL,
    average_check numeric(27,4) NOT NULL
);


ALTER TABLE public.restaurants OWNER TO postgres;

--
-- Name: restaurants_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.restaurants_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.restaurants_id_seq OWNER TO postgres;

--
-- Name: restaurants_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.restaurants_id_seq OWNED BY public.restaurants.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Name: tables; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tables (
    id integer NOT NULL,
    restaurant_id integer NOT NULL,
    seats_number integer NOT NULL
);


ALTER TABLE public.tables OWNER TO postgres;

--
-- Name: tables_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.tables_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.tables_id_seq OWNER TO postgres;

--
-- Name: tables_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.tables_id_seq OWNED BY public.tables.id;


--
-- Name: bookings id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings ALTER COLUMN id SET DEFAULT nextval('public.bookings_id_seq'::regclass);


--
-- Name: bookings_tables id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings_tables ALTER COLUMN id SET DEFAULT nextval('public.bookings_tables_id_seq'::regclass);


--
-- Name: restaurants id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.restaurants ALTER COLUMN id SET DEFAULT nextval('public.restaurants_id_seq'::regclass);


--
-- Name: tables id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tables ALTER COLUMN id SET DEFAULT nextval('public.tables_id_seq'::regclass);


--
-- Data for Name: bookings; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.bookings (id, client_name, client_phone, booked_date, booked_time_from, booked_time_to) FROM stdin;
1	Павел	89278575533	2022-06-17	15:33:00	17:33:00
\.


--
-- Data for Name: bookings_tables; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.bookings_tables (id, booking_id, table_id) FROM stdin;
1	1	13
\.


--
-- Data for Name: restaurants; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.restaurants (id, name, average_waiting_time, average_check) FROM stdin;
1	Каравелла	30	2000.0000
2	Молодость	15	1000.0000
3	Мясо и Салат	60	1500.0000
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.schema_migrations (version, dirty) FROM stdin;
20220610192523	f
\.


--
-- Data for Name: tables; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tables (id, restaurant_id, seats_number) FROM stdin;
1	1	4
2	1	4
3	1	4
4	1	4
5	1	4
6	1	4
7	1	3
8	1	3
9	1	2
10	1	2
11	2	3
12	2	3
13	2	3
14	3	8
15	3	8
16	3	3
17	3	3
18	3	3
19	3	3
\.


--
-- Name: bookings_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.bookings_id_seq', 1, true);


--
-- Name: bookings_tables_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.bookings_tables_id_seq', 1, true);


--
-- Name: restaurants_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.restaurants_id_seq', 3, true);


--
-- Name: tables_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.tables_id_seq', 19, true);


--
-- Name: bookings bookings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings
    ADD CONSTRAINT bookings_pkey PRIMARY KEY (id);


--
-- Name: bookings_tables bookings_tables_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings_tables
    ADD CONSTRAINT bookings_tables_pkey PRIMARY KEY (id);


--
-- Name: restaurants restaurants_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.restaurants
    ADD CONSTRAINT restaurants_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: tables tables_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tables
    ADD CONSTRAINT tables_pkey PRIMARY KEY (id);


--
-- Name: bookings_tables fk_bookings_tables_bookings; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings_tables
    ADD CONSTRAINT fk_bookings_tables_bookings FOREIGN KEY (booking_id) REFERENCES public.bookings(id) ON DELETE CASCADE;


--
-- Name: bookings_tables fk_bookings_tables_tables; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings_tables
    ADD CONSTRAINT fk_bookings_tables_tables FOREIGN KEY (table_id) REFERENCES public.tables(id) ON DELETE CASCADE;


--
-- Name: tables fk_tables_restaurants; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tables
    ADD CONSTRAINT fk_tables_restaurants FOREIGN KEY (restaurant_id) REFERENCES public.restaurants(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

