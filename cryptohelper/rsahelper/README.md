# RAS 加/解密, 签名/验证

```golang
package main

import (
	"app/rashelper"
	"encoding/hex"
	"fmt"
)

// 保存公/私钥匙到文件
func genKey() {
	helper := rashelper.New()
	err := helper.SaveKeyPair(".", "key")
	fmt.Println(err)
}

// 公钥加密/私钥解密
func stand() {
	helper := rashelper.New()
	_, err := helper.LoadPrivateKey("./key")
	if err != nil {
		fmt.Println(err)
		return
	}

	str := "abcd1234"
	msg, err := helper.EncryptPKCS1v15([]byte(str))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("密文:", hex.EncodeToString(msg))

	raw, err := helper.DecryptPKCS1v15(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("明文:", string(raw))
}

// 扩展的私钥加密/公钥解密，标准库中没的这两人个方法
func custom() {
	helper := rashelper.New()
	_, err := helper.LoadPrivateKey("./key")
	if err != nil {
		fmt.Println(err)
		return
	}

	str := "abcd1234"
	msg, err := helper.EncryptPKCS1v15WithPrivatekey([]byte(str))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("密文:", hex.EncodeToString(msg))

	raw, err := helper.DecryptPKCS1v15WithPublicKey(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("明文:", string(raw))
}

func main() {
	genKey()
	stand()
	custom()
}

```