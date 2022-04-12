package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

/**
 * @Author: zhaobin
 * @Description TODO
 * @Date: 2022-04-08 16:33
 */

func CreatePostHandler(c *gin.Context) {
	// 1.获取参数以及参数的校验
	p := new(models.Post)
	if err := c.ShouldBind(p); err != nil {
		zap.L().Error("error", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//从c.content里面取到当前用户的ID值
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	fmt.Println("帖子的数据：", p)
	// 2.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 3.创建成功返回响应
	ResponseSuccess(c, nil)

}

// 获取帖子详情
func GetPostDetailHandler(c *gin.Context) {
	// 获取参数
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 根据id获取数据
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表的接口
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)
	// 获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandlerV2 根据时间或者分数查询所有帖子 V2
/*
	- 根据前端传来的参数动态的获取帖子列表
	- 按创建时间排序 或者 按照 分数排序
		1.获取请求的 query string 参数
		2.去redis查询id列表
		3.根据id去数据库查询帖子详细信息
*/
func GetPostListHandlerV2(c *gin.Context) {

	// 初始化结构体时指定初始参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime, // magic string
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	zap.L().Debug("GetPostListHandlerV2", zap.Any("查询帖子列表参数", p))
	// 进行逻辑判断，判断是否根据社区查询帖子列表，从而调用不同的方法
	data, err := logic.GetPostListNew(p)

	if err != nil {
		zap.L().Error("logic.GetPostList2() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)

}
