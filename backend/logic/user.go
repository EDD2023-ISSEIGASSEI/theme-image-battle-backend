package logic

import (
	"fmt"
	"line-bot-otp-back/db"
	"line-bot-otp-back/model"

	log "github.com/sirupsen/logrus"
)

type UserLigic struct {
	User *model.User
}

func (ul *UserLigic) Create() error {
	log.Debugln("Start crate user")
	user := ul.User

	query := fmt.Sprintf("insert into users (name, password) values (?, ?)")
	log.Debugln("--- bulk insert query ---")
	log.Debugln(query)
	log.Debugln("-------------------------")
	stmt, err := db.Db.Prepare(query)
	if err != nil {
		log.Errorln("Prepare error: ", err)
		return err
	}

	if _, err := stmt.Exec(user.Name, user.Password); err != nil {
		log.Errorln("Exec error: ", err)
		return err
	}

	return nil
}
