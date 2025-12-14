package urlshortenerconfig

type IDGeneratorConfig struct {
	NodeID int64 `env:"ID_GEN_NODE_ID, default=1"`

	// It should be YYYY-MM-DD
	EpochTimeStartFrom string `env:"ID_GEN_EPOCH_TIME_START_FROM, default=2025-12-14"`
}
