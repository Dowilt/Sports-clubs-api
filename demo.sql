CREATE TABLE cities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE clubs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    city_id INTEGER NOT NULL REFERENCES cities(id),
    titles_count INTEGER NOT NULL CHECK (titles_count >= 0),
    avg_age INTEGER NOT NULL CHECK (avg_age > 0 AND avg_age < 100)
);

INSERT INTO cities (name) VALUES
('Москва'),
('Санкт-Петербург'),
('Екатеринбург'),
('Казань'),
('Сочи'),
('Новосибирск'),
('Ростов-на-Дону'),
('Краснодар'),
('Воронеж'),
('Пермь'),
('Омск'),
('Уфа'),
('Челябинск'),
('Волгоград'),
('Калининград');

INSERT INTO clubs (name, city_id, titles_count, avg_age) VALUES
('Динамо', 1, 12, 26),
('Зенит', 2, 9, 28),
('Урал', 3, 2, 24),
('Рубин', 4, 6, 25),
('Краснодар', 8, 3, 23),
('ЦСКА', 1, 15, 27),
('Спартак', 1, 10, 25),
('Локомотив', 1, 8, 26),
('Ахмат', 5, 1, 24),
('Ротор', 6, 2, 29),
('Факел', 7, 0, 22),
('Торпедо', 1, 5, 23),
('Уфа', 12, 0, 21),
('Химки', 1, 0, 25),
('Ростов', 7, 1, 24);