CREATE TABLE transfers (
    id VARCHAR(36) PRIMARY KEY NOT NULL,
    account_origin_id VARCHAR NOT NULL,
    account_destination_id VARCHAR NOT NULL,
    amount BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL
);