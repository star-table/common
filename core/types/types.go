package types

import (
	"fmt"
	"github.com/star-table/common/core/consts"
	"io"
	"strings"
	"time"
)

type Time time.Time

func NowTime() Time {
	return Time(time.Now())
}

func Time0() Time {
	return Time(time.Unix(0, 0))
}

func (t *Time) IsNotNull() bool {
	return time.Time(*t).After(consts.BlankTimeObject)
}

func (t *Time) IsNull() bool {
	return !t.IsNotNull()
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	//fmt.Println("UnmarshalJSON:", string(data))
	date := strings.ReplaceAll(string(data), "\"", "")
	if date == "" {
		return
	}
	now, err := time.Parse(consts.AppTimeFormat, date)
	if err != nil {
		now, err = time.Parse(consts.AppSystemTimeFormat, date)
		if err != nil {
			now, err = time.Parse(consts.AppSystemTimeFormat8, date)
			if err != nil {
				*t = Time(consts.BlankTimeObject)
				return
			}
		}
	}
	*t = Time(now)
	//fmt.Println("UnmarshalJSON:", t)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	//fmt.Println("MarshalJSON:", t)
	b := make([]byte, 0, len(consts.AppTimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, consts.AppTimeFormat)
	b = append(b, '"')
	//fmt.Println("MarshalJSON:", string(b))
	return b, nil
}

// UnmarshalGQL implements the graphql.Marshaler interface
func (t *Time) UnmarshalGQL(v interface{}) error {
	date, ok := v.(string)

	date = strings.ReplaceAll(date, "\"", "")
	if date == "" {
		return nil
	}
	//fmt.Println("UnmarshalGQL:", date)
	if !ok {
		return fmt.Errorf("points must be strings")
	}
	now, err := time.Parse(consts.AppTimeFormat, date)
	if err != nil {
		now, err = time.Parse(consts.AppSystemTimeFormat, date)
		if err != nil {
			now, err = time.Parse(consts.AppSystemTimeFormat8, date)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	*t = Time(now)
	//fmt.Println("UnmarshalGQL:", t)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (t Time) MarshalGQL(w io.Writer) {
	//fmt.Println("MarshalGQL:", t)
	b := make([]byte, 0, len(consts.AppTimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, consts.AppTimeFormat)
	b = append(b, '"')
	//fmt.Println("MarshalGQL:", string(b))
	w.Write(b)
}

func (t Time) String() string {
	return time.Time(t).Format(consts.AppTimeFormat)
}
