package md5

import (
	"crypto/md5"
)

type Calculator struct{}

func NewCalculator() *Calculator {
	return &Calculator{}
}

func (c *Calculator) Calculate(data []byte) []byte {
	sum := md5.Sum(data)

	return sum[:]
}
