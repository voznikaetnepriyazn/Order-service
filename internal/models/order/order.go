package order

import "github.com/google/uuid"

type Order struct {
	Id           uuid.UUID
	IdOfCustomer uuid.UUID
}
