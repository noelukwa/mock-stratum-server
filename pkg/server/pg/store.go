package pg

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
	sq squirrel.StatementBuilderType
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
		sq: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// func (s *Store) SaveRequest(req *Request) error {
// 	_, err := s.sq.Insert("requests").
// 		Columns("id", "method", "params", "id").
// 		Values(req.Id, req.Method, req.Params, req.Id).
// 		Exec(s.db)
// 	return err
// }
