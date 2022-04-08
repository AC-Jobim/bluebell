package controller

import (
	"bluebell/logic"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

/**
 * @Author: zhaobin
 * @Description 社区相关的
 * @Date: 2022-04-08 14:56
 */

func CommunityHandler(c *gin.Context) {
	// 查询到所有的社区  community_id community_name
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		// 不轻易把服务端报错报给用户
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

func CommunityDetailHandler(c *gin.Context) {
	communityID := c.Param("id")
	id, err := strconv.ParseInt(communityID, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
	}
	fmt.Println("id是", id)
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		//bu'qing'yi把服务端报错报给用户
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
