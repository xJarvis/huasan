package exhttp

import (
	"github.com/xJarvis/huashan/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * API 对外接口
 */
type ApiResult struct {
	HRet 		int											// 结果 0 失败  | 1 成功| 2 异常
	Result		string										// 结果 success | fail  | error
	Code   		string										// 状态码 fail(500-599) error(400-499) success(200-299)
	Msg    		string										// 消息
	Data   		interface{}								// 数据
}

type PageData struct {
	Total 		int 	`json:"total"`						//总条目数
	Page 		int 	`json:"page"`						//总页面数
	PageSize 	int 	`json:"page_size"`					//当前页面容量
	PageNum 	int 	`json:"page_num"`					//当前页数
	List 		[]interface{}	`json:"list"`				//列表
}

type AdminPage struct {
	PageNum    int  //当前页
	PageSize   int  //每页显示条数
	TotalPage  int  //总页数
	Total      int  //总条数
	FirstPage  bool //是否是首页
	LastPage   bool //是否是最后一页
	List       interface{}
	Sources    map[string]interface{}
	Field      interface{}
}


func NewApiResult(hret int, code string, data interface{}, msg string) *ApiResult {
	result := new(ApiResult)
	result.HRet = hret
	result.Code = code
	result.Data = data
	result.Msg = msg

	switch result.HRet {
	case HTTP_HRET_SUCCESS:
		result.Result = HTTP_RESULT_SUCCESS
		if result.Msg == "" {
			result.Msg = "request success!"
		}
	case HTTP_HRET_FAIL:
		result.Result = HTTP_RESULT_FAIL
		if result.Msg == "" {
			result.Msg = "request failed!"
		}
	case HTTP_HRET_ERROR:
		result.Result = HTTP_RESULT_ERROR
		if result.Msg == "" {
			result.Msg = "request error!"
		}
	default :
		result.HRet = HTTP_HRET_ERROR
		result.Result = HTTP_RESULT_ERROR
		result.Data = nil
		if result.Msg == "" {
			result.Msg = "not support!"
		}
	}
	return result
}

/**
 *  total 总条目数
 *  page_size 页面支持条目数
 *  page_num  页码
 */
func NewPageData(total int,page_size int,page_num int,list []interface{}) *PageData {
	result := new(PageData)
	result.Total = total
	page := total / page_size
	if total % page_size > 0 {
		page +=  1
	}
	result.Page = page
	result.List = list

	result.PageNum = page_num
	result.PageSize = len(list) //页面实际条目数
	return result
}

func NewAdminPage(total int, pageNum int, pageSize int, list interface{}) AdminPage {
	sources := make(map[string]interface{})
	tp := total / pageSize
	if total % pageSize > 0 {
		tp = total / pageSize + 1
	}
	return AdminPage {
		PageNum: pageNum,
		PageSize: pageSize,
		TotalPage: tp,
		Total: total,
		FirstPage: pageNum == 1,
		LastPage: pageNum == tp,
		List: list,
		Sources: sources,
	}
}

// API返回数据模型
func ApiJsonResult(c *gin.Context,data *ApiResult) {
	c.JSON(http.StatusOK, gin.H{
		"hRet"		: data.HRet,
		"result"	: data.Result,
		"code"		: data.Code,
		"data"		: data.Data,
		"msg"		: data.Msg,
	})
	return
}

func ApiJsonResultParamError(c *gin.Context, params []string,msg string) {
	errMsg := fmt.Sprintf("param error! %v %v!",params,msg)
	logger.Error("request error:" + errMsg)
	result := NewApiResult(HTTP_HRET_ERROR,HTTP_CODE_ERROR_PARAM,nil,errMsg)
	ApiJsonResult(c,result)
	return
}

func ApiJsonResultError(c *gin.Context, data interface{},msg string) {
	logger.Error("request error:" + msg)
	result := NewApiResult(HTTP_HRET_ERROR,HTTP_CODE_ERROR,nil,msg)
	ApiJsonResult(c,result)
	return
}

func ApiJsonResultFailed(c *gin.Context, data interface{},msg string) {
	logger.Error("request failed:" + msg)
	result := NewApiResult(HTTP_HRET_FAIL,HTTP_CODE_FAIL,data,msg)
	ApiJsonResult(c,result)
	return
}

func ApiJsonResultSuccess(c *gin.Context, data interface{},msg string) {
	logger.Debug("request success!")
	result := NewApiResult(HTTP_HRET_SUCCESS,HTTP_CODE_SUCCESS,data,msg)
	ApiJsonResult(c,result)
	return
}





// API返回数据模型-PHP返回值模式
func ApiJsonResultPhp(c *gin.Context,data *ApiResult) {
	c.JSON(http.StatusOK, gin.H{
		"isSuccess"		: data.HRet,
		"type"			: data.Result,
		"code"			: data.Code,
		"data"			: data.Data,
		"Msg"			: data.Msg,
	})
	return
}

func ApiJsonResultParamPhpError(c *gin.Context, params []string,msg string) {
	errMsg := fmt.Sprintf("param error! %v %v!",params,msg)
	logger.Error("request error:" + errMsg)
	result := NewApiResult(HTTP_HRET_ERROR,HTTP_CODE_ERROR_PARAM,nil,errMsg)
	ApiJsonResultPhp(c,result)
	return
}

func ApiJsonResultPhpError(c *gin.Context, data interface{},msg string) {
	logger.Error("request error:" + msg)
	result := NewApiResult(HTTP_HRET_ERROR,HTTP_CODE_ERROR,nil,msg)
	ApiJsonResultPhp(c,result)
	return
}

func ApiJsonResultPhpFailed(c *gin.Context, data interface{},msg string) {
	logger.Error("request failed:" + msg)
	result := NewApiResult(HTTP_HRET_FAIL,HTTP_CODE_FAIL,data,msg)
	ApiJsonResultPhp(c,result)
	return
}

func ApiJsonResultPhpSuccess(c *gin.Context, data interface{},msg string) {
	logger.Debug("request success!")
	result := NewApiResult(HTTP_HRET_SUCCESS,HTTP_CODE_SUCCESS,data,msg)
	ApiJsonResultPhp(c,result)
	return
}