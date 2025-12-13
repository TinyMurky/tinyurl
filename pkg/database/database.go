// Copyright 2020 the Exposure Notifications Server authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package database manage SQLite3,
// use https://gitlab.com/cznic/sqlite ("modernc.org/sqlite") as driver
package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	// ErrNotFound indicates that the requested record was not found in the database.
	ErrNotFound = errors.New("record not found")

	// ErrKeyConflict indicates that there was a key conflict inserting a row.
	ErrKeyConflict = errors.New("key conflict")
)

// InTx runs the given function f within a transaction with the provided
// sql TxOption.
func (db *DB) InTx(ctx context.Context, opts *sql.TxOptions, f func(tx *sql.Tx) error) error {
	conn, err := db.Pool.Conn(ctx)

	if err != nil {
		return fmt.Errorf("acquiring connection: %w", err)
	}

	defer conn.Close()

	tx, err := conn.BeginTx(ctx, opts)

	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	if err := f(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("rolling back transaction: %v (original error: %w)", rollbackErr, err)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
