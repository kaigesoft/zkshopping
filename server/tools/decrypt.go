package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unsafe"
)

var key = []byte{12, 7, 21, 9, 8, 21, 12, 15, 7, 8, 1, 84, 95, 87, 84, 87}

func init() {
	for i := range key {
		key[i] = key[i] ^ 0x66
	}
}

//DecDecrypt dec
func DecDecrypt(cipherstring string) (string, error) {
	ciphertext, err := hex.DecodeString(cipherstring)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCBCDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.CryptBlocks(ciphertext, ciphertext)
	return string(bytes.TrimRight(ciphertext, string([]byte{0}))), nil
}

//Encrypt enc
func Encrypt(plainstring string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext := []byte(plainstring)
	if len(plaintext)%aes.BlockSize != 0 {
		plaintext = append(plaintext, bytes.Repeat([]byte{0}, aes.BlockSize-len(plaintext)%aes.BlockSize)...)
	}
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCBCEncrypter(block, iv)
	stream.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.
	return hex.EncodeToString(ciphertext), nil
}

//AuthToken token
func AuthToken() string {
	userid := "userid=ringliu"
	time := String(time.Now().Unix())
	tt := "timestamp=" + time
	m := "37505d287a" + time
	token := "token=" + md5v1(m)
	return userid + string("&") + string(tt) + string("&") + token
}

func md5v1(str string) string {
	data := []byte(str)
	ha := md5.Sum(data)
	md5str := fmt.Sprintf("%x", ha)
	return md5str
}

func GetLocalIPAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		ipaddr, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			continue
		}
		if ipaddr.IsLoopback() {
			continue
		}
		if ipaddr.To4() != nil {
			if runtime.GOOS == "darwin" {
				if !strings.HasPrefix(ipaddr.String(), "192") {
					continue
				}
			}
			return ipaddr.String()
		}
	}
	return ""
}

//FilteredSQLInject 正则过滤sql注入的方法
func FilteredSQLInject(toMatchStr string) bool {
	str := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	re, err := regexp.Compile(str)
	if err != nil {
		return false
	}
	return re.MatchString(toMatchStr)
}

//IsNumeric 验证数字类型
func IsNumeric(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128:
		return true
	case string:
		str := val.(string)
		if len(str) == 0 {
			return false
		}

		str = strings.Trim(str, " \t\r\n\v\f")
		if len(str) == 0 {
			return false
		}

		if str[0] == '-' || str[0] == '+' {
			if len(str) == 1 {
				return false
			}
			str = str[1:]
		}

		if len(str) > 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X') {
			for _, h := range str[2:] {
				if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'f') || (h >= 'A' && h <= 'F')) {
					return false
				}
			}
			return true
		}
		// 0-9
		p, s, l := 0, 0, len(str)
		for i, v := range str {
			if v == '.' {
				if p > 0 || s > 0 || i+1 == l {
					return false
				}
				p = i
			} else if v == 'e' || v == 'E' {
				if i == 0 || s > 0 || i+1 == l {
					return false
				}
				s = i
			} else if v < '0' || v > '9' {
				return false
			}
		}
		return true
	}
	return false
}

// MD5 加密各种类型数据
func MD5(source interface{}) (md5str string) {
	sourceLen := unsafe.Sizeof(source) // nolint
	data := make([]byte, sourceLen)
	switch source.(type) {
	case string:
		data = []byte(source.(string))
		break
	case int16:
		binary.BigEndian.PutUint16(data, uint16(source.(int16)))
		break
	case int32:
		binary.BigEndian.PutUint32(data, uint32(source.(int32)))
		break
	case int64:
		binary.BigEndian.PutUint64(data, uint64(source.(int64)))
		break
	case uint16:
		binary.BigEndian.PutUint16(data, source.(uint16))
		break
	case uint32:
		binary.BigEndian.PutUint32(data, source.(uint32))
		break
	case uint64:
		binary.BigEndian.PutUint64(data, source.(uint64))
		break
	case int, int8:
	case uint, uint8:
	case float32, float64, complex64, complex128:
		data = IntToByte(source.(int), sourceLen)
		break
	}
	md5str = fmt.Sprintf("%x", md5.Sum(data))
	return
}

// IntToByte 整型转字节数组
func IntToByte(data int, len uintptr) (ret []byte) {
	ret = make([]byte, len)
	var tmp = 0xff
	var index uint
	for index = 0; index < uint(len); index++ {
		ret[index] = byte((tmp << (index * 8) & data) >> (index * 8))
	}
	return ret
}
