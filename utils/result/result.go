package result

type Result struct {
	Code    int         `json:"code"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Msg     string      `json:"msg"`
}

func (result *Result) Ok(code int, data interface{}, msg string) *Result {
	result.Code = code
	result.Success = true
	result.Data = data
	result.Msg = msg
	return result
}

func (result *Result) Error(code int, data interface{}, msg string) *Result {
	result.Code = code
	result.Success = false
	result.Data = data
	result.Msg = msg
	return result
}

