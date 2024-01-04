CREATE TABLE accounts (
    id VARCHAR(36) PRIMARY KEY NOT NULL,
    name VARCHAR NOT NULL,
    cpf VARCHAR UNIQUE NOT NULL,
    balance BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- CREATE TABLE IF NOT EXISTS accountstodo (
--     id SERIAL NOT NULL,
--     name VARCHAR(15) NOT NULL,
--     created_at TIMESTAMP NOT NULL,
--     updated_at TIMESTAMP NOT NULL,
--     PRIMARY KEY (id)
-- );

-- COMMENT ON COLUMN accountstodo.id IS 'アカウントID';
-- COMMENT ON COLUMN accountstodo.name IS '名前';
-- COMMENT ON COLUMN accountstodo.created_at IS '作成日時';
-- COMMENT ON COLUMN accountstodo.updated_at IS '更新日時';