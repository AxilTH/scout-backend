// internal/repository/transaction.go
package repository

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

// TransactionManager управляет транзакциями базы данных
type TransactionManager struct {
	db *sqlx.DB
}

// NewTransactionManager создает новый менеджер транзакций
func NewTransactionManager(db *sqlx.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// TransactionFunc представляет функцию, которая выполняется в транзакции
type TransactionFunc func(tx *sqlx.Tx) error

// WithTransaction выполняет функцию в транзакции
// Если функция возвращает ошибку, транзакция откатывается
// Если функция завершается успешно, транзакция фиксируется
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn TransactionFunc) error {
	tx, err := tm.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	// Убеждаемся, что транзакция будет откачена при панике
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r) // повторно генерируем панику
		}
	}()

	// Выполняем функцию в транзакции
	if err := fn(tx); err != nil {
		// Откатываем транзакцию при ошибке
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.Join(err, rbErr)
		}
		return err
	}

	// Фиксируем транзакцию при успехе
	return tx.Commit()
}

// WithTx выполняет функцию в транзакции (без контекста)
// Упрощенная версия для обратной совместимости
func (tm *TransactionManager) WithTx(fn TransactionFunc) error {
	return tm.WithTransaction(context.Background(), fn)
}

// TxRepository представляет репозиторий, который может работать с транзакциями
type TxRepository interface {
	// WithTx возвращает новый экземпляр репозитория, работающий с транзакцией
	WithTx(tx *sqlx.Tx) interface{}
}

// TransactionalRepository представляет базовый репозиторий с поддержкой транзакций
type TransactionalRepository struct {
	db *sqlx.DB
	tm *TransactionManager
}

// NewTransactionalRepository создает новый репозиторий с поддержкой транзакций
func NewTransactionalRepository(db *sqlx.DB) *TransactionalRepository {
	return &TransactionalRepository{
		db: db,
		tm: NewTransactionManager(db),
	}
}

// DB возвращает базовое подключение к БД
func (r *TransactionalRepository) DB() *sqlx.DB {
	return r.db
}

// TransactionManager возвращает менеджер транзакций
func (r *TransactionalRepository) TransactionManager() *TransactionManager {
	return r.tm
}

// WithTransaction выполняет функцию в транзакции
func (r *TransactionalRepository) WithTransaction(ctx context.Context, fn TransactionFunc) error {
	return r.tm.WithTransaction(ctx, fn)
}

// WithTx выполняет функцию в транзакции (без контекста)
func (r *TransactionalRepository) WithTx(fn TransactionFunc) error {
	return r.tm.WithTx(fn)
}