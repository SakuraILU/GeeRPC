package codec

type ICodec interface {
	ReadHead(head *Head) error
	ReadBody(body interface{}) error

	Write(head *Head, body interface{}) error
}
