package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// generateOrderNum returns a server-generated, hard-to-guess order number.
// Format: AP<yyyymmddHHMMSS><6 random hex chars>, e.g. AP20260601120000a1b2c3.
// The timestamp keeps it sortable and the random suffix makes collisions and
// client forgery impractical. It fits in the OrderNum varchar(32) column.
func generateOrderNum() string {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		// Fall back to a nanosecond suffix if the RNG is unavailable.
		return fmt.Sprintf("AP%s%06d", time.Now().Format("20060102150405"), time.Now().Nanosecond()%1000000)
	}
	return fmt.Sprintf("AP%s%s", time.Now().Format("20060102150405"), hex.EncodeToString(b))
}
