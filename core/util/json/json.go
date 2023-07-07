package json

import (
	"time"

	"github.com/star-table/go-common/utils/unsafe"
	"github.com/star-table/common/core/consts"
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

type TimeDecoder struct {
}
type TimeEncoder struct {
	precision time.Duration
}

var loc, _ = time.LoadLocation("Local")

func (td *TimeDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	str := iter.ReadString()

	mayBlank, _ := time.Parse(consts.AppTimeFormat, str)
	now, err := time.ParseInLocation(consts.AppTimeFormat, str, loc)

	if err != nil {
		*((*time.Time)(ptr)) = time.Unix(0, 0)
	} else if mayBlank.IsZero() {
		*((*time.Time)(ptr)) = mayBlank
	} else {
		*((*time.Time)(ptr)) = now
	}
}

func (codec *TimeEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	ts := *((*time.Time)(ptr))
	return ts.IsZero()
}

func (codec *TimeEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	ts := *((*time.Time)(ptr))
	if !ts.IsZero() {
		timestamp := ts.Unix()
		tm := time.Unix(timestamp, 0)
		format := tm.Format(consts.AppTimeFormat)
		stream.WriteString(format)
	} else {
		mayBlank, _ := time.Parse(consts.AppTimeFormat, consts.BlankString)
		stream.WriteString(mayBlank.Format(consts.AppTimeFormat))
	}

}

var json = jsoniter.ConfigFastest

func init() {
	extra.RegisterFuzzyDecoders()
	jsoniter.RegisterTypeEncoder("time.Time", &TimeEncoder{})
	jsoniter.RegisterTypeDecoder("time.Time", &TimeDecoder{})
}

func Marshal(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func Unmarshal(bs []byte, obj interface{}) error {
	return json.Unmarshal(bs, obj)
}

func ToJson(obj interface{}) (string, error) {
	bs, err := Marshal(obj)
	if err != nil {
		return "", err
	}
	return unsafe.BytesString(bs), nil
}

func ToJsonIgnoreError(obj interface{}) string {
	if obj == nil {
		return ""
	}
	bs, err := Marshal(obj)
	if err != nil {
		return ""
	}
	return unsafe.BytesString(bs)
}

func ToJsonBytesIgnoreError(obj interface{}) []byte {
	if obj == nil {
		return nil
	}
	bs, err := Marshal(obj)
	if err != nil {
		return nil
	}
	return bs
}

func FromJson(jsonStr string, obj interface{}) error {
	return Unmarshal(unsafe.StringBytes(jsonStr), obj)
}
