-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tenants.user_tenants (
  id uuid PRIMARY KEY,
  user_id uuid NOT NULL,
  tenant_id uuid NOT NULL,
  default_tenant boolean DEFAULT false,
  created_at timestamptz DEFAULT now(),
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES tenants.users_directory (id) ON DELETE CASCADE,
  CONSTRAINT fk_tenant FOREIGN KEY (tenant_id) REFERENCES tenants.tenants (id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tenants.user_tenants;
-- +goose StatementEnd
