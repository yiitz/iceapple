package api

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/yiitz/iceapple/log"
	"github.com/yiitz/iceapple/storage"
	"github.com/yiitz/persistent-cookiejar"
	"math/big"
	"net/http"
	"strings"
	"net/url"
	"io/ioutil"
	"github.com/yiitz/iceapple/config"
)

var client http.Client

var logger = log.NewLogger("api")

const httpRoot = "http://music.163.com"

const modulus = "00e0b509f6259df8642dbc35662901477df22677ec152b5ff68ace615bb7b725152b3ab17a876aea8a5aa76d2e417629ec4ee341f56135fccf695280104e0312ecbda92557c93870114af6c9d05c4f7f0c3685b7a46bee255932575cce10b424d813cfe4875d3e82047b97ddef52741d546b8e289dc6935b3ece0462db0a22b8e7"
const nonce = "0CoJUm6Qyw8W8jud"
const pubKey = "010001"

const iv = "0102030405060708"

func get(uri string) map[string]interface{} {
	return decodeResp(client.Get(uri))
}

func post(uri string, data map[string]interface{}) map[string]interface{} {

	var text string
	if data == nil {
		data = map[string]interface{}{}
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		logger.Error(err)
		return nil
	}
	text = string(bytes)

	secKey := newRandom(16)
	encText := aesEncrypt(aesEncrypt(text, nonce), secKey)
	encSecKey := rsaEncrypt(secKey)

	req, err := http.NewRequest("POST", uri, strings.NewReader(fmt.Sprintf("params=%s&encSecKey=%s", url.QueryEscape(encText), url.QueryEscape(encSecKey))))
	if err != nil {
		logger.Error(err)
		return nil
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8,gl;q=0.6,zh-TW;q=0.4")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", "http://music.163.com")
	req.Header.Add("Host", "music.163.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.152 Safari/537.36")

	logger.Debugf("http request\nuri:%s\nbody:%s,%s", uri, encText, encSecKey)

	return decodeResp(client.Do(req))
}

func decodeResp(resp *http.Response, err error) map[string]interface{} {
	if err != nil {
		logger.Error(err)
		return nil
	}
	r := make(map[string]interface{})
	rs, err := ioutil.ReadAll(resp.Body)
	logger.Debugf("http response\nheader:%+v\nbody:%s", resp.Header, string(rs))
	err = json.Unmarshal(rs, &r)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return r
}

func aesEncrypt(text string, secKey string) string {
	pad := 16 - len(text)%16
	text = text + strings.Repeat(string(byte(pad)), pad)
	plaintext := []byte(text)
	if len(plaintext)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}
	block, err := aes.NewCipher([]byte(secKey))
	if err != nil {
		panic(err)
	}

	cipherText := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(cipherText, plaintext)
	return base64.StdEncoding.EncodeToString(cipherText)
}

func rsaEncrypt(text string) string {
	text = reverse(text)
	a := parseBigInt(hex.EncodeToString([]byte(text)))
	b := parseBigInt(pubKey)
	c := parseBigInt(modulus)
	a = a.Exp(a, b, c)
	return fmt.Sprintf("%0256x", a)
}

func parseBigInt(s string) *big.Int {
	i := new(big.Int)
	i, ok := i.SetString(s, 16)
	if !ok {
		log.LoggerRoot.Warnf("big int SetString error: %s", hex.EncodeToString([]byte(s)))
	}
	return i
}

func reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

func InitClient() {
	jar, err := cookiejar.New(&cookiejar.Options{Filename: storage.AppDir() + "/.cookie"})
	if err != nil {
		panic(err)
	}
	client.Jar = jar
	if len(config.Proxy) > 0 {
		client.Transport = &http.Transport{Proxy: func(request *http.Request) (*url.URL, error) {
			return url.Parse(config.Proxy)
		}}
	}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if strings.HasSuffix(req.URL.Host, "music.126.net") {
			return http.ErrUseLastResponse
		}
		return nil
	}
}

func SaveCookie() {
	client.Jar.(*cookiejar.Jar).Save()
}

var stdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// https://github.com/dchest/uniuri/blob/master/uniuri.go
func newRandom(length int) string {
	return newLenChars(length, stdChars)
}

// NewLenChars returns a new random string of the provided length, consisting of the provided byte slice of allowed characters(maximum 256).
func newLenChars(length int, chars []byte) string {
	if length == 0 {
		return ""
	}
	clen := len(chars)
	if clen < 2 || clen > 256 {
		panic("Wrong charset length for NewLenChars()")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("Error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				continue // Skip this number to avoid modulo bias.
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}
