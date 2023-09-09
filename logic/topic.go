package logic

import (
	"edd2023-back/db"
	"edd2023-back/model"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type TopicLogic struct {
	Topic model.Topic
}

func (tl *TopicLogic) GenerateTopic() (bool, error) {
	log.Debugln("Start select format")

	query := fmt.Sprintf("select format from formats order by rand() limit 1")
	log.Debugln("--- select user query ---")
	log.Debugln(query)
	log.Debugln("-------------------------")

	rows, err := db.Db.Query(query)
	if err != nil {
		log.Errorln("Exec error: ", err)
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&tl.Topic.Format)
		if err != nil {
			log.Errorln("Exec error: ", err)
			return false, err
		}
	}

	log.Debugln("Start select genre")

	query = fmt.Sprintf("select genre, word from genres order by rand() limit 2")
	log.Debugln("--- select user query ---")
	log.Debugln(query)
	log.Debugln("-------------------------")

	rows, err = db.Db.Query(query)
	if err != nil {
		log.Errorln("Exec error: ", err)
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var genre model.Genre
		var word string
		err = rows.Scan(&genre, &word)
		if err != nil {
			log.Errorln("Exec error: ", err)
			return false, err
		}
		tl.Topic.Blanks = append(tl.Topic.Blanks, genre)
		tl.Topic.Words = append(tl.Topic.Words, word)
	}

	return true, nil
}
