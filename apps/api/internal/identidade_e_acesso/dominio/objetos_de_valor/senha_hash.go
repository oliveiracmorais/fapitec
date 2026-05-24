package objetos_de_valor

type SenhaHash struct {
	hash string
}

func NovaSenhaHash(hash string) SenhaHash {
	return SenhaHash{hash: hash}
}

func (s SenhaHash) String() string {
	return s.hash
}
