package Random

import (
	"math/rand"
)

const letterBytes ="abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//生成随机UserAgent
func RandomString() string  {
	//rand.Seed(time.Now().Unix())
	b:=make([]byte,rand.Intn(10)+10)
	for i:=range b{
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}