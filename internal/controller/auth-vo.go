package controller

type GenerateCaptchaReq struct{}

type GenerateCaptchaResp struct {
	Code int    `json:"code" binding:"required"`
	Msg  string `json:"msg" binding:"required"`
	Data struct {
		Image              string `json:"image" binding:"required"`
		VerificationCodeID string `json:"verification_code_id" binding:"required"`
	}
}

type LoginReq struct {
	UserID             string `json:"user_id" binding:"required"`
	Password           string `json:"password" binding:"required"`
	VerificationCodeID string `json:"verification_code_id" binding:"required"`
	VerificationCode   string `json:"verification_code" binding:"required"`
}

type LoginResp struct {
	Code int    `json:"code" binding:"required"`
	Msg  string `json:"msg" binding:"required"`
	Data struct {
		Token string `json:"token" binding:"required"`
	}
}
