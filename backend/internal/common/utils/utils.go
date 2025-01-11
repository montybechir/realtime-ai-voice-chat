package utils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func GenerateID(prefix string) string {
	timestamp := time.Now().Format("20250104150405")
	uid := uuid.New().String()[:8]
	return fmt.Sprintf("%s_%s_%s", prefix, timestamp, uid)
}
