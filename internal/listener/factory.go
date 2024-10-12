package listener

type Factory interface {
	Create(conn any) (listener Listener, err error)
}
