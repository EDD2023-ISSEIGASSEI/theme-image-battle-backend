package logic

import (
	"edd2023-back/model"
	"edd2023-back/util"
)

type RoomLogic struct {
	Room model.Room
}

func (rl *RoomLogic) CreateRoom(req model.CreateRoomRequest, userId string) error {
	id, err := util.GenerateOtp()
	if err != nil {
		return err
	}

	rl.Room = model.Room{
		Id:            *id,
		Name:          req.Name,
		Password:      req.Password,
		PlayerNum:     0,
		MaxPlayerNum:  req.MaxMember,
		OwnerPlayerId: userId,
	}
	return nil
}
