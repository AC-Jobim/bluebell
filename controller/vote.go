package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

/**
 * @Author: 投票
 * @Description TODO
 * @Date: 2022-04-11 12:44
 */

func PostVoteController(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		// 获取validator.ValidationErrors类型的errors
		// 可能没有触发校验规则就报错了
		errs, ok := err.(validator.ValidationErrors)
		zap.L().Error("vote.go", zap.Error(err))
		zap.L().Error("vote.go", zap.Error(errs))
		if !ok {
			// 如果不是参数校验错误
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
	}

	userId, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error(" getCurrentUserID(c)", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}

	if err := logic.VoteForPost(userId, p); err != nil {
		zap.L().Error(" logic.VoteForPost", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)

}
