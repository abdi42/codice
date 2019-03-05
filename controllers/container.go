package controllers

type CreateContainReq struct {
	Source string `json:"source"`
	Lang   string `json:"lang"`
}
