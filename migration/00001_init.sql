-- +goose Up
-- +goose StatementBegin
CREATE COLLATION case_insensitive (
    PROVIDER = ICU, LOCALE = 'und-u-ks-level2', DETERMINISTIC = false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP COLLATION IF EXISTS case_insensitive;
-- +goose StatementEnd
