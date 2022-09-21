package BloomFilter

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	mathrand "math/rand"
	"sync"
	"testing"
	"time"
)

var rander = rand.Reader // random function

func TestConcurrency(t *testing.T) {
	mathrand.Seed(time.Now().UnixNano())
	conCount := 100
	elementAll := getElements(100)
	Probability := 0.00000001
	bf, err := CreateWithFPP(uint64(200000), Probability)
	if err != nil {
		t.Fatal(err)
	}
	wg := new(sync.WaitGroup)
	for i := 0; i < conCount; i++ {
		wg.Add(1)
		go func(index int, elements []string, wg *sync.WaitGroup, bf *BloomFilter) {
			for {
				randi := mathrand.Intn(len(elements))
				randj := mathrand.Intn(len(elements))
				bf.Put(elements[randi])
				log.Println(index, "now put success")
				if !bf.MightContain(elements[randj]) {
					log.Println("not contain")
				}
				log.Println(index, "now get success")
			}

		}(i, elementAll, wg, &bf)
	}
	wg.Wait()
}

//仅测试性能
// create elements success with 3000000
// create bf with 1.643833ms
// put bf with 297.341583ms
// mightContain bf with 512.716916ms
// success num is  3000000
// error num is  0
// error percent is  0
func TestValidity(t *testing.T) {
	elementCount := 3000000
	elementAll := getElements(elementCount)
	log.Println("create elements success with", len(elementAll))
	elements := elementAll[0 : len(elementAll)/3*2]
	othersElements := elementAll[len(elementAll)/3*2:]
	total := uint64(2000000)
	Probability := 0.00000001

	beforeTime := time.Now()
	bf, err := CreateWithFPP(total, Probability)
	//bf, err := Create(total) //0.009996466666666667
	log.Println("create bf with", time.Since(beforeTime))
	if err != nil {
		t.Fatal(err)
	}

	beforeTime = time.Now()
	for i := range elements {
		bf.Put(elements[i])
	}
	log.Println("put bf with", time.Since(beforeTime))
	successNum := 0
	errorNum := 0
	beforeTime = time.Now()
	for i := range elements {
		if bf.MightContain(elements[i]) {
			successNum++
		} else {
			errorNum++
		}
	}
	log.Println("mightContain bf with", time.Since(beforeTime))

	for i := range othersElements {
		if bf.MightContain(othersElements[i]) {
			errorNum++
		} else {
			successNum++
		}
	}
	log.Println("success num is ", successNum)
	log.Println("error num is ", errorNum)
	log.Println("error percent is ", float64(errorNum)/float64(successNum+errorNum))
}
func getElements(num int) []string {
	maps := make(map[string]int)
	for len(maps) != num {
		id := NewRandom()
		if _, ok := maps[id]; !ok {
			maps[id] = 1
		}
	}
	elements := make([]string, num)
	elementsIndex := 0
	for k, _ := range maps {
		elements[elementsIndex] = k
		elementsIndex++
	}
	return elements
}

func NewRandom() string {
	var uuid [16]byte
	io.ReadFull(rander, uuid[:])
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	//var buf []byte
	buf := make([]byte, 36)
	hex.Encode(buf, uuid[:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], uuid[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], uuid[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], uuid[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], uuid[10:])
	return string(buf[:])
}
