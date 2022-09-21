package BloomFilter

import (
	"errors"
	"math"
)

func optimalNumOfBits(elementNum uint64, falsePositiveProbability float64) (uint64, error) {
	if elementNum < 0 {
		return 0, errors.New("the number of elements must be greater than zero")
	}
	if falsePositiveProbability <= 0 {
		falsePositiveProbability = defaultfalsePositiveProbability
	}
	//log.Println(float64(elementNum) * math.Log(falsePositiveProbability))
	//log.Println(math.Log(2) * math.Log(2))
	//log.Println(float64(elementNum) * math.Log(falsePositiveProbability) / (math.Log(2) * math.Log(2)))
	res := float64(elementNum) * math.Log(falsePositiveProbability) / (math.Log(2) * math.Log(2))
	if res <= 0 {
		res = res * -1
	}
	//log.Println(uint64(res))
	return uint64(res), nil
}
func optimalNumOfHashFunctions(elementNum uint64, bitslen uint64) uint64 {
	return uint64(math.Max(1.0, math.Round(float64(bitslen)/float64(elementNum)*math.Log(2))))
}

func (b *BloomFilter) put(element []byte) {
	hash64 := murmurHash64A(element)
	hash1 := hash64
	hash2 := hash64 >> 32
	for i := uint64(1); i <= b.numOfHashFunctions; i++ {
		nextHash := hash1 + int64(i)*hash2
		if nextHash < 0 {
			nextHash = ^nextHash
		}
		b.set(uint64(nextHash) % b.numOfBits)
	}
}
func (b *BloomFilter) set(index uint64) {
	b.data[index>>6] |= uint64(1) << (index % 64)
	//b.data[int(index>>6)/defaultDataLen%len(b.data)][int(index>>6)%defaultDataLen] |= uint64(1) << (index % 64)
}

func (b BloomFilter) mightContain(element []byte) bool {
	hash64 := murmurHash64A(element)
	hash1 := hash64
	hash2 := hash64 >> 32
	for i := uint64(1); i <= b.numOfHashFunctions; i++ {
		nextHash := hash1 + int64(i)*hash2
		if nextHash < 0 {
			nextHash = ^nextHash
		}

		if !b.get(uint64(nextHash) % b.numOfBits) {
			return false
		}
		//if getBit(b.data[uint64(nextHash)%b.NumOfBits>>6], uint64(nextHash)%b.NumOfBits%64) != 1 {
		//	return false
		//}
	}
	return true
}

func (b BloomFilter) get(index uint64) bool {
	return getBit(b.data[index>>6], index%64) == 1
}

//func (b BloomFilter) get(index uint64) bool {
//	return getBit(b.data[int(index>>6)/defaultDataLen%len(b.data)][int(index>>6)%defaultDataLen], index%64) == 1
//}
func getBit(c uint64, i uint64) uint64 {
	//b.data[index>>6] |= uint64(1) << (index % 64)
	return c >> i & 0x1
	// return (this.data[index >> 6] & 1L << index) != 0L;
	//return c & (0x1 << i)
}

const (
	BIG_M = 0xc6a4a7935bd1e995
	BIG_R = 47
	SEED  = 0x1234ABCD
)

func murmurHash64A(data []byte) (h int64) {
	var k int64
	h = SEED ^ int64(uint64(len(data))*BIG_M)

	var ubigm uint64 = BIG_M
	var ibigm = int64(ubigm)
	for l := len(data); l >= 8; l -= 8 {
		k = int64(int64(data[0]) | int64(data[1])<<8 | int64(data[2])<<16 | int64(data[3])<<24 |
			int64(data[4])<<32 | int64(data[5])<<40 | int64(data[6])<<48 | int64(data[7])<<56)

		k := k * ibigm
		k ^= int64(uint64(k) >> BIG_R)
		k = k * ibigm

		h = h ^ k
		h = h * ibigm
		data = data[8:]
	}

	switch len(data) {
	case 7:
		h ^= int64(data[6]) << 48
		fallthrough
	case 6:
		h ^= int64(data[5]) << 40
		fallthrough
	case 5:
		h ^= int64(data[4]) << 32
		fallthrough
	case 4:
		h ^= int64(data[3]) << 24
		fallthrough
	case 3:
		h ^= int64(data[2]) << 16
		fallthrough
	case 2:
		h ^= int64(data[1]) << 8
		fallthrough
	case 1:
		h ^= int64(data[0])
		h *= ibigm
	}

	h ^= int64(uint64(h) >> BIG_R)
	h *= ibigm
	h ^= int64(uint64(h) >> BIG_R)
	return
}
