package model

import (
	"encoding/base64"
	"fmt"
	"io"
)

type Bytes []byte

func (b *Bytes) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("Bytes must be a string")
	}

	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}

	*b = decoded
	return nil
}

// MarshalGQL is a custom function to handle encoding
func (b Bytes) MarshalGQL(w io.Writer) {
	_, _ = fmt.Fprint(w, "\""+base64.StdEncoding.EncodeToString(b)+"\"")
}
