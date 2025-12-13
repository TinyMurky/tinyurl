-- BEGIN;
    CREATE TABLE IF NOT EXISTS urls (
        -- id 設定為 INTEGER PRIMARY KEY
        -- 在 SQLite 中，INTEGER 可以直接儲存 64-bit (int64) 的數值，
        -- 這足以容納雪花演算法產生的 ID。
        id INTEGER PRIMARY KEY,

        -- shortURL 通常需要唯一性索引，以便快速查詢且不重複
        shortURL TEXT NOT NULL UNIQUE,

        -- longURL 存放原本的長網址
        longURL TEXT NOT NULL,

        -- (選用) 建議加上建立時間，方便日後維護
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
-- COMMIT;
