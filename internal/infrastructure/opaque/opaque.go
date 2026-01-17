package opaque

type OpaqueService interface {
	Init() error
	RegisterInit(message []byte) (response []byte, credID []byte, err error)
	RegisterFinalize(message, credID []byte, username string) ([]byte, error)
	LoginInit(message, userRecord []byte) ([]byte, error)
	LoginFinalize(message []byte) ([]byte, error)
}
