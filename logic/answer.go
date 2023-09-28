package logic

import (
	"crypto/rand"
	"edd2023-back/model"
	"math/big"
)

type AnswerLogic struct {
	Answer model.Answer
}

func (al *AnswerLogic) CalcScore() error {
	// TODO Calculate Score
	n, err := rand.Int(rand.Reader, big.NewInt(500))
	if err != nil {
		return err
	}
	al.Answer.Score = int(n.Int64())
	return nil
}
