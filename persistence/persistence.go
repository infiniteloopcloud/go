package persistence

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

// Storable defines whether a struct can be persisted to database or not
type Storable interface {
	Create(transactionID string) (DatabaseCreate, error)
}

// Updatable defines whether a struct can be updated if it's already exists in the database
type Updatable interface {
	Update(transactionID string) (DatabaseUpdate, error)
}

// Selectable defines whether a struct can be retrieved from the database
type Selectable interface {
	Select(opts ...interface{}) (DatabaseSelect, error)
}

// Prepared is a definition for anything implement the PrepareContext. This can be a single database statement or a
// database transaction statement as well.
type Prepared interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

// Query is a definition for anything implement the PrepareContext. This can be a single database statement or a
// database transaction statement as well.
type Query interface {
	QueryContext(ctx context.Context, query string, elements ...interface{}) (*sql.Rows, error)
}

// DatabaseCreate is the result of a Storable call
type DatabaseCreate struct {
	Fields           []interface{}
	Values           []interface{}
	ReturningColumns []interface{}
	Receivers        []interface{}
	Table            string
}

// DatabaseUpdate is the result of a Updatable call
type DatabaseUpdate struct {
	Records map[string]interface{}
	Where   []exp.Expression
	Table   string
}

// DatabaseSelect is the result of a Selectable call
type DatabaseSelect struct {
	Fields    []interface{}
	Receivers []interface{}
	Table     string
}

// Create is a general solution for adding values to a Prepared implementation
func Create(ctx context.Context, p Prepared, dialect goqu.DialectWrapper, dc DatabaseCreate) error {
	s, params, err := dialect.Insert(dc.Table).Prepared(true).Cols(dc.Fields...).Vals(dc.Values).
		Returning(dc.ReturningColumns...).ToSQL()
	if err != nil {
		return err
	}
	stmt, err := p.PrepareContext(ctx, s)
	if err != nil {
		return err
	}
	if dc.Receivers != nil {
		err = stmt.QueryRowContext(ctx, params...).Scan(dc.Receivers...)
	} else {
		_, err = stmt.ExecContext(ctx, params...)
	}
	if err != nil {
		return err
	}
	return nil
}

// Update is a general solution for updating values to a Prepared implementation
func Update(ctx context.Context, p Prepared, dialect goqu.DialectWrapper, du DatabaseUpdate) error {
	s, params, err := dialect.Update(du.Table).Prepared(true).Set(goqu.Record(du.Records)).
		Where(du.Where...).ToSQL()
	if err != nil {
		return err
	}
	stmt, err := p.PrepareContext(ctx, s)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, params...)
	if err != nil {
		return err
	}
	return nil
}

// Select is a general solution for getting data from a Query implementation
func Select(ctx context.Context, q Query, dialect goqu.DialectWrapper, ds DatabaseSelect, where ...exp.Expression) error {
	s, params, err := dialect.From(ds.Table).Select(ds.Fields...).Prepared(true).Where(where...).ToSQL()
	if err != nil {
		return err
	}
	rows, err := q.QueryContext(ctx, s, params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(ds.Receivers...); err != nil {
			return err
		}
	}

	return err
}
