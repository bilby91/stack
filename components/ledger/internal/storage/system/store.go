package system

import (
	"github.com/uptrace/bun"
)

const Schema = "_system"

type Store struct {
	db bun.IDB
}

func NewStore(db bun.IDB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) DB() bun.IDB {
	return s.db
}
