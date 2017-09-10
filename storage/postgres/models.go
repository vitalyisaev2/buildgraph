package postgres

type model struct {
	id int
}

func (m *model) GetID() int   { return m.id }
func (m *model) SetID(id int) { m.id = id }

type author struct {
	model
	name  string
	email string
}

func (m *author) GetID() int       { return m.id }
func (m *author) GetName() string  { return m.name }
func (m *author) GetEmail() string { return m.email }

type project struct {
	model
	namespace string
	name      string
}

func (m *project) GetID() int           { return m.id }
func (m *project) GetNamespace() string { return m.namespace }
func (m *project) GetName() string      { return m.name }
