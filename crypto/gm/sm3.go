package gm

import "github.com/tjfoc/gmsm/sm3"

//SM3Hash 加密算法
func SM3Hash(msg []byte) []byte {
	c := sm3.New()
	c.Write(msg)
	return c.Sum(nil)
}
