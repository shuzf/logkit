package mutate

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/qiniu/logkit/utils/models"
)

func TestKV_Transform(t *testing.T) {
	tests := []struct {
		line       []Data
		expectData []Data
		expectErr  error
		splitter   string
		new        string
	}{
		{
			expectData: []Data(nil),
			splitter:   "=",
			new:        "raw",
		},
		{
			line: []Data{{"raw": "foo=\"bar\""}},
			expectData: []Data{
				{"raw": Data{"foo": "bar"}},
			},
			splitter: "=",
			new:      "raw",
		},
		{
			line:       []Data{{"raw": "foo=\"bar\"\n"}},
			expectData: []Data{{"raw": Data{"foo": "bar"}}},
			splitter:   "=",
			new:        "raw",
		},
		{
			line: []Data{{"raw": "ts=2018-01-02T03:04:05.123Z  lvl=info msg=\"http request\" 	method=PUT\nduration=1.23 log_id=123456abc"}},
			expectData: []Data{
				{
					"raw": Data{
						"lvl":      "info",
						"msg":      "http request",
						"method":   "PUT",
						"ts":       "2018-01-02T03:04:05.123Z",
						"duration": 1.23,
						"log_id":   "123456abc",
					},
				},
			},
			splitter: "=",
			new:      "raw",
		},
		{
			line: []Data{{"raw": `ts=2018-01-02T03:04:05.123Z lvl=info  method=PUT msg="http request"`}},
			expectData: []Data{
				{
					"raw": Data{
						"lvl":    "info",
						"msg":    "http request",
						"method": "PUT",
						"ts":     "2018-01-02T03:04:05.123Z",
					},
				},
			},
			splitter: "=",
			new:      "raw",
		},
		{
			line:       []Data{{"raw": `no data.`}},
			expectData: []Data{{"raw": `no data.`}},
			expectErr:  errors.New("find total 1 erorrs in transform keyvalue, last error info is no value matched in key value transform in raw"),
			splitter:   "=",
			new:        "raw",
		},
		{
			line:       []Data{{"raw": `foo="" bar=`}},
			expectData: []Data{{"raw": `foo="" bar=`}},
			expectErr:  errors.New("find total 1 erorrs in transform keyvalue, last error info is no value matched in key value transform in raw"),
			splitter:   "=",
			new:        "raw",
		},
		{
			line:       []Data{{"raw": `abc=abc foo="def`}},
			expectData: []Data{{"raw": `abc=abc foo="def`}},
			expectErr:  errors.New("find total 1 erorrs in transform keyvalue, last error info is logfmt syntax error at pos 17 on line 1: unterminated quoted value"),
			splitter:   "=",
			new:        "raw",
		},
		{
			line:       []Data{{"raw": `"foo=" bar=abc`}},
			expectData: []Data{{"raw": `"foo=" bar=abc`}},
			expectErr:  errors.New(`find total 1 erorrs in transform keyvalue, last error info is logfmt syntax error at pos 1 on line 1: unexpected '"'`),
			splitter:   "=",
			new:        "raw",
		},
		{
			line: []Data{{"raw": "ts:2018-01-02T03:04:05.123Z lvl:info msg:\"http request\" method:PUT\nduration:1.23 log_id:123456abc"}},
			expectData: []Data{
				{
					"raw":      "ts:2018-01-02T03:04:05.123Z lvl:info msg:\"http request\" method:PUT\nduration:1.23 log_id:123456abc",
					"lvl":      "info",
					"msg":      "http request",
					"method":   "PUT",
					"ts":       "2018-01-02T03:04:05.123Z",
					"duration": 1.23,
					"log_id":   "123456abc",
				},
			},
			splitter: ":",
		},
	}

	k := &KV{
		Key: "raw",
	}

	for _, test := range tests {
		k.New = test.new
		k.Init()
		k.Splitter = test.splitter
		actual, err := k.Transform(test.line)
		assert.EqualValues(t, test.expectErr, err)
		assert.EqualValues(t, test.expectData, actual)
	}
}
