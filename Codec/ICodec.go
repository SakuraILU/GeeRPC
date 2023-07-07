package codec

import "io"

type ICodec interface {
	ReadHead(head *Head) error
	ReadBody(body interface{}) error

	Write(head *Head, body interface{}) error
}

var CodecNewFuncs = map[CodecType]func(io.ReadWriter) ICodec{
	GOBTYPE:  NewGobCodec,
	JSONTYPE: NewJsonCodec,
}
