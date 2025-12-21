# 1. Design

> [!Note]
> - folder structure follow [google/exposure-notifications-server](https://github.com/google/exposure-notifications-server)
> - [The standard library now has all you need for advanced routing in Go](https://www.youtube.com/watch?v=H7tbjKFSg58)

## 1.1 API

Four APIs as follow:
- `POST /api/v1/data/shorten`:
    - with body `{longUrl: longURLString}`
    - return: shortURL
- `GET /api/v1/shortUrl`:
    - return status 302
    - return longUrl for redirect
- `GET /shortUrl`: redirect to `GET /api/v1/shortUrl`
- `GET /`: UI

## 1.2 

- `POST /api/v1/data/shorten`:
    - when longURL comes in, check if usl `not exists in` database
        - if not exists in => use snowflake ID to generate Index and insert => return shortURL
        - if might exist => search longURL
            - if exist url => return shortURl
            - not exist => insert

- `GET /api/v1/shortUrl`:
    - Plan 1: use bloom filter to check wether url exist in database first
    - Plan 2: use catch
    - Plan 3: do nothing

# Note

Each api should:
1. check http method
2. check status
3. check input

Middleware should:
1. log input
2. recovery

# How to route net/http

[![The standard library now has all you need for advanced routing in Go.](https://img.youtube.com/vi/H7tbjKFSg58/0.jpg)](https://www.youtube.com/watch?v=H7tbjKFSg58)

# 未來流程規劃

```mermaid
graph TD
    User((使用者)) --> LB[Load Balancer]
    LB --> Server[Web Servers]

    subgraph "Write Path（產生短網址)"
        Server --> Blacklist{1. 黑名單檢查<br/>Redis/Local Cache}
        Blacklist -- 命中黑名單 --> Refuse[回傳 403 Forbidden]
        Blacklist -- 安全 --> KeyGen[2. 產生短網址 Key]
        KeyGen --> WriteDB[(3. 寫入 Database)]
        WriteDB --> UpdateBF[4. 更新 Bloom Filter]
        UpdateBF --> UpdateCache[5. 寫入 Cache]
    end

    subgraph "Read Path(100000000/24/60/60還原網址)"
        Server --> BF_Check{A. Bloom Filter<br/>檢查是否存在?}
        BF_Check -- 絕對不存在 --> 404[直接回傳 404]
        BF_Check -- 可能存在 --> CacheCheck{B. Cache 查詢<br/>Redis}
        
        CacheCheck -- Hit --> Redirect[回傳 302/301 導向]
        CacheCheck -- Miss --> DB_Query[(C. Database 查詢)]
        
        DB_Query -- 存在 --> SyncCache[更新 Cache]
        SyncCache --> Redirect
        DB_Query -- 不存在 --> 404_Real[回傳 404]
    end

    %% 零件定義
    subgraph "資料層"
        Redis[(Redis<br/>Cache + BloomFilter)]
        DB[(PostgreSQL/MySQL)]
    end

    %% 連接資料層
    UpdateCache -.-> Redis
    CacheCheck -.-> Redis
    WriteDB -.-> DB
    DB_Query -.-> DB
    UpdateBF -.-> Redis
```
