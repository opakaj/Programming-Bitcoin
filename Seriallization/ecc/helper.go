package ecc

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"unsafe"

	"github.com/btcsuite/btcutil"
)

var BASE58ALPHABET = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func hash160(s string) string {
	return hex.EncodeToString(btcutil.Hash160([]byte(s)))
}

func hash256(s string) string {
	//two rounds of sha256
	hash := sha256.Sum256([]byte(s))
	hash = sha256.Sum256(hash[:])
	return hex.EncodeToString(hash[:]) //convet to [] by slicicng it
}

func divmod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}

func ByteArrayToInt(arr []byte) int64 {
	val := int64(0)
	size := len(arr)
	for i := 0; i < size; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}

func encodeBase58(s string) string {
	count := 0
	for _, c := range s {
		if c == 0 {
			count += 1
		} else {
			break
		}
	}
	num := binary.BigEndian.Uint32([]byte(s)) //bytes to int
	prefix := strings.Repeat("1", count)
	result := ""
	for num > 0 {
		_, mod := divmod(int64(num), 58)
		result = string(BASE58ALPHABET[mod]) + result
	}
	return prefix + result
}

func encodeBase58Checksum(b string) string {
	return encodeBase58(b + hash256(b)[:4])
}

func decodeBase58(s string) string {
	num := 0
	for _, c := range s {
		num *= 58
		num += func() int {
			for i, val := range BASE58ALPHABET {
				if val == c {
					return i
				}
			}
			panic(errors.New("ValueError: element not found"))
		}()
	}
	combined := make([]byte, 25)
	binary.BigEndian.PutUint64(combined, uint64(num)) //int to bytes
	checksum := combined[len(combined)-4:]
	if hash256(string(combined[:len(combined)-4]))[:4] != string(checksum) {
		panic(
			fmt.Errorf("bad address: %s %s", checksum, hash256(string(combined[:len(combined)-4]))[:4]),
		)
	}
	return string(combined[1 : len(combined)-4])
}
