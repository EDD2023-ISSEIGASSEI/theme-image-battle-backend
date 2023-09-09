package logic

import "edd2023-back/model"

type GeneratePhaseStateLigic struct {
	State *model.GeneratePhaseState
}

func (gl *GeneratePhaseStateLigic) FromGameSession(session model.GameSession) {
	gl.State = &model.GeneratePhaseState{
		PlayerStates: session.PlayerStates,
		Time:         session.Time,
	}
}
