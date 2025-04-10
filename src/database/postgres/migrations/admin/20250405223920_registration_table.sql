-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tenants.registration (
  id uuid PRIMARY KEY,
  email varchar NOT NULL,
  name varchar NOT NULL,
  partial_data jsonb,
  user_id uuid NOT NULL,
  tenant_id uuid NOT NULL,
  is_completed boolean DEFAULT false,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  CONSTRAINT fk_reg_user FOREIGN KEY (user_id) REFERENCES tenants.users_directory (id) ON DELETE CASCADE,
  CONSTRAINT fk_reg_tenant FOREIGN KEY (tenant_id) REFERENCES tenants.tenants (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tenants.registration;
-- +goose StatementEnd
