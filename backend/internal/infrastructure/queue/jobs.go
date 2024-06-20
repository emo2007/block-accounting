package queue

type JobDeployMultisig struct {
	OwnersPubKeys []string `json:"pub_keys"`
	Confirmations int      `json:"confirmations"`
}
