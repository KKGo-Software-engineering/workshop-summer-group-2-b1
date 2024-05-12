-- +goose Up
-- +goose StatementBegin
ALTER TABLE "transaction"
ADD COLUMN "spender_id" int4;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE DROP COLUMN "spender_id";
-- +goose StatementEnd
