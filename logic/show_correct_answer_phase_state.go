package logic

import "edd2023-back/model"

type ShowCorrectAnswerPhaseStateLigic struct {
	State *model.ShowCorrectAnswerPhaseState
}

func (sl *ShowCorrectAnswerPhaseStateLigic) FromGameSession(session model.GameSession) {
	var q model.GeneratedQuestion
	for _, gq := range session.GeneratedQuestions {
		if gq.Player.Id == session.DealerPlayerId {
			q = gq
			break
		}
	}

	answers := []model.Answer{}
	for _, pa := range session.PlayerAnswers {
		if pa.DealerPlayerId == session.DealerPlayerId {
			answers = append(answers, pa.Answer)
		}
	}

	var score int
	for _, ps := range session.PlayerStates {
		if ps.Player.Id == session.DealerPlayerId {
			score = ps.Score
			break
		}
	}

	sl.State = &model.ShowCorrectAnswerPhaseState{
		Question:    q,
		Answers:     answers,
		DealerScore: score,
	}
}
