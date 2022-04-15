package main

import (
	"github.com/260721735/GoBloomFilter/BloomFilter"
	"log"
	"strconv"
	"time"
)

func main() {
	test1()
}
func test1() {

	total := uint64(2000000)
	Probability := 0.0000000001
	skip := 1
	beforeTime := time.Now()
	bf, err := BloomFilter.Create(total)
	log.Println("createtime", time.Since(beforeTime))
	if err != nil {
		log.Println(err.Error())
	}
	beforeTime = time.Now()
	for i := 0; uint64(i) < total; i += skip {
		bf.Put(strconv.Itoa(i))
	}
	log.Println("puttime", time.Since(beforeTime))
	// 判断值是否存在过滤器中
	beforeTime = time.Now()
	count := 0
	for i := 0; uint64(i) < total*2; i += skip {
		if bf.MightContain(strconv.Itoa(i)) {
			count++
		}
	}
	log.Println("gettime", time.Since(beforeTime))
	log.Println("已匹配数量", count)
	beforeTime = time.Now()
	bf, err = BloomFilter.CreateWithFPP(total, Probability)
	log.Println("指定误判率createtime", time.Since(beforeTime))
	if err != nil {
		log.Println(err.Error())
	}
	beforeTime = time.Now()
	for i := 0; uint64(i) < total; i += skip {
		bf.Put(strconv.Itoa(i))
	}
	log.Println("指定误判率puttime", time.Since(beforeTime))
	// 判断值是否存在过滤器中
	beforeTime = time.Now()
	count = 0
	for i := 0; uint64(i) < total*2; i += skip {
		if bf.MightContain(strconv.Itoa(i)) {
			count++
		}
	}
	log.Println("指定误判率gettime", time.Since(beforeTime))
	log.Println("指定误判率已匹配数量", count)
}
func test2() {
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
