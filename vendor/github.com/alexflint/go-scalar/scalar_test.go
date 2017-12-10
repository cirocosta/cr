package scalar

import (
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertParse(t *testing.T, expected interface{}, str string) {
	v := reflect.New(reflect.TypeOf(expected)).Elem()
	err := ParseValue(v, str)
	if assert.NoError(t, err) {
		assert.Equal(t, expected, v.Interface())
	}

	ptr := reflect.New(reflect.PtrTo(reflect.TypeOf(expected))).Elem()
	err = ParseValue(ptr, str)
	if assert.NoError(t, err) {
		assert.Equal(t, expected, ptr.Elem().Interface())
	}

	assert.True(t, CanParse(v.Type()))
	assert.True(t, CanParse(ptr.Type()))
}

func TestParseValue(t *testing.T) {
	// strings
	assertParse(t, "abc", "abc")

	// booleans
	assertParse(t, true, "true")
	assertParse(t, false, "false")

	// integers
	assertParse(t, int(123), "123")
	assertParse(t, int8(123), "123")
	assertParse(t, int16(123), "123")
	assertParse(t, int32(123), "123")
	assertParse(t, int64(123), "123")

	// unsigned integers
	assertParse(t, uint(123), "123")
	assertParse(t, byte(123), "123")
	assertParse(t, uint8(123), "123")
	assertParse(t, uint16(123), "123")
	assertParse(t, uint32(123), "123")
	assertParse(t, uint64(123), "123")
	assertParse(t, uintptr(123), "123")
	assertParse(t, rune(123), "123")

	// floats
	assertParse(t, float32(123), "123")
	assertParse(t, float64(123), "123")

	// durations
	assertParse(t, 3*time.Hour+15*time.Minute, "3h15m")

	// IP addresses
	assertParse(t, net.IPv4(1, 2, 3, 4), "1.2.3.4")

	// MAC addresses
	assertParse(t, net.HardwareAddr("\x01\x23\x45\x67\x89\xab"), "01:23:45:67:89:ab")

	// MAC addresses
	assertParse(t, net.HardwareAddr("\x01\x23\x45\x67\x89\xab"), "01:23:45:67:89:ab")
}

func TestParse(t *testing.T) {
	var v int
	err := Parse(&v, "123")
	require.NoError(t, err)
	assert.Equal(t, 123, v)
}
