package codec

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
)

const MAGICNUM = 0x114514

type CodecType int

const (
	GOBTYPE CodecType = iota
	JSONTYPE
)

func (c CodecType) String() (name string) {
	switch c {
	case GOBTYPE:
		name = "GobType"
	case JSONTYPE:
		name = "JsonType"
	default:
		name = "UnknownType"
	}
	return
}

type Option struct {
	Magic_num      int
	Codec_type     CodecType
	Conn_timeout   time.Duration // 0 means no timeout
	Handle_timeout time.Duration
}

func NewOption(tp CodecType, conn_timeout, handle_timeout time.Duration) (*Option, error) {
	if (tp != GOBTYPE) && (tp != JSONTYPE) {
		err_msg := fmt.Sprintf("ERROR: unsupported type %s, only support %s and %s", tp, GOBTYPE, JSONTYPE)
		log.Println(err_msg)
		return nil, fmt.Errorf(err_msg)
	}

	return &Option{
		Magic_num:      MAGICNUM,
		Codec_type:     tp,
		Conn_timeout:   conn_timeout,
		Handle_timeout: handle_timeout,
	}, nil
}

func ParseOption(stream io.ReadWriter) (*Option, error) {
	var opt Option
	if err := json.NewDecoder(stream).Decode(&opt); err != nil || !opt.IsValid() {
		log.Println("ERROR: decode error")
		return nil, err
	}
	return &opt, nil
}

func (o *Option) IsValid() bool {
	if o.Magic_num != MAGICNUM {
		log.Printf("ERROR: Magic number wrong, unknown type")
		return false
	}

	if (o.Codec_type != GOBTYPE) && (o.Codec_type != JSONTYPE) {
		log.Printf("ERROR: unspported type %s, only support %s and %s", o.Codec_type, GOBTYPE, JSONTYPE)
		return false
	}
	return true
}
