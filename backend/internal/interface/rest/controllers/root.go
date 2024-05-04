package controllers

type RootController struct {
	Ping PingController
	Auth AuthController
}

func NewRootController(
	ping PingController,
	auth AuthController,
) *RootController {
	return &RootController{
		Ping: ping,
		Auth: auth,
	}
}
