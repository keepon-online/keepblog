package copier

import (
	"encoding/json"
	"fmt"
)

func Copy(to, from any) error {
	b, err := json.Marshal(from)
	if err != nil {
		return fmt.Errorf("marshal from data err: %w", err)
	}
	if err = json.Unmarshal(b, to); err != nil {
		return fmt.Errorf("unmarshal to data err: %w", err)
	}
	return nil
}
