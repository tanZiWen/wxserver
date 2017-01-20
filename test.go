package main

import (
	"fmt"
	"golang.org/x/crypto/scrypt"
	"encoding/base64"
	"errors"
	"math/rand"
	"time"
)

func Afuntion(ch chan int) {
	fmt.Println("finish")
	<-ch
}

func main() {
	passrd, err := EncryptPassword("271226")
	fmt.Println(passrd, err);

	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6}

	fmt.Println(m)



	//c := make(chan int, 3)

	go func(p map[string]int) {
		p["a"] = 10
	}(m)

	go func(p map[string]int) {
		p["a"] = 11
	}(m)

	go func(p map[string]int) {
		p["a"] = 12
	}(m)

	fmt.Println(m)

	a := []int{1, 2}

	myAppends(a)

	fmt.Println(a)

	s := make(map[string]string)
	s["foo"] = "foo"

	maptest(s)

	fmt.Println(s)

	timeout := make (chan bool, 5)
	go func() {
		time.Sleep(5e9) // sleep one second
		timeout <- true
	}()
	ch := make (chan int, 5)
	ch <- 5
	ch <- 4
	select {
	case i := <- ch:
		fmt.Println(i)
	case <- timeout:
		fmt.Println("timeout!")
	}

	if err := returnsError(); err.(*MyError) != nil {
		panic(err)
	}
	quenue := make(chan int, 10)
	quenue <- 1
	quenue <- 2
	quenue <- 3
	close(quenue)
	for elem := range quenue {
		fmt.Println(elem)
	}


	go print3()
	go print2()
	go print1()

	for i:=0;i<3;i++ {
		fmt.Println(<- quit)
	}

	//c := make(chan int)
	//
	//c <- 1
	//
	//fmt.Println(<- c)

	//
	//c := make(chan int)
	//
	//go func() {
	//	fmt.Println("goroutine ok")
	//	c <- 1
	//	fmt.Println("我会被阻塞")
	//	close(c)
	//}()
	//
	//fmt.Println(<-c)
	//fmt.Println("main ok")

	c := make(chan int, 1)
	go func(){
		c <- 1
		fmt.Println("第一次发送成功")
		c <- 2
		fmt.Println("第二次发送成功")
		c <- 3
		fmt.Println("第三次发送成功")
		c <- 4
		fmt.Println("第四次发送成功")
	}()

	fmt.Println("准备接收")
	k := <- c
	fmt.Println("第一次接收", k)
	k = <- c
	fmt.Println("第二次接收", k)
	k = <- c
	fmt.Println("第三次接收", k)
	k = <- c
	fmt.Println("第四次接收", k)
}

var quit chan int = make(chan int, 3000)

func print1() {
	quit <- 1
}
func print2() {
	quit <- 2
}
func print3() {
	quit <- 3
}

type MyError struct {}
func (* MyError) Error() string {
	return "MyError"
}
func returnsError() error {
	var p *MyError = nil
	return p
}


func maptest(m map[string]string) {
	m["foo"] = "bar"
}

func myAppends(s []int) {
	_ = append(s, 3)
}

const (
	ENCRYPT_ALGORITHM = "pbkdf2"
	ENCRYPT_HASHER = "sha256"
	ENCRYPT_TIMES = 16384
)

func EncryptPassword(password string, salts ...string) (encryptedPassword string, err error) {
	var salt string

	if len(salts) > 0 {
		salt = salts[0]

		if len(salt) != 6 {
			salt = RandomStr(6)
		}
	} else {
		salt = RandomStr(6)
	}

	dk, err := scrypt.Key([]byte(password), []byte(salt), ENCRYPT_TIMES, 8, 1, 32)

	if err != nil {
		return "", err
	}

	encryptedPassword = fmt.Sprintf("%s_%s_%s_%d_%s", ENCRYPT_ALGORITHM, ENCRYPT_HASHER, salt, ENCRYPT_TIMES, base64.StdEncoding.EncodeToString(dk))
	return encryptedPassword, nil
}

var charset = []byte{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
	'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
	'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
}

const (
	UPPER_CHARSET = iota
	LOWER_CHARSET
	INT_CHARSET
	UPPER_LOWER_CHARSET
	FULL_CHARSET
)

func RandomSpecStr(num int, set uint) string {
	var subcharset []byte

	switch(set) {
	case UPPER_CHARSET: subcharset = charset[:25]
	case LOWER_CHARSET: subcharset = charset[26:51]
	case INT_CHARSET: subcharset = charset[52:]
	case UPPER_LOWER_CHARSET: subcharset = charset[:51]
	case FULL_CHARSET: subcharset = charset
	default: panic(errors.New("Unknown charset identify:"+ string(set)))
	}

	var buf = make([]byte, num)
	for i := 0; i < num; i++ {
		index := rand.Intn(len(subcharset))
		buf[i] = subcharset[index]
	}
	return string(buf)
}

func RandomStr(num int) string {
	return RandomSpecStr(num, FULL_CHARSET)
}