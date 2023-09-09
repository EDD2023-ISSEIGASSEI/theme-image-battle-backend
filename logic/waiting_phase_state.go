package logic

import "edd2023-back/model"

type WaitingPhaseStateLigic struct {
	State *model.WaitingPhaseState
}

func (wl *WaitingPhaseStateLigic) FromGameSession(session model.GameSession) {
	wl.State = &model.WaitingPhaseState{
		RoomInfo: session.Room,
		Players:  session.Players,
	}
}
