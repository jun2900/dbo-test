// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dal

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

var (
	Q        = new(Query)
	Customer *customer
	LoginLog *loginLog
	Order    *order
	User     *user
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
	Customer = &Q.Customer
	LoginLog = &Q.LoginLog
	Order = &Q.Order
	User = &Q.User
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:       db,
		Customer: newCustomer(db, opts...),
		LoginLog: newLoginLog(db, opts...),
		Order:    newOrder(db, opts...),
		User:     newUser(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	Customer customer
	LoginLog loginLog
	Order    order
	User     user
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:       db,
		Customer: q.Customer.clone(db),
		LoginLog: q.LoginLog.clone(db),
		Order:    q.Order.clone(db),
		User:     q.User.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:       db,
		Customer: q.Customer.replaceDB(db),
		LoginLog: q.LoginLog.replaceDB(db),
		Order:    q.Order.replaceDB(db),
		User:     q.User.replaceDB(db),
	}
}

type queryCtx struct {
	Customer ICustomerDo
	LoginLog ILoginLogDo
	Order    IOrderDo
	User     IUserDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		Customer: q.Customer.WithContext(ctx),
		LoginLog: q.LoginLog.WithContext(ctx),
		Order:    q.Order.WithContext(ctx),
		User:     q.User.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
