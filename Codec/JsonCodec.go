package codec

import (
	"encoding/json"
	"io"
)

type JsonCodec struct {
	stream io.ReadWriter
	dec    *json.Decoder
	enc    *json.Encoder
}

func NewJsonCodec(stream io.ReadWriter) *JsonCodec {
	return &JsonCodec{
		stream: stream,
		dec:    json.NewDecoder(stream),
		enc:    json.NewEncoder(stream),
	}
}

func (j *JsonCodec) ReadHead(head *Head) (err error) {
	err = j.dec.Decode(head)
	return
}

func (j *JsonCodec) ReadBody(body interface{}) (err error) {
	err = j.dec.Decode(body)
	return
}

func (j *JsonCodec) Write(head *Head, body interface{}) (err error) {
	if err = j.enc.Encode(head); err != nil {
		return
	}

	if err = j.enc.Encode(body); err != nil {
		return
	}

	return
}
