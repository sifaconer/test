-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.productos ADD created_at timestamptz DEFAULT now() NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.productos DROP COLUMN created_at;
-- +goose StatementEnd
