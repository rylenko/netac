package speaker

import (
	"fmt"

	"github.com/google/uuid"
)

func generateRandomUUIDBytes() (bytes []byte, err error) {
	// Generate a new copy identifactor.
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate a new identificator: %v", err)
	}

	// Marshal identificator to bytes.
	idBytes, err := id.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal identificator to bytes: %v", err)
	}
	return idBytes, nil
}
