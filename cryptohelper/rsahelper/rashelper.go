package rashelper

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"strings"
)

type RsaHelper struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func New() *RsaHelper {
	return &RsaHelper{}
}

// 初始化公/私钥
func (s *RsaHelper) GenerateKey() *rsa.PrivateKey {
	s.privateKey, _ = rsa.GenerateKey(rand.Reader, 2048)
	s.publicKey = &s.privateKey.PublicKey
	return s.privateKey
}

// 获取公/私钥对
func (s *RsaHelper) QueryKeyPair() *KeyPair {
	if s.privateKey == nil {
		s.GenerateKey()
	}

	p := &KeyPair{}

	p.Private = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(s.privateKey),
	})

	p.Public = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(s.publicKey),
	})

	return p
}

// 设置私钥
func (s *RsaHelper) SetPrivateKey(privateKey *rsa.PrivateKey) {
	s.privateKey = privateKey
}

// 设置公钥
func (s *RsaHelper) SetPublicKey(publicKey *rsa.PublicKey) {
	s.publicKey = publicKey
}

// 设置私钥(base64 string)
func (s *RsaHelper) SetPrivateKeyStr(base64Str string) (*rsa.PrivateKey, error) {
	block, err := s.decodeKey([]byte(base64Str))
	if err != nil {
		return nil, err
	}

	return s.toPrivateKey(block)
}

// 设置公钥(base64 string)
func (s *RsaHelper) SetPublicKeyStr(base64Str string) (*rsa.PublicKey, error) {
	block, err := s.decodeKey([]byte(base64Str))
	if err != nil {
		return nil, err
	}
	return s.toPublicKey(block)
}

// 从文件加载私钥
func (s *RsaHelper) LoadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	block, err := s.LoadFile(filename)
	if err != nil {
		return nil, err
	}

	return s.toPrivateKey(block)
}

// 从文件加载公钥
func (s *RsaHelper) LoadPublicKey(filename string) (*rsa.PublicKey, error) {
	block, err := s.LoadFile(filename)
	if err != nil {
		return nil, err
	}

	return s.toPublicKey(block)
}

// pem.Block转私钥
func (s *RsaHelper) toPrivateKey(block *pem.Block) (*rsa.PrivateKey, error) {
	var err error
	s.privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	s.publicKey = &s.privateKey.PublicKey
	return s.privateKey, nil
}

// pem.Block转公钥
func (s *RsaHelper) toPublicKey(block *pem.Block) (*rsa.PublicKey, error) {
	var err error
	s.publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return s.publicKey, nil
}

// base64 字节转pem.Block
func (s *RsaHelper) decodeKey(bytes []byte) (*pem.Block, error) {
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, errors.New("content is not a pem")
	}
	return block, nil
}

// 从文件加载pem.Block
func (s *RsaHelper) LoadFile(filename string) (*pem.Block, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return s.decodeKey(bytes)
}

// 保存公/私钥
func (s *RsaHelper) SaveKeyPair(dir, name string) error {
	if s.privateKey == nil {
		s.GenerateKey()
	}
	p := s.QueryKeyPair()

	tmpDir := strings.TrimRight(dir, "/") + "/"

	if err := s.saveFile(tmpDir+name, p.Private); err != nil {
		return err
	}
	if err := s.saveFile(tmpDir+name+".pub", p.Public); err != nil {
		return err
	}
	return nil
}

func (s *RsaHelper) saveFile(name string, bytes []byte) error {
	return ioutil.WriteFile(name, bytes, 0755)
}

// 公钥加密
func (s *RsaHelper) EncryptPKCS1v15(msg []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, s.publicKey, msg)
}

// 私钥解密
func (s *RsaHelper) DecryptPKCS1v15(msg []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, s.privateKey, msg)
}

// 公钥加密
func (s *RsaHelper) EncryptOAEP(msg, label []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, s.publicKey, msg, label)
}

// 私钥解密
func (s *RsaHelper) DecryptOAEP(msg, label []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, s.privateKey, msg, label)
}

// sha256签名
func (s *RsaHelper) SignPKCS1v15(msg []byte) ([]byte, error) {
	hashed := sha256.Sum256(msg)
	return rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hashed[:])
}

// sha256签名验证
func (s *RsaHelper) VerifyPKCS1v15(msg, sign []byte) error {
	hashed := sha256.Sum256(msg)
	return rsa.VerifyPKCS1v15(s.publicKey, crypto.SHA256, hashed[:], sign)
}

// 私钥加密
func (s *RsaHelper) EncryptPKCS1v15WithPrivatekey(msg []byte) ([]byte, error) {
	k := (s.privateKey.N.BitLen() + 7) / 8
	tLen := len(msg)
	if tLen > k-11 {
		return nil, errors.New("input size too large")
	}

	em := make([]byte, k)
	em[1] = 1
	for i := 2; i < k-tLen-1; i++ {
		em[i] = 0xff
	}

	copy(em[k-tLen:k], msg)
	c := new(big.Int).SetBytes(em)
	if c.Cmp(s.privateKey.N) > 0 {
		return nil, errors.New("encryption error")
	}

	var m *big.Int
	var ir *big.Int
	if s.privateKey.Precomputed.Dp == nil {
		m = new(big.Int).Exp(c, s.privateKey.D, s.privateKey.N)
	} else {
		m = new(big.Int).Exp(c, s.privateKey.Precomputed.Dp, s.privateKey.Primes[0])
		m2 := new(big.Int).Exp(c, s.privateKey.Precomputed.Dq, s.privateKey.Primes[1])
		m.Sub(m, m2)
		if m.Sign() < 0 {
			m.Add(m, s.privateKey.Primes[0])
		}

		m.Mul(m, s.privateKey.Precomputed.Qinv)
		m.Mod(m, s.privateKey.Primes[0])
		m.Mul(m, s.privateKey.Primes[1])
		m.Add(m, m2)

		for i, values := range s.privateKey.Precomputed.CRTValues {
			prime := s.privateKey.Primes[2+i]
			m2.Exp(c, values.Exp, prime)
			m2.Sub(m2, m)
			m2.Mul(m2, values.Coeff)
			m2.Mod(m2, prime)
			if m2.Sign() < 0 {
				m2.Add(m2, prime)
			}
			m2.Mul(m2, values.R)
			m.Add(m, m2)
		}
	}

	if ir != nil {
		m.Mul(m, ir)
		m.Mod(m, s.privateKey.N)
	}
	return m.Bytes(), nil
}

// 公钥解密
func (s *RsaHelper) DecryptPKCS1v15WithPublicKey(msg []byte) ([]byte, error) {
	c := new(big.Int)
	m := new(big.Int)
	m.SetBytes(msg)

	e := big.NewInt(int64(s.publicKey.E))
	out := c.Exp(m, e, s.publicKey.N).Bytes()
	skip := 0
	for i := 2; i < len(out); i++ {
		if i+1 >= len(out) {
			break
		}
		if out[i] == 0xff && out[i+1] == 0 {
			skip = i + 2
			break
		}
	}

	return out[skip:], nil
}

type (
	KeyPair struct {
		Private []byte
		Public  []byte
	}
)
