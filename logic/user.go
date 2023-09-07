package logic

import (
	"edd2023-back/db"
	"edd2023-back/model"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type UserLigic struct {
	User *model.User
}

func (ul *UserLigic) Create() error {
	log.Debugln("Start crate user")
	user := ul.User

	query := fmt.Sprintf("insert into users (id, name, password, line_uid, icon_image_url) values (?, ?, ?, ?, ?)")
	log.Debugln("--- insert query ---")
	log.Debugln(query)
	log.Debugln("-------------------------")
	stmt, err := db.Db.Prepare(query)
	if err != nil {
		log.Errorln("Prepare error: ", err)
		return err
	}

	_, err = stmt.Exec(user.Id, user.Name, user.Password, user.LineUid, user.IconImageUrl)
	if err != nil {
		log.Errorln("Exec error: ", err)
		return err
	}

	return nil
}

func (ul *UserLigic) SelectById() (bool, error) {
	log.Debugln("Start select user")

	query := fmt.Sprintf("select id, name, password, line_uid, icon_image_url from users where id = ?")
	log.Debugln("--- select user query ---")
	log.Debugln(query)
	log.Debugln("-------------------------")

	rows, err := db.Db.Query(query, ul.User.Id)
	if err != nil {
		log.Errorln("Exec error: ", err)
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&ul.User.Id, &ul.User.Name, &ul.User.Password, &ul.User.LineUid, &ul.User.IconImageUrl)
		if err != nil {
			log.Errorln("Exec error: ", err)
			return false, err
		}
	}

	return true, nil
}

func (ul *UserLigic) VaridatePassword(password string) bool {
	return ul.User.Password == password
}

func (ul *UserLigic) IdIsExists() (bool, error) {
	query := fmt.Sprintf("select id from users where id = ?")
	log.Debugln("--- select user query ---")
	log.Debugln(query)
	log.Debugln("-------------------------")

	rows, err := db.Db.Query(query, ul.User.Id)
	if err != nil {
		log.Errorln("Exec error: ", err)
		return false, err
	}
	defer rows.Close()
	count := 0
	var id string
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Errorln("Exec error: ", err)
			return false, err
		}
		count++
	}
	if count == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
