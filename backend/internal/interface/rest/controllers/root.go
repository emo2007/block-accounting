package controllers

type RootController struct {
	Ping          PingController
	Auth          AuthController
	Organizations OrganizationsController
}

func NewRootController(
	ping PingController,
	auth AuthController,
	// organizations OrganizationsController,
) *RootController {
	return &RootController{
		Ping: ping,
		Auth: auth,
		// Organizations: organizations,
	}
}
