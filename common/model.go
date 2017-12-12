package common

// ObjectID represnts unique identifier of an object
// (which is likely to be used as a primary key in RDBMS)
type ObjectID = int

// Model describes something that can be stored within a RDBMS
type Model interface {
	GetObjectID() ObjectID
	SetObjectID(ObjectID)
}

// Object is a default implementation of Model interface
type Object struct {
	id ObjectID
}

func (o *Object) GetObjectID() ObjectID {
	return o.id
}

func (o *Object) SetObjectID(val ObjectID) {
	o.id = val
}
