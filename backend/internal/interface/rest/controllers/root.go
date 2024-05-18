package controllers

type RootController struct {
	Ping          PingController
	Auth          AuthController
	Organizations OrganizationsController
	Transactions  TransactionsController
}

func NewRootController(
	ping PingController,
	auth AuthController,
	organizations OrganizationsController,
	transactions TransactionsController,
) *RootController {
	return &RootController{
		Ping:          ping,
		Auth:          auth,
		Organizations: organizations,
		Transactions:  transactions,
	}
}
