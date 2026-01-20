package utils

import "encoding/json"

type NullString struct {
	Value   *string // Ponteiro para string (nil = null)
	Present bool    // Foi enviado no JSON?
}

// UnmarshalJSON é chamado AUTOMATICAMENTE pelo json.Unmarshal
// quando ele encontra um campo do tipo NullString
func (n *NullString) UnmarshalJSON(data []byte) error {
	// Marca que o campo FOI enviado no JSON
	n.Present = true

	// Caso 1: JSON contém "null" (literal)
	if string(data) == "null" {
		n.Value = nil // Value fica nil, mas Present = true
		return nil
	}

	// Caso 2: JSON contém uma string
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err // Erro de parsing (ex: número ao invés de string)
	}

	n.Value = &value // Guarda o ponteiro para o valor
	return nil
}

func (n NullString) MarshalJSON() ([]byte, error) {
	if !n.Present || n.Value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(*n.Value)
}
