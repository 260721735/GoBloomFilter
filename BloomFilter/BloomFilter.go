package BloomFilter

import (
	"encoding/binary"
	"errors"
	"strconv"
)

const (
	defaultfalsePositiveProbability = 0.03
	defaultDataLen                  = 2000000
	maxElement                      = 20000000
	//maxNumOfBits                    = 1000000000
)

type BloomFilter struct {
	bitCount uint64
	data     [][defaultDataLen]uint64
	//data               []uint64
	numOfBits          uint64
	numOfHashFunctions uint64
}

//ln(falsePositiveProbability)=len/expectedInsertions
func CreateWithFPP(expectedInsertions uint64, falsePositiveProbability float64) (BloomFilter, error) {
	return create(expectedInsertions, falsePositiveProbability)
}
func Create(expectedInsertions uint64) (BloomFilter, error) {
	return create(expectedInsertions, defaultfalsePositiveProbability)
}
func create(expectedInsertions uint64, falsePositiveProbability float64) (BloomFilter, error) {
	if expectedInsertions > maxElement {
		return BloomFilter{}, errors.New("maxElement is " + strconv.Itoa(maxElement))
	}
	var err error
	b := BloomFilter{}
	b.numOfBits, err = optimalNumOfBits(expectedInsertions, falsePositiveProbability)
	if err != nil {
		return b, err
	}
	b.numOfHashFunctions = optimalNumOfHashFunctions(expectedInsertions, b.numOfBits)
	for {
		b.data = append(b.data, [defaultDataLen]uint64{0})
		if uint64(len(b.data)*defaultDataLen) >= b.numOfBits {
			break
		}
		//log.Println("NumOfBits", b.NumOfBits, "now len", len(b.data)*defaultDataLen)
	}
	//for i := uint64(0); i < b.NumOfBits; i++ {
	//	b.data = append(b.data, 0)
	//}
	return b, err
}
func (b *BloomFilter) Put(res interface{}) {
	switch res.(type) {
	case string:
		b.put([]byte(res.(string)))
		break
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		b.put(intToBytes(res.(int64)))
	}
}
func (b *BloomFilter) MightContain(res interface{}) bool {
	switch res.(type) {
	case string:
		return b.mightContain([]byte(res.(string)))
		break
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		return b.mightContain(intToBytes(res.(int64)))
	}
	return false
}
func intToBytes[T intall](i T) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

type intall interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}
