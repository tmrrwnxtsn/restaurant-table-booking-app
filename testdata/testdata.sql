SET client_encoding TO 'utf8';

-- Изначально есть 3 ресторана
INSERT INTO restaurants (name, average_waiting_time, average_check)
VALUES ('Каравелла', 30, 2000.00),
       ('Молодость', 15, 1000.00),
       ('Мясо и Салат', 60, 1500.00);

/*
 В ресторане "Каравелла" 10 столиков:
    6 столиков вмещает до 4 человек,
    2 столика вмещают до 3 человек,
    2 столика вмещает до 2 человек
 */
INSERT INTO tables (restaurant_id, seats_number)
VALUES (1, 4),
       (1, 4),
       (1, 4),
       (1, 4),
       (1, 4),
       (1, 4),
       (1, 3),
       (1, 3),
       (1, 2),
       (1, 2);

-- В ресторане "Молодость" 3 столика, каждый из которых вмещает 3 человека
INSERT INTO tables (restaurant_id, seats_number)
VALUES (2, 3),
       (2, 3),
       (2, 3);

/*
 В ресторане "Мясо и Салат" 6 столиков:
    2 столика вмещает до 8 человек,
    4 столика вмещают до 3 человек
 */
INSERT INTO tables (restaurant_id, seats_number)
VALUES (3, 8),
       (3, 8),
       (3, 3),
       (3, 3),
       (3, 3),
       (3, 3);