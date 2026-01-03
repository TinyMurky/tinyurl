package database

import "testing"

func TestToDSN(t *testing.T) {
	testCases := []struct {
		wantedDsn string
		config    Config
	}{
		{
			wantedDsn: "file://home?_busy_timeout=5000&_cache_size=-2000&_foreign_keys=on&_journal_mode=WAL&_synchronous=NORMAL",
			config: Config{
				Path:        "home",
				JournalMode: "WAL",
				BusyTimeout: 5000,     // 5秒
				SyncMode:    "NORMAL", // 在 WAL 模式下夠安全且快
				ForeignKeys: true,     // 強制檢查關聯
				CacheSize:   -2000,    // 預設 cache 大小
			},
		},
		{
			wantedDsn: "file://home?_busy_timeout=5000&_foreign_keys=off&_journal_mode=WAL",
			config: Config{
				Path:        "home",
				JournalMode: "WAL",
				BusyTimeout: 5000,  // 5秒
				SyncMode:    "",    // 在 WAL 模式下夠安全且快
				ForeignKeys: false, // 強制檢查關聯
				CacheSize:   0,     // 預設 cache 大小
			},
		},
	}

	for _, tc := range testCases {
		got := tc.config.ToFileDSN()
		want := tc.wantedDsn

		if got != want {
			t.Errorf("Expect %q, got %q", want, got)
		}

	}
}
