package routerbase

/**
 * @Author: lee
 * @Description:
 * @File: model
 * @Date: 2022/4/1 6:48 下午
 */

type AckBase struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
