package handler

import (
	"edd2023-back/logic"
	"edd2023-back/model"
	"edd2023-back/util"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type PromptHandler struct{}

type SubmitPromptRequest struct {
	Prompt string `json:"prompt"`
}

func (*PromptHandler) SubmitPrompt(ctx *gin.Context) {
	var req SubmitPromptRequest
	err := ctx.Bind(&req)
	if err != nil {
		s := err.Error()
		r := util.BadRequest(&s)
		log.Errorln("[Error]request parse error: ", s)
		ctx.JSON(r.StatusCode, r.Message)
		return
	}

	pl := logic.PromptLogic{
		Prompt: req.Prompt,
	}
	imageUrl, err := pl.GenerateImage()
	if err != nil {
		log.Errorln("[Error]exec error: ", err.Error())
		r := util.InternalServerError(nil)
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

	var topic model.Topic
	for _, t := range gsl.Session.PlayerTopics {
		if t.Player.Id == player.Id {
			topic = t.Question
			break
		}
	}

	gq := model.GeneratedQuestion{
		Player:         player,
		Topic:          topic,
		Prompt:         req.Prompt,
		ResultImageUrl: imageUrl,
	}

	f := true
	gsl.Session.GeneratedQuestions = append(gsl.Session.GeneratedQuestions, gq)
	for idx, ps := range gsl.Session.PlayerStates {
		if ps.Player.Id == player.Id {
			gsl.Session.PlayerStates[idx].IsCompleted = true
		}
		if !ps.IsCompleted {
			f = false
		}
	}

	if f {
		gsl.Session.Phase = model.GuessPhase
		gsl.Session.DealerPlayerId = gsl.Session.Players[0].Id
		for idx := range gsl.Session.PlayerStates {
			gsl.Session.PlayerStates[idx].IsCompleted = false
		}
	}

	gsl.UpdateByUuId()
	// ToDo: broadcast

	ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
}
