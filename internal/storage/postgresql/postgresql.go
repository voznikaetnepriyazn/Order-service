package postgresql

import (
	"Order/internal/models/order"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Storage struct {
	db *sql.DB //connection string
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("pgx", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: failed to ping db: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddURL(order order.Order) (uuid.UUID, error) {
	const op = "storage.postgresql.addURL"

	newID := uuid.New()

	stmt, err := s.db.Prepare(
		`INSERT INTO Order ("Id", "idOfCustomer") 
		VALUES ($1, $2)
		`)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var insertedID uuid.UUID
	err = stmt.QueryRow(newID, order.IdOfCustomer).Scan(&insertedID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: execute query: %w", op, err)
	}

	return insertedID, nil
}

func (s *Storage) DeleteURL(id uuid.UUID) error {
	const op = "storage.postgresql.deleteURL"

	stmt, err := s.db.Prepare(
		`DELETE FROM Order WHERE id=Id`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetAllURL() ([]order.Order, error) {
	const op = "storage.postgresql.getAllURL"

	stmt, err := s.db.Prepare(`
		SELECT Order.Id
		FROM Order 
		INNER JOIN dbo.GoodInOrder ON Order.IdOfClient = GoodInOrder.IdOfClient 
		INNER JOIN Good ON GoodInOrder.Id = Good.Id
		`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	row, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var orders []order.Order
	for row.Next() {
		var order order.Order
		err := row.Scan(&order)
		if err != nil {
			return nil, fmt.Errorf("%s: scann failed: %w", op, err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (s *Storage) GetByIdURL(id uuid.UUID) (uuid.UUID, error) {
	const op = "storage.postgresql.getByIdURL"

	stmt, err := s.db.Prepare(`
	SELECT * FROM dbo.Order WHERE id=Id'
	`)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var order uuid.UUID
	err = stmt.QueryRow(id).Scan(&order)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.UUID{}, fmt.Errorf("%s: order not found", op)
		}
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	return order, nil
}

func (s *Storage) UpdateURL(order order.Order) error {
	const op = "storage.postgresql.updateURL"

	newID := uuid.New()

	stmt, err := s.db.Prepare(
		`INSERT INTO Order ("Id", "idOfCustomer") 
		VALUES ($1, $2)
		`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var insertedID uuid.UUID
	err = stmt.QueryRow(newID, order.IdOfCustomer).Scan(&insertedID)
	if err != nil {
		return fmt.Errorf("%s: execute query: %w", op, err)
	}

	return nil
}

func (s *Storage) IsOrderCreatedURL(id uuid.UUID) (bool, error) {
	ord, err := s.GetByIdURL(id)
	if err != nil {
		return false, err
	}

	if ord == (uuid.UUID{}) {
		return false, err
	}

	return true, nil
}
