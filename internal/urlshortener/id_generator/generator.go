// Package idgenerator can generate snowflake ID
package idgenerator

import (
	"fmt"
	"time"

	"github.com/TinyMurky/snowflake"

	urlshortenerconfig "github.com/TinyMurky/tinyurl/internal/urlshortener/config"
)

const epochTimeFormat = "2006-01-02"

// Generator will generate snowflake ID
type Generator struct {
	generator *snowflake.Generator
	nodeID    int64
}

// NewGenerator will create a new snowflakeID generator
func NewGenerator(cfg *urlshortenerconfig.Config) (*Generator, error) {
	nodeID := cfg.IDGenerator.NodeID
	epochTimeStartFrom := cfg.IDGenerator.EpochTimeStartFrom

	epochStartDate, err := time.Parse(epochTimeFormat, epochTimeStartFrom)

	if err != nil {
		return nil, fmt.Errorf("invalid epoch time format, required: %s, got: %s", epochTimeFormat, epochTimeStartFrom)
	}

	epochTime := time.Date(
		epochStartDate.Year(),
		epochStartDate.Month(),
		epochStartDate.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)

	generator, err := snowflake.NewGenerator(
		snowflake.WithEpoch(epochTime),
	)

	if err != nil {
		return nil, fmt.Errorf("new snowflake generator error: %w", err)
	}

	return &Generator{
		generator: generator,
		nodeID:    nodeID,
	}, nil
}

// NextID will get 1 snowflakeID
func (g *Generator) NextID() (snowflake.SID, error) {
	return g.generator.NextID(g.nodeID)
}
