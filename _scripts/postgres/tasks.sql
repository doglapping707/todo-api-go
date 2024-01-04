-- テーブルを作成する
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL NOT NULL,
    title VARCHAR(15) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

-- コメントを設定する
COMMENT ON COLUMN tasks.id IS 'タスクID';
COMMENT ON COLUMN tasks.title IS 'タイトル';
COMMENT ON COLUMN tasks.created_at IS '作成日時';
COMMENT ON COLUMN tasks.updated_at IS '更新日時';

-- 関数を作成する
CREATE FUNCTION trigger_set_timestamp() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- トリガーを作成する
CREATE TRIGGER set_timestamp BEFORE UPDATE ON tasks FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

-- ダミーデータをインサートする
INSERT INTO tasks (
    title
) 
VALUES 
(
    'task1'
), 
(
    'task2'
), 
(
    'task3'
);

-- account_id INTEGER NOT NULL,
-- COMMENT ON COLUMN tasks.account_id IS 'アカウントID';

-- completed INTEGER CHECK(id = 0 OR id = 1) NOT NULL DEFAULT 0,
-- COMMENT ON COLUMN tasks.completed IS '完了フラグ 0:未完了 1:完了';