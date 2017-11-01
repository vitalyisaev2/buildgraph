package storage

type Model interface {
	GetID() int
	SetID(id int)
}

type Author interface {
	Model
	GetName() string
	GetEmail() string
}

type Project interface {
	Model
	GetNamespace() string
	GetName() string
}
