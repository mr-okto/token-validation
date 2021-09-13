package id_generator

type IdGenerator interface {
	GenerateId() int32
}

type idGenerator struct {
	curId int32
}

var (
	instance *idGenerator
)

func GetInstance() IdGenerator {
	if instance == nil {
		instance = &idGenerator{}
	}
	return instance
}

func (s *idGenerator) GenerateId() int32 {
	s.curId++
	return s.curId
}
