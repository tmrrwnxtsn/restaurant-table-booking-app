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

CREATE TABLE IF NOT EXISTS clients
(
    id    SERIAL PRIMARY KEY,
    name  VARCHAR(255) NOT NULL,
    phone VARCHAR(11)  NOT NULL
);

CREATE TABLE IF NOT EXISTS bookings
(
    id          SERIAL PRIMARY KEY,
    client_id   INTEGER   NOT NULL,
    booked_from TIMESTAMP NOT NULL,
    booked_to   TIMESTAMP NOT NULL,
    CONSTRAINT fk_bookings_clients FOREIGN KEY (client_id) REFERENCES clients (id)
);

CREATE TABLE IF NOT EXISTS bookings_tables
(
    id         SERIAL PRIMARY KEY,
    booking_id INTEGER NOT NULL,
    table_id   INTEGER NOT NULL,
    CONSTRAINT fk_bookings_tables_bookings FOREIGN KEY (booking_id) REFERENCES bookings (id) ON DELETE CASCADE,
    CONSTRAINT fk_bookings_tables_tables FOREIGN KEY (table_id) REFERENCES tables (id) ON DELETE CASCADE
);