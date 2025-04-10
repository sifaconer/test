-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS tenants;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS tenants;
-- +goose StatementEnd
