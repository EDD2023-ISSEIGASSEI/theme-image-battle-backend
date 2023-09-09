package util

import (
	"bytes"
	"context"
	"crypto/rand"
	"edd2023-back/db"
	"math/big"
)

func GenerateOtp() (*string, error) {
	ctx := context.Background()
	var buffer bytes.Buffer
	for {
		for i := 0; i < 6; i++ {
			n, err := rand.Int(rand.Reader, big.NewInt(10))
			if err != nil {
				return nil, err
			}
			buffer.WriteString(n.String())
		}
		otp := buffer.String()
		res := db.Redis.Exists(ctx, otp).Val()
		if res == 0 {
			return &otp, nil
		}
	}
}
