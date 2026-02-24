package goodinorder

import "github.com/google/uuid"

type GoodInOrder struct {
	IdOfGood  uuid.UUID
	IdOfOrder uuid.UUID
}
