package crypto

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/shomali11/util/xhashes"
	"golang.org/x/crypto/bcrypt"
	"io"
	r2 "math/rand"
	"strconv"
	"time"
)

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//密码加密
func GeneratePassword(password string) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPwd), nil
}

//密码验证
func ComparePassword(hashPassword string, password string) bool {
	// 进行密码验证
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		return false
	}
	return true
}

// 生成随机编码
func GenerateSn(prefix ...string) string {
	sn := xhashes.FNV64(UniqueId())
	snStr := strconv.FormatUint(sn, 10)
	if prefix != nil {
		snStr = prefix[0] + snStr
	}
	return snStr
}

// 生成Guid字串
func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return MD5(base64.URLEncoding.EncodeToString(b))
}

func CreateCaptcha() string {
	return fmt.Sprintf("%06v", r2.New(r2.NewSource(time.Now().UnixNano())).Int31n(1000000))
}
