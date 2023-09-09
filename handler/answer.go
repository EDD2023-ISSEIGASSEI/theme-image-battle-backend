package handler

import (
	"edd2023-back/logic"
	"edd2023-back/model"
	"edd2023-back/util"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type AnswerHandler struct{}

type SubmitAnswerRequest struct {
	Answers []string `json:"answers"`
}

func (*AnswerHandler) SubmitAnswer(ctx *gin.Context) {
	var req SubmitAnswerRequest
	err := ctx.Bind(&req)
	if err != nil {
		s := err.Error()
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	sessionId, _ := ctx.Cookie("sessionId")
	asl := logic.AuthSessionLogic{
		Session: model.AuthSession{Uuid: sessionId},
	}
	asl.GetByUuid()
	player := model.UserToPlayer(asl.Session.User)

	gameSessionId, _ := ctx.Cookie("gameSessionId")
	gsl := logic.GameSessionLogic{}
	gsl.GetByUuid(gameSessionId)

	ans := model.Answer{
		Player:       player,
		BlankAnswers: req.Answers,
		Score:        0,
	}
	al := logic.AnswerLogic{
		Answer: ans,
	}
	al.CalcScore()

	var q model.GeneratedQuestion
	for _, gq := range gsl.Session.GeneratedQuestions {
		if gq.Player.Id == gsl.Session.DealerPlayerId {
			q = gq
			break
		}
	}
	aq := model.AnswerForQuestion{
		DealerPlayerId:   gsl.Session.DealerPlayerId,
		QuestionImageUrl: q.ResultImageUrl,
		Answer:           ans,
	}
	gsl.Session.PlayerAnswers = append(gsl.Session.PlayerAnswers, aq)
	f := true
	dealerNum := 0
	gsl.Session.GeneratedQuestions = append(gsl.Session.GeneratedQuestions, q)
	for idx, ps := range gsl.Session.PlayerStates {
		if ps.Player.Id == player.Id {
			gsl.Session.PlayerStates[idx].IsCompleted = true
		}
		if !ps.IsCompleted {
			f = false
		}
		if ps.Player.Id == gsl.Session.DealerPlayerId {
			dealerNum = idx
		}
	}

	// 全員回答が終わったら
	if f {
		for idx := range gsl.Session.PlayerStates {
			gsl.Session.PlayerStates[idx].IsCompleted = false
		}
		// 次の親がいるなら
		if dealerNum+1 < len(gsl.Session.Players) {
			gsl.Session.DealerPlayerId = gsl.Session.Players[dealerNum+1].Id
			gsl.Session.PlayerAnswers = []model.AnswerForQuestion{}
		} else {
			// 次の親がいないのなら
			gsl.Session.Phase = model.ShowScorePhase
			gsl.Session.DealerPlayerId = gsl.Session.Players[0].Id
			gsl.Session.ShowingPlayerId = gsl.Session.Players[1].Id
		}
	}

	gsl.UpdateByUuId()
	// ToDo: broadcast

	ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
}
