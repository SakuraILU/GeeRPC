package codec

import (
	"encoding/json"
	"errors"
	"io"
)

type baseCodec struct {
	stream io.ReadWriter
}

func (bc *baseCodec) ReadOption() (*Option, error) {
	var opt Option
	if err := json.NewDecoder(bc.stream).Decode(&opt); err != nil {
		return nil, err
	}

	if !opt.IsValid() {
		return nil, errors.New("invalid option format")
	}
	return &opt, nil
}

func (bc *baseCodec) WriteOption(tp CodecType) (err error) {
	opt, err := NewOption(GOBTYPE)
	if err != nil {
		return err
	}

	if json.NewEncoder(bc.stream).Encode(&opt); err != nil {
		return err
	}
	return
}
