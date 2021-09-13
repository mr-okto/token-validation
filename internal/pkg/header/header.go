package header

type Header struct {
	SvcId      int32
	BodyLength int32
	RequestId  int32
}

const (
	Size = 12 // sizeOf(SvcId) + sizeOf(BodyLength) + sizeOf(RequestId)
)
