package hasher

type ImageHasher interface {
	Hash([]byte) (string, error)
}
