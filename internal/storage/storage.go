package storage

import (
	"Order/internal/models/order"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUrlNotFound = errors.New("url not found")
	ErrUrlExist    = errors.New("url exist")
)

type OrderService interface {
	AddURL(order order.Order) (uuid.UUID, error)

	DeleteURL(id uuid.UUID) error

	GetAllURL() ([]order.Order, error)

	GetByIdURL(id uuid.UUID) (uuid.UUID, error)

	UpdateURL(order order.Order) error

	IsOrderCreatedURL(id uuid.UUID) (bool, error)
}
