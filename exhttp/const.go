package exhttp

/**
 *  HTTP-基础状态数据
 */
const (
	HTTP_RESULT_SUCCESS 	string = "success"			// 请求成功
	HTTP_RESULT_FAIL  		string = "fail"        		// 请求失败
	HTTP_RESULT_ERROR  		string = "error"			// 请求错误

	HTTP_HRET_FAIL  		int = 0    					// 失败
	HTTP_HRET_SUCCESS 		int = 1						// 成功
	HTTP_HRET_ERROR   		int = 2						// 异常

	HTTP_CODE_SUCCESS		string    = "200"			// 默认成功
	HTTP_CODE_ERROR			string    = "400"			// 默认错误
	HTTP_CODE_FAIL			string    = "500"			// 默认失败
	HTTP_CODE_NOT_SUPPORT	string    = "700"			// 不支持
	HTTP_CODE_ERROR_PARAM	string    = "401"			// 参数错误

	HTTP_PAGE_SIZE_DEFAULT int =  100					// 默认页面条目数
)


