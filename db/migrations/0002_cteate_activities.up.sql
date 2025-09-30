-- таблица активностей пользователя
CREATE TABLE IF NOT EXISTS activities (
  id           BIGSERIAL PRIMARY KEY,
  user_id      BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name         TEXT        NOT NULL,
  emoji        TEXT        NULL,
  is_archived  BOOLEAN     NOT NULL DEFAULT FALSE,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- уникальность имени активности в рамках пользователя (без учета регистра)
CREATE UNIQUE INDEX IF NOT EXISTS uq_activities_user_lower_name
  ON activities (user_id, lower(name));

-- (полезно для выборок по пользователю)
CREATE INDEX IF NOT EXISTS idx_activities_user
  ON activities (user_id);







