package controllers

type RootController struct {
	Ping PingController
}

func NewRootController(
	ping PingController,
) *RootController {
	return &RootController{
		Ping: ping,
	}
}
