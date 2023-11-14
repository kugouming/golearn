package utils

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"net"
)

// GetMd5 获取key的md5
func GetMd5(key string) []byte {
	d5 := md5.New()
	d5.Write([]byte(key))
	return d5.Sum(nil)
}

// GetHexString 获取16进制字符串
func GetHexString(b []byte) string {
	return hex.EncodeToString(b)
}

// GetKeyIv 通过提供的加密字符串通过md5计算出key iv
func GetKeyIv(passwd string) (key []byte, iv []byte) {
	var pb = []byte(passwd)
	var split = len(passwd) / 2
	var d5 = md5.New()
	d5.Write(pb[:split])
	key = d5.Sum(nil)
	d5 = md5.New()
	d5.Write(pb[split:])
	iv = d5.Sum(nil)
	return key, iv
}

func Recover() {
	if err := recover(); err != nil {
		log.Printf("Recover err:%+v", err)
	}
}

// WCopy 写的一端加密，读不加密
func WCopy(dst *NCopy, src net.Conn) {
	defer func() {
		src.Close()
		dst.Close()
	}()
	buf := make([]byte, 10240)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			(*dst).Write(buf[:n])
		}
		if err != nil {
			return
		}
	}
}

// RCopy 读的一端解密，写不加密
func RCopy(dst net.Conn, src *NCopy) {
	defer func() {
		src.Close()
		dst.Close()
	}()
	buf := make([]byte, 10240)
	for {
		n, err := (*src).Read(buf)
		if n > 0 {
			dst.Write(buf[:n])
		}
		if err != nil {
			return
		}
	}
}

// NetCopy 流复制处理
func NetCopy(dst, src net.Conn, msg string) {
	defer func() {
		src.Close()
		dst.Close()
	}()
	buf := make([]byte, 10240)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			dst.Write(buf[:n])
		}
		if err != nil {
			return
		}
	}
}
