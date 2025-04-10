-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tenants.users_directory (
  id uuid PRIMARY KEY,
  email varchar UNIQUE NOT NULL,
  name varchar NOT NULL,
  password varchar NOT NULL,
  is_active boolean DEFAULT true,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tenants.users_directory;
-- +goose StatementEnd
