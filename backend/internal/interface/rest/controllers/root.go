package controllers

type RootController struct {
	Ping          PingController
	Auth          AuthController
	Organizations OrganizationsController
	Transactions  TransactionsController
	Participants  ParticipantsController
}

func NewRootController(
	ping PingController,
	auth AuthController,
	organizations OrganizationsController,
	transactions TransactionsController,
	participants ParticipantsController,
) *RootController {
	return &RootController{
		Ping:          ping,
		Auth:          auth,
		Organizations: organizations,
		Transactions:  transactions,
		Participants:  participants,
	}
}
