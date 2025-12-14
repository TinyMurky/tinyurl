// Package database is a package for url shortener
package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/TinyMurky/snowflake"

	"github.com/TinyMurky/tinyurl/internal/urlshortener/model"
	"github.com/TinyMurky/tinyurl/pkg/database"
)

type URLShortenerDB struct {
	db *database.DB
}

func New(db *database.DB) *URLShortenerDB {
	return &URLShortenerDB{
		db: db,
	}
}

// GetFirstByID will get first url by sid
func (db *URLShortenerDB) GetFirstByID(ctx context.Context, sid snowflake.SID) (model.URL, error) {
	var urlFromDB model.URL
	query := `
		SELECT id, long_url, created_at
		FROM urls
		WHERE id = ?
		LIMIT 1;
	`

	row := db.db.Pool.QueryRowContext(ctx, query, int64(sid))

	err := row.Scan(
		&urlFromDB.ID,
		&urlFromDB.LongURL,
		&urlFromDB.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.URL{}, nil
		}
		return model.URL{}, fmt.Errorf("GetFirstByID scan error: %w", err)
	}

	return urlFromDB, nil
}

// GetFirstByLongURL will get first url by longURL
func (db *URLShortenerDB) GetFirstByLongURL(ctx context.Context, longURL string) (model.URL, error) {
	var urlFromDB model.URL
	query := `
		SELECT id, long_url, created_at
		FROM urls
		WHERE long_url = ?
		LIMIT 1;
	`

	row := db.db.Pool.QueryRowContext(ctx, query, longURL)

	err := row.Scan(
		&urlFromDB.ID,
		&urlFromDB.LongURL,
		&urlFromDB.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.URL{}, nil
		}
		return model.URL{}, fmt.Errorf("GetFirstByID scan error: %w", err)
	}

	return urlFromDB, nil
}

func (db *URLShortenerDB) CreateURL(ctx context.Context, u model.URL) error {
	if u.ID == 0 {
		return errors.New("create URL need to provide ID")
	}

	if u.LongURL == "" {
		return errors.New("create url need to provide longURL")
	}

	query := `
        INSERT INTO urls (id, long_url)
        VALUES (?, ?)
    `

	_, err := db.db.Pool.ExecContext(ctx, query, int64(u.ID), u.LongURL)

	if err != nil {
		return fmt.Errorf("create url error: %w", err)
	}
	return nil
}

// GetFirstByID will get first url by sid
// return URL in zero value if not found
// func (db *URLShortenerDB) GetFirstByID(ctx context.Context, sid snowflake.SID) (model.URL, error) {
// 	var urlFromDB model.URL
// 	id := int64(sid)
//
// 	query := `
// 		SELECT id, long_url, created_at
// 		FROM urls
// 		WHERE id = ?
// 		LIMIT 1;
// 	`
//
// 	if err := db.db.InTx(
// 		ctx,
// 		&sql.TxOptions{
// 			Isolation: sql.LevelReadCommitted,
// 			ReadOnly:  true,
// 		},
// 		func(tx *sql.Tx) error {
// 			stmt, err := tx.Prepare(query)
//
// 			if err != nil {
// 				return err
// 			}
//
// 			rows, err := stmt.QueryContext(ctx, id)
//
// 			if err != nil {
// 				return err
// 			}
//
// 			if rows.Next() {
// 				if err := rows.Scan(
// 					&urlFromDB.ID,
// 					&urlFromDB.LongURL,
// 					&urlFromDB.CreatedAt,
// 				); err != nil {
// 					return fmt.Errorf("rows scan: %w", err)
// 				}
// 			}
//
// 			if err := rows.Err(); err != nil {
// 				return fmt.Errorf("rows iteration: %w", err)
// 			}
//
// 			return nil
// 		}); err != nil {
//
// 		return model.URL{}, fmt.Errorf("get first url: %w", err)
// 	}
//
// 	return urlFromDB, nil
// }
