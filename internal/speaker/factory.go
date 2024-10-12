package speaker

type Factory interface {
	Create(conn any) (speaker Speaker, err error)
}
