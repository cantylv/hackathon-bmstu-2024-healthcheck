-------- EXTENSIONS --------
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE user_sex AS ENUM ('W', 'M');

-- НЕТ ФИЗИЧЕСКОЙ АКТИВНОСТИ - NFA (No Physical Activity)
-- ЛЕГКАЯ АКТИВНОСТЬ - LA (Light Activity)
-- СРЕДНЯЯ АКТИВНОСТЬ - MA (Moderate Activity)
-- ТЯЖЕЛАЯ АКТИВНОСТЬ - HA (Heavy Activity)
-- ЭКСТРЕМАЛЬНАЯ АКТИВНОСТЬ - EA (Extreme Activity)

CREATE TYPE user_activity AS ENUM ('NFA', 'LA', 'MA', 'HA', 'EA');

-------- DDL table 'user' --------
CREATE TABLE "user" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT,
    username TEXT,
    first_name TEXT,
    weight FLOAT,
    height INTEGER,
    age INTEGER,
    sex user_sex,
    physical_activity user_activity,
    password TEXT,
    day_calories FLOAT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

ALTER TABLE "user"
    ADD CONSTRAINT user_unique_email UNIQUE (email),
    ADD CONSTRAINT user_email_length CHECK (LENGTH(email) <= 50 AND LENGTH(email) >= 6),

    ADD CONSTRAINT user_unique_username UNIQUE (username),
    ADD CONSTRAINT user_username_length CHECK (LENGTH(username) <= 30 AND LENGTH(username) >= 2),

    ADD CONSTRAINT user_first_name_length CHECK (LENGTH(first_name) <= 30 AND LENGTH(first_name) >= 2);

ALTER TABLE "user"
    ALTER COLUMN username SET NOT NULL,
    ALTER COLUMN first_name SET NOT NULL,
    ALTER COLUMN weight SET NOT NULL,
    ALTER COLUMN height SET NOT NULL,
    ALTER COLUMN age SET NOT NULL,
    ALTER COLUMN sex SET NOT NULL,
    ALTER COLUMN physical_activity SET NOT NULL,
    ALTER COLUMN password SET NOT NULL,
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN updated_at SET NOT NULL;


-- -- Эта таблица содержит данные о группах
-- CREATE TABLE "group" (
--     id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--     name TEXT,
--     owner_id UUID REFERENCES "user"(id) ON DELETE RESTRICT,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
--     updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
-- );

-- CREATE TABLE agent (
--     id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--     name TEXT,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
-- );

-- CREATE TYPE status_type AS ENUM ('in_progress', 'rejected', 'approved');

-- -- Эта таблица содержит данные о заявках на создание группу
-- CREATE TABLE bid (
--     id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--     group_name TEXT,
--     user_id UUID REFERENCES "user"(id) ON DELETE CASCADE,
--     status status_type,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
--     updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
-- );

-- -- Эта таблица содержит доступ групп к агентам
-- CREATE TABLE privelege_group (
--     id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--     agent_id INT REFERENCES agent(id) ON DELETE CASCADE,
--     group_id INT REFERENCES "group"(id) ON DELETE CASCADE,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
-- );

-- -- Эта таблица содержит доступ пользователей к агентам
-- CREATE TABLE privelege_user (
--     id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--     agent_id INT REFERENCES agent(id) ON DELETE CASCADE,
--     user_id UUID REFERENCES "user"(id) ON DELETE CASCADE,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
-- );

-- -- Эта таблица содержит принадлежность пользователей к группам
-- CREATE TABLE participation (
--     id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--     user_id UUID REFERENCES "user"(id) ON DELETE CASCADE,
--     group_id INT REFERENCES "group"(id) ON DELETE CASCADE
-- );

-- -- table 'privelege_group'
-- ALTER TABLE privelege_group
-- ALTER COLUMN created_at SET NOT NULL;

-- -- table 'privelege_user'
-- ALTER TABLE privelege_user
-- ALTER COLUMN created_at SET NOT NULL;

-- -- table 'group'
-- ALTER TABLE "group"
-- ADD CONSTRAINT group_unique_name UNIQUE (name),
-- ADD CONSTRAINT group_name_length CHECK (LENGTH(name)>=2 AND LENGTH(name) <= 30);

-- ALTER TABLE "group"
-- ALTER COLUMN name SET NOT NULL,
-- ALTER COLUMN created_at SET NOT NULL,
-- ALTER COLUMN updated_at SET NOT NULL;

-- -- table 'agent'
-- ALTER TABLE agent
-- ADD CONSTRAINT agent_unique_name UNIQUE (name),
-- ADD CONSTRAINT agent_name_length CHECK (LENGTH(name)>=2 AND LENGTH(name) <= 50);

-- ALTER TABLE "group"
-- ALTER COLUMN name SET NOT NULL,
-- ALTER COLUMN created_at SET NOT NULL;


-- -- table 'bid'
-- ALTER TABLE bid
-- ALTER COLUMN group_name SET NOT NULL,
-- ALTER COLUMN status SET NOT NULL,
-- ALTER COLUMN created_at SET NOT NULL,
-- ALTER COLUMN updated_at SET NOT NULL;

-------- FUNCTIONS AND TRIGGERS --------
-- table 'user'
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_user_updated_at
BEFORE UPDATE ON "user"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- CREATE TRIGGER update_group_updated_at
-- BEFORE UPDATE ON "group"
-- FOR EACH ROW
-- EXECUTE FUNCTION update_updated_at_column();


