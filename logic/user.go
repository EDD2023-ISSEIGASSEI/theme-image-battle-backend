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

	query := fmt.Sprintf("insert into users (id, name, password, line_uid) values (?, ?, ?, ?)")
	log.Debugln("--- insert query ---")
	log.Debugln(query)
	log.Debugln("-------------------------")
	stmt, err := db.Db.Prepare(query)
	if err != nil {
		log.Errorln("Prepare error: ", err)
		return err
	}

	_, err = stmt.Exec(user.Id, user.Name, user.Password, user.LineUid)
	if err != nil {
		log.Errorln("Exec error: ", err)
		return err
	}

	return nil
}

func (ul *UserLigic) SelectByNameAndPass() error {
	log.Debugln("Start select user")

	query := fmt.Sprintf("select id, name, password, line_uid from users where name = ? and password = ?")
	log.Debugln("--- select user query ---")
	log.Debugln(query)
	log.Debugln("-------------------------")

	rows, err := db.Db.Query(query, ul.User.Name, ul.User.Password)
	if err != nil {
		log.Errorln("Exec error: ", err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var lineUid *string
		err = rows.Scan(&ul.User.Id, &ul.User.Name, &ul.User.Password, &lineUid)
		if err != nil {
			log.Errorln("Exec error: ", err)
			return err
		}
		ul.User.LineUid = lineUid
	}

	return nil
}
