package controllers

type FilController struct {
	service *services.FilService
}

func NewFilController(service *services.FilService) *FilController {
	return &FilController{service: service}
}
