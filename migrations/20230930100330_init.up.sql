-- Создаем таблицы
CREATE TABLE IF NOT EXISTS ns_roles (
	id SERIAL PRIMARY KEY,
	role VARCHAR(32) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS ns_users (
	id UUID PRIMARY KEY,
	login VARCHAR(32) NOT NULL UNIQUE,
	password VARCHAR(128) NOT NULL,
	salt VARCHAR(8) NOT NULL,
	role VARCHAR(32) DEFAULT 'user' REFERENCES ns_roles ("role") ON UPDATE CASCADE ON DELETE CASCADE,
	registered_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	last_login_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ns_notes (
	id UUID PRIMARY KEY,
	author_id UUID REFERENCES ns_users ("id") ON UPDATE CASCADE ON DELETE CASCADE,
	title VARCHAR(128) NOT NULL,
	content TEXT,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ns_notes_config (
	id SERIAL PRIMARY KEY,
	user_id UUID REFERENCES ns_users ("id") ON UPDATE CASCADE ON DELETE CASCADE,
	notes_limit INT NOT NULL DEFAULT 30
);

-- Добавляем триггер на добавление пользователя в конфиг с заметками
CREATE OR REPLACE FUNCTION add_user_to_notes_config() RETURNS TRIGGER AS $$
BEGIN
    -- Добавляем нового пользователя в ns_notes_config с лимитом по умолчанию
    INSERT INTO ns_notes_config (user_id) VALUES (NEW.id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_add_user_to_notes_config
AFTER INSERT ON ns_users
FOR EACH ROW EXECUTE PROCEDURE add_user_to_notes_config();

-- Добавляем триггер на проверку лимита заметок
CREATE OR REPLACE FUNCTION check_notes_limit()
RETURNS TRIGGER AS $$
DECLARE
    limit_count INT;
    current_count INT;
BEGIN
    -- Получаем лимит заметок для пользователя
    SELECT notes_limit INTO limit_count 
    FROM ns_notes_config 
    WHERE user_id = NEW.author_id;
    
    -- Если лимит не установлен, то прекращаем выполнение триггера
    IF limit_count IS NULL THEN
        RAISE NOTICE 'limit is not set for user %', NEW.author_id;
        RETURN NEW;
    END IF;
    
    -- Получаем текущее количество заметок пользователя
    SELECT COUNT(*) INTO current_count 
    FROM ns_notes 
    WHERE author_id = NEW.author_id;
    
    -- Если текущее количество заметок равно или превышает лимит, то отклоняем вставку
    IF current_count >= limit_count THEN
        RAISE EXCEPTION 'notes limit exceeded for user %', NEW.author_id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_check_notes_limit
BEFORE INSERT ON ns_notes
FOR EACH ROW EXECUTE PROCEDURE check_notes_limit();


-- Вставляем роли по умолчанию
INSERT INTO ns_roles ("role") VALUES ('admin'); 
INSERT INTO ns_roles ("role") VALUES ('user'); 

-- Вставляем пользователя по умолчанию
INSERT INTO ns_users ("id", "login", "password", "salt", "role")
VALUES ('07f3c5a1-70ea-4e3f-b9b5-110d29891673', 'admin', 'i2RekjynLOONpUgaxpHj3Y0/DmYuhp24THRkAnt4AyX/MHqY8wEfYIIZdK3WKks1cLXEU8rCEax2JHrParZkZg==', 'jCrzpcyc', 'admin')