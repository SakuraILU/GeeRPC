package codec

import (
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	baseCodec

	stream io.ReadWriter
	dec    *gob.Decoder
	enc    *gob.Encoder
}

func NewGobCodec(stream io.ReadWriter) *GobCodec {
	return &GobCodec{
		stream: stream,
		dec:    gob.NewDecoder(stream),
		enc:    gob.NewEncoder(stream),
	}
}

func (g *GobCodec) ReadHead(head *Head) (err error) {
	err = g.dec.Decode(head)
	if err != nil {
		log.Printf("Error: read head, %s", err)
	}
	return
}

func (g *GobCodec) ReadBody(body interface{}) (err error) {
	err = g.dec.Decode(body)
	if err != nil {
		log.Printf("Error: read body, %s", err)
	}
	return
}

func (g *GobCodec) Write(head *Head, body interface{}) (err error) {
	if err = g.enc.Encode(head); err != nil {
		log.Printf("Error: write head, %s", err)
		return
	}
	if err = g.enc.Encode(body); err != nil {
		log.Printf("Error: write body, %s", err)
		return
	}
	return
}
