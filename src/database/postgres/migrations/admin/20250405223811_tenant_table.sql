-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tenants.tenants (
  id uuid PRIMARY KEY,
  name varchar NOT NULL,
  db_host varchar NOT NULL,
  db_port int NOT NULL,
  db_name varchar(63) UNIQUE NOT NULL,
  db_user varchar NOT NULL,
  db_password bytea NOT NULL,
  iv bytea NOT NULL,
  version varchar NOT NULL,
  is_active boolean DEFAULT true,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tenants.tenants;
-- +goose StatementEnd
