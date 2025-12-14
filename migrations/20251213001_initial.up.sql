-- BEGIN;
    CREATE TABLE IF NOT EXISTS urls (
        -- id 設定為 INTEGER PRIMARY KEY
        -- 在 SQLite 中，INTEGER 可以直接儲存 64-bit (int64) 的數值，
        -- 這足以容納雪花演算法產生的 ID。
        id INTEGER PRIMARY KEY,

        long_url TEXT NOT NULL UNIQUE,

        -- (選用) 建議加上建立時間，方便日後維護
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
-- COMMIT;
