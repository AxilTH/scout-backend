-- Таблица пользователей
CREATE TABLE users (
   id BIGSERIAL PRIMARY KEY,
   first_name TEXT NOT NULL,
   last_name TEXT NOT NULL,
   middle_name TEXT,
   phone_number TEXT,
   email TEXT UNIQUE NOT NULL,
   date_of_birth DATE,
   vk_profile_url TEXT,
   password_hash TEXT NOT NULL,
   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Таблица учебных заведений
CREATE TABLE education_institutions (
   id BIGSERIAL PRIMARY KEY,
   title TEXT NOT NULL,
   short_name TEXT,
   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Таблица приглашений
CREATE TABLE invitations (
   id BIGSERIAL PRIMARY KEY,
   email TEXT NOT NULL,
   squad_id BIGINT NOT NULL,
   expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
   used_at TIMESTAMP WITH TIME ZONE,
   used_by BIGINT,
   created_by BIGINT NOT NULL,
   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Внешние ключи
ALTER TABLE invitations
   ADD CONSTRAINT fk_invitations_created_by
   FOREIGN KEY (created_by) REFERENCES users(id);