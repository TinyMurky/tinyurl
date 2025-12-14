// Copyright 2020 the Exposure Notifications Server authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package database

import (
	"fmt"
	"net/url"
	"path/filepath"
	"reflect"
)

// Config 封裝了連線所需的參數
// Most configs supported by "github.com/mattn/go-sqlite3"
// are also supported by "modernc.org/sqlite"
type Config struct {
	Path        string `env:"DB_PATH"`                                              // 資料庫檔案路徑 (例如: "./data.db")
	JournalMode string `env:"DB_JOURNAL_MODE, default=WAL" sqlite:"_journal_mode"`  // 建議: "WAL"
	BusyTimeout int    `env:"DB_BUSY_TIMEOUT, default=5000" sqlite:"_busy_timeout"` // 建議: 5000 (毫秒)
	SyncMode    string `env:"DB_SYNC_MODE, default=NORMAL" sqlite:"_synchronous"`   // 建議: "NORMAL"
	ForeignKeys bool   `env:"DB_FOREIGN_KEYS, default=true" sqlite:"_foreign_keys"` // 建議: true
	CacheSize   int    `env:"DB_CACHE_SIZE, default=-2000" sqlite:"_cache_size"`    // 建議: -2000 (代表約 2MB)
}

// DefaultConfig 回傳一組建議的預設值
// func DefaultConfig(path string) *Config {
// 	return &Config{
// 		Path:        path,
// 		JournalMode: "WAL",
// 		BusyTimeout: 5000,
// 		SyncMode:    "NORMAL",
// 		ForeignKeys: true,
// 		CacheSize:   -2000,
// 	}
// }

// DatabaseConfig return the config of database
func (c *Config) DatabaseConfig() *Config {
	return c
}

// ToFileDSN transfer config to file DSN
func (c *Config) ToFileDSN() string {
	return c.ToDSN("file")
}

// ToSQLiteDSN transfer config to sqlite DSN
func (c *Config) ToSQLiteDSN() string {
	return c.ToDSN("sqlite")
}

// ToDSN transfer config to  DSN
// driver can be "sqlite" or "file"
func (c *Config) ToDSN(driver string) string {
	if c == nil || c.Path == "" {
		return ""
	}

	val := reflect.ValueOf(c)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	t := val.Type()
	query := make(url.Values)

	for i := 0; i < t.NumField(); i++ {
		fieldDef := t.Field(i)
		fieldVal := val.Field(i)

		// 1. 略過 Path 欄位 (它不在 Query String 中)
		if fieldDef.Name == "Path" {
			continue
		}

		// 2. 取得 Tag
		sqliteTag, ok := fieldDef.Tag.Lookup("sqlite")
		if !ok {
			continue
		}

		// 3. 處理數值轉換
		var stringValue string

		switch fieldVal.Kind() {
		case reflect.Bool:
			// 優化 1: 將 bool 轉為 "on" / "off" (相容性更好)
			if fieldVal.Bool() {
				stringValue = "on"
			} else {
				// 優化 2: 解決 IsZero 問題
				// 如果你是 false，我們明確傳送 "off"，而不是略過
				// 這樣可以確保設定被強制覆蓋
				stringValue = "off"
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if fieldVal.Int() == 0 {
				continue // 數值為 0 時略過 (依需求調整)
			}
			stringValue = fmt.Sprintf("%d", fieldVal.Int())
		case reflect.String:
			if fieldVal.String() == "" {
				continue
			}
			stringValue = fieldVal.String()
		default:
			// 其他類型回退到預設格式
			if fieldVal.IsZero() {
				continue
			}
			stringValue = fmt.Sprintf("%v", fieldVal.Interface())
		}

		query.Add(sqliteTag, stringValue)
	}

	// 優化 3: 確保路徑分隔符統一 (處理 Windows 路徑問題)
	cleanPath := filepath.ToSlash(c.Path)

	dsn := fmt.Sprintf("%s://%s?%s", driver, cleanPath, query.Encode())
	return dsn
}
