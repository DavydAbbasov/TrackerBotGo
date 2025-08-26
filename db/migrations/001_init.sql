-- === СХЕМА ДЛЯ tracker_bot ===

-- пользователи
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    tg_id BIGINT UNIQUE NOT NULL,
    locale TEXT NOT NULL DEFAULT 'en',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- активности пользователя
CREATE TABLE IF NOT EXISTS activities (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    is_archived BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, name)
);

-- сессии активности (старт/стоп)
CREATE TABLE IF NOT EXISTS sessions (
    id BIGSERIAL PRIMARY KEY,
    activity_id BIGINT NOT NULL REFERENCES activities (id) ON DELETE CASCADE,
    started_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ,
    CHECK (
        ended_at IS NULL
        OR ended_at >= started_at
    )
);

-- частый фильтр: по активности и времени
CREATE INDEX IF NOT EXISTS idx_sessions_activity_started ON sessions (activity_id, started_at);

-- === ПРАВА ДЛЯ РОЛИ tracker (твоё приложение) ===
GRANT
SELECT, INSERT,
UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO tracker;

GRANT USAGE,
SELECT,
UPDATE ON ALL SEQUENCES IN SCHEMA public TO tracker;

-- чтобы будущие таблицы/последовательности сразу имели права
ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT
SELECT, INSERT,
UPDATE, DELETE ON TABLES TO tracker;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT USAGE,
SELECT,
UPDATE ON SEQUENCES TO tracker;