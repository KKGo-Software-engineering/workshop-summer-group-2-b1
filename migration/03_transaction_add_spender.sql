-- +goose Up
-- +goose StatementBegin
ALTER TABLE "transaction" ADD COLUMN "spender_id" int4;
-- +goose StatementEnd
