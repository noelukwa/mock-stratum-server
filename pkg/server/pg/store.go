package pg

import (
	"context"

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

func (s *Store) SaveAuthRequest(method, req_id, id, user_id string) error {

	query, args, err := s.gen_save_query(method, req_id, id, user_id)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) SaveSubRequest(method, agent, id, req_id, extra_nonce string) error {
	ctx := context.Background()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return err
	}

	query1, args1, _ := s.gen_sub_query(method, agent, id, req_id, extra_nonce)

	_, err = tx.ExecContext(ctx, query1, args1...)
	if err != nil {
		tx.Rollback()
		return err
	}

	query2, args2, _ := s.gen_save_query(method, req_id, id, "")

	_, err = tx.ExecContext(ctx, query2, args2...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) gen_save_query(method, req_id, id, user_id string) (query string, args []interface{}, err error) {
	return s.sq.Insert("requests").
		Columns("method", "req_id", "id", "user_id").
		Values(method, req_id, id, user_id).ToSql()
}

func (s *Store) gen_sub_query(method, agent, id, req_id, extra_nonce string) (query string, args []interface{}, err error) {
	return s.sq.Insert("subscriptions").
		Columns("method", "user_agent", "id", "req_id", "extra_nonce").
		Values(method, agent, id, req_id, extra_nonce).ToSql()
}
