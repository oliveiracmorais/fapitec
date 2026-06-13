package objetos_de_valor

import "strings"

type Endereco struct {
	Logradouro  string
	Numero      string
	Complemento string
	Bairro      string
	Cidade      string
	UF          string
	CEP         string
}

func NovoEndereco(logradouro, numero, complemento, bairro, cidade, uf, cep string) Endereco {
	return Endereco{
		Logradouro:  strings.TrimSpace(logradouro),
		Numero:      strings.TrimSpace(numero),
		Complemento: strings.TrimSpace(complemento),
		Bairro:      strings.TrimSpace(bairro),
		Cidade:      strings.TrimSpace(cidade),
		UF:          strings.TrimSpace(uf),
		CEP:         strings.TrimSpace(cep),
	}
}

func (e Endereco) Completo() string {
	partes := make([]string, 0, 4)

	logradouro := e.Logradouro
	if e.Numero != "" {
		logradouro += ", " + e.Numero
	}
	if e.Complemento != "" {
		logradouro += " - " + e.Complemento
	}
	if logradouro != "" {
		partes = append(partes, logradouro)
	}

	if e.Bairro != "" {
		partes = append(partes, e.Bairro)
	}

	cidadeUF := e.Cidade
	if e.UF != "" {
		if cidadeUF != "" {
			cidadeUF += "/" + e.UF
		} else {
			cidadeUF = e.UF
		}
	}
	if cidadeUF != "" {
		partes = append(partes, cidadeUF)
	}

	if e.CEP != "" {
		if len(e.CEP) == 8 {
			partes = append(partes, e.CEP[:5]+"-"+e.CEP[5:])
		} else {
			partes = append(partes, e.CEP)
		}
	}

	return strings.Join(partes, ", ")
}

func (e Endereco) Vazio() bool {
	return e.Logradouro == "" && e.Numero == "" &&
		e.Complemento == "" && e.Bairro == "" &&
		e.Cidade == "" && e.UF == "" && e.CEP == ""
}
