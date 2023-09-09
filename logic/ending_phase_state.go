package logic

import (
	"edd2023-back/model"
	"sort"
)

type EndingPhaseStateLigic struct {
	State *model.EndingPhaseState
}

func (sl *EndingPhaseStateLigic) FromGameSession(session model.GameSession) {
	ps := session.PlayerStates
	sort.Sort(model.PlayerStateByScore(ps))

	sl.State = &model.EndingPhaseState{
		Ranking: ps,
	}
}
