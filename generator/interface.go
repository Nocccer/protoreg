package generator

type Generator interface {
	Marshaler() (code string)
	Unmarshaler() (code string)
}
