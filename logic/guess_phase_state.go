package logic

import "edd2023-back/model"

type GuessPhaseStateLigic struct {
	State *model.GuessPhaseState
}

func (gl *GuessPhaseStateLigic) FromGameSession(session model.GameSession) {
	var q model.GeneratedQuestionForGuess
	for _, gq := range session.GeneratedQuestions {
		if gq.Player.Id == session.DealerPlayerId {
			q = model.GeneratedQuestionForGuess{
				TopicForGuess: model.TopicForGuess{
					Format: gq.Topic.Format,
					Blanks: gq.Topic.Blanks,
				},
				ResultImageUrl: gq.ResultImageUrl,
			}
			break
		}
	}

	gl.State = &model.GuessPhaseState{
		PlayerStates:   session.PlayerStates,
		DealerPlayerId: session.DealerPlayerId,
		Question:       q,
		Time:           session.Time,
	}
}
