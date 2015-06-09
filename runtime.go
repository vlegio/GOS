package runtime

import (
	. "gotype"
	"memory"
)

const (
	sizeOfInt = 4
)

//extern __unsafe_get_addr
func pointer2byteSlice(ptr uint32) *[]byte

//extern __unsafe_get_addr
func pointer2uint32(ptr interface{}) uint32

func memset(buf_ptr uint32, value byte, count uint32) {
	var buf *[]byte
	buf = pointer2byteSlice(buf_ptr)
	for i := uint32(0); i < count; i++ {
		(*buf)[i] = value
	}
}

func memcpy(dst, src uint32, size uint32) {
	var dest, source *[]byte
	dest = pointer2byteSlice(dst)
	source = pointer2byteSlice(src)
	for i := uint32(0); i < size; i++ {
		(*dest)[i] = (*source)[i]
	}
}

func memcmp(buf1_ptr, buf2_ptr uint32, size uint32) int {
	var buf1 *[]byte
	var buf2 *[]byte
	buf1 = pointer2byteSlice(buf1_ptr)
	buf2 = pointer2byteSlice(buf2_ptr)
	for i := uint32(0); i < size; i++ {
		if (*buf1)[i] < (*buf2)[i] {
			return -1
		}
		if (*buf2)[i] < (*buf1)[i] {
			return 1
		}
	}
	return 0
}

func New(typeDescriptor uint32, size uint32) uint32 {
	buf_ptr := memory.Kmalloc(size)
	memset(buf_ptr, 0, size)
	return buf_ptr
}

func ByteArrayToString(buf_ptr uint32, length uint32) String {
	var str String
	if length == 0 {
		str.Length = 0
		return str
	}
	var res_p uint32
	res_p = New(0, length)
	memcpy(res_p, buf_ptr, length)
	str.Str = *(pointer2byteSlice(res_p))
	str.Length = uint32(length)
	return str
}

func StringPlus(str1, str2 String) String {
	if str1.Length == 0 {
		return str2
	}
	if str2.Length == 0 {
		return str1
	}
	var result String
	str1_ptr := pointer2uint32(&(str1.Str))
	str2_ptr := pointer2uint32(&(str2.Str))
	result_ptr := pointer2uint32(&(result.Str))
	memcpy(result_ptr, str1_ptr, str1.Length)
	memcpy(result_ptr+str1.Length, str2_ptr, str2.Length)
	result.Length = str1.Length + str2.Length
	return result

}

func TypeEqualIdentity(k1, k2 uint32, size uint32) bool {
	return memcmp(k1, k2, size) == 0
}

func TypeHashIdentity(key uint32, key_size uint32) uint32 {
	var ret uint32
	var i uint32
	var p *[]byte
	if key_size <= 8 {
		var u struct {
			V uint32
			A [8]byte
		}
		u.V = 0
		memcpy(pointer2uint32(&(u.A)), key, key_size)
		return ((u.V >> 32) ^ (u.V & 0xffffffff))
	}
	ret = 5381
	p = pointer2byteSlice(key)
	for i = 0; i < key_size; i++ {
		ret = ret*33 + uint32((*p)[i])
	}
	return ret
}
