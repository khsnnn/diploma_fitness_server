-- Создание таблиц для базы данных спортивных клубов

CREATE TABLE clubs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    description TEXT,
    working_hours VARCHAR(255),
    rating FLOAT DEFAULT 0.0,
    coordinates POINT,
    type VARCHAR(50) CHECK (type IN ('commercial', 'university')) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE subcategories (
    id SERIAL PRIMARY KEY,
    category_id INTEGER REFERENCES categories(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    UNIQUE (category_id, name)
);

CREATE TABLE club_categories (
    club_id INTEGER REFERENCES clubs(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (club_id, category_id)
);

CREATE TABLE club_subcategories (
    club_id INTEGER REFERENCES clubs(id) ON DELETE CASCADE,
    subcategory_id INTEGER REFERENCES subcategories(id) ON DELETE CASCADE,
    PRIMARY KEY (club_id, subcategory_id)
);

CREATE TABLE schedules (
    id SERIAL PRIMARY KEY,
    club_id INTEGER REFERENCES clubs(id) ON DELETE CASCADE,
    day_of_week VARCHAR(20) CHECK (day_of_week IN ('Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday')) NOT NULL,
    time VARCHAR(50) NOT NULL,
    activity VARCHAR(255) NOT NULL,
    instructor VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_clubs_name ON clubs(name);
CREATE INDEX idx_clubs_coordinates ON clubs USING GIST (coordinates);
CREATE INDEX idx_clubs_type ON clubs(type);
CREATE INDEX idx_schedules_club_id ON schedules(club_id);
CREATE INDEX idx_schedules_day_of_week ON schedules(day_of_week);