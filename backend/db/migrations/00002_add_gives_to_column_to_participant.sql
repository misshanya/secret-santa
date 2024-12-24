-- +goose Up
-- +goose StatementBegin
ALTER TABLE participants
ADD gives_to BIGINT REFERENCES participants(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE participants
DROP COLUMN gives_to;
-- +goose StatementEnd
