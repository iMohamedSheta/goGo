package enums

type TodoStatusEnum int8

const (
	IN_PROGRESS TodoStatusEnum = 1
	DONE        TodoStatusEnum = 2
	CANCELLED   TodoStatusEnum = 3
)
