package gate

type Gate struct{}

func (gate *Gate) Allowed() bool {
	return true
}
