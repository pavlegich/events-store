-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS events (
    eventID Int64,
    eventType String,
    userID Int64,
    eventTime DateTime,
    payload String
) ENGINE = MergeTree
ORDER BY (eventID, eventTime);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE events;