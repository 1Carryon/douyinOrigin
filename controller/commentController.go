package controller

import (
	"douyinOrigin/dao"
	"douyinOrigin/middleware/jwt"
	"douyinOrigin/service"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CommentActionResponse 评论操作返回响应
type CommentActionResponse struct {
	Response
	Comment service.CommentInfo `json:"comment"` // 评论成功返回评论内容，不需要重新拉取整个列
}

// CommentListResponse 评论列表返回响应
type CommentListResponse struct {
	Response
	CommentList []service.CommentInfo `json:"comment_list"` // 评论列表
}

// CommentAction
// 发表 or 删除评论 comment/action/
func CommentAction(c *gin.Context) {
	log.Println("CommentController-Comment_Action: running") //函数已运行

	//通过token获取userId
	tokenString := c.Query("token")
	myClaims, _ := jwt.ParseToken(tokenString)         //解析token
	userId, _ := strconv.ParseInt(myClaims.ID, 10, 64) //通过解析token，拿到userid

	fmt.Printf("userId:%v\n", userId)

	//获取videoId
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)

	//获取操作类型
	actionType, _ := strconv.ParseInt(c.Query("action_type"), 10, 32)

	//调用service层评论函数
	cService := new(service.CommentServiceImpl)
	if actionType == 1 { //actionType为1，则进行发表评论操作
		content := c.Query("comment_text")
		// 垃圾评论过滤。
		//666content = util.Filter.Replace(content, '#')

		// find, _ := util.Filter.FindIn(content)
		/*if find {
			log.Println("垃圾评论")
			c.JSON(http.StatusOK, CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "垃圾评论",
			})
			return
			content = "*****"
		}
		*/
		//发表评论数据准备
		var sendComment dao.Comment
		sendComment.UserId = userId
		sendComment.VideoId = videoId
		sendComment.CommentText = content
		timeNow := time.Now()
		sendComment.CreateDate = timeNow
		//发表评论
		commentInfo, err := cService.Send(sendComment)
		//发表评论失败
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{
					StatusCode: -1,
					StatusMsg:  "send comment failed",
				},
			})
			log.Println("CommentController-Comment_Action: return send comment failed") //发表失败
			return
		}

		//发表评论成功:
		//返回结果
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "send comment success",
			},
			Comment: commentInfo,
		})
		log.Println("CommentController-Comment_Action: return Send success") //发表评论成功，返回正确信息
		return
	} else { //actionType为2，则进行删除评论操作
		//获取要删除的评论的id
		commentId, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{
					StatusCode: -1,
					StatusMsg:  "delete commentId invalid",
				},
			})
			log.Println("CommentController-Comment_Action: return commentId invalid") //评论id格式错误
			return
		}
		//删除评论操作
		err = cService.DelComment(commentId)
		if err != nil { //删除评论失败
			str := err.Error()
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{
					StatusCode: -1,
					StatusMsg:  str,
				},
			})
			log.Println("CommentController-Comment_Action: return delete comment failed") //删除失败
			return
		}
		//删除评论成功
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "delete comment success",
			},
		})

		log.Println("CommentController-Comment_Action: return delete success") //函数执行成功，返回正确信息
		return
	}
}

// CommentList
// 查看评论列表 comment/list/
func CommentList(c *gin.Context) {
	log.Println("CommentController-Comment_List: running") //函数已运行
	//获取userId
	id, _ := c.Get("userId")
	userid, _ := id.(string)
	userId, err := strconv.ParseInt(userid, 10, 64)
	//log.Printf("err:%v", err)
	//log.Printf("userId:%v", userId)

	//获取videoId
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	//错误处理
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: -1,
			StatusMsg:  "comment videoId json invalid",
		})
		log.Println("CommentController-Comment_List: return videoId json invalid") //视频id格式有误
		return
	}
	log.Printf("videoId:%v", videoId)

	//调用service层评论函数
	commentService := new(service.CommentServiceImpl)
	commentList, err := commentService.GetList(videoId, userId)
	//commentList, err := commentService.GetListFromRedis(videoId, userId)
	if err != nil { //获取评论列表失败
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
		})
		log.Println("CommentController-Comment_List: return list false") //查询列表失败
		return
	}

	//获取评论列表成功
	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "get comment list success",
		},
		CommentList: commentList,
	})
	log.Println("CommentController-Comment_List: return success") //成功返回列表
	return
}
