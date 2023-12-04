// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: words.sql

package parser

import (
	"context"
)

const createWord = `-- name: CreateWord :one
INSERT INTO words (
    eng, rus, learn_at
) VALUES (
    $1, $2, CURRENT_TIMESTAMP
) ON CONFLICT (eng) DO UPDATE
SET rus = EXCLUDED.rus || ' ; ' || words.rus
RETURNING id, eng, rus, learn_at
`

type CreateWordParams struct {
	Eng string
	Rus string
}

func (q *Queries) CreateWord(ctx context.Context, arg CreateWordParams) (Word, error) {
	row := q.db.QueryRow(ctx, createWord, arg.Eng, arg.Rus)
	var i Word
	err := row.Scan(
		&i.ID,
		&i.Eng,
		&i.Rus,
		&i.LearnAt,
	)
	return i, err
}

const listWords = `-- name: ListWords :many
SELECT id, eng, rus, learn_at FROM words
`

func (q *Queries) ListWords(ctx context.Context) ([]Word, error) {
	rows, err := q.db.Query(ctx, listWords)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Word
	for rows.Next() {
		var i Word
		if err := rows.Scan(
			&i.ID,
			&i.Eng,
			&i.Rus,
			&i.LearnAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
