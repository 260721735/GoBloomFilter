# GoBloomFilter
bloom filter for golang
```
package main

import (
	"github.com/260721735/GoBloomFilter/BloomFilter"
	"log"
)
func main() {
	total := uint64(2000000) //max total
	Probability := 0.0000000001
	bf, err := BloomFilter.Create(total)
	if err != nil {
		log.Println(err.Error())
	}
	bf.Put("BloomFilter")
	log.Println(bf.MightContain("BloomFilter"))
	//create with fpp
	bf, err = BloomFilter.CreateWithFPP(total, Probability)
	if err != nil {
		log.Println(err.Error())
	}
	bf.Put("BloomFilter")
	log.Println(bf.MightContain("BloomFilter"))
}

```
