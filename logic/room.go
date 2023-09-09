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
		PlayerNum:     1,
		MaxPlayerNum:  req.MaxMember,
		OwnerPlayerId: userId,
	}
	return nil
}

func (rl *RoomLogic) VaridatePassword(password string) bool {
	return rl.Room.Password == password
}

func (rl *RoomLogic) CanJoin() bool {
	return rl.Room.PlayerNum+1 <= rl.Room.MaxPlayerNum
}

func (rl *RoomLogic) Join() {
	rl.Room.PlayerNum += 1
}
