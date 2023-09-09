package logic

import "edd2023-back/model"

type ShowScorePhaseStateLigic struct {
	State *model.ShowScorePhaseState
}

func (sl *ShowScorePhaseStateLigic) FromGameSession(session model.GameSession) {
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

	var ans model.AnswerForQuestion
	for _, pa := range session.PlayerAnswers {
		if pa.DealerPlayerId == session.DealerPlayerId &&
			pa.Answer.Player.Id == session.ShowingPlayerId {
			ans = pa
			break
		}
	}

	sl.State = &model.ShowScorePhaseState{
		PlayerStates:    session.PlayerStates,
		DealerPlayerId:  session.DealerPlayerId,
		ShowingPlayerId: session.ShowingPlayerId,
		Question:        q,
		PlayerAnswer:    ans,
	}
}
