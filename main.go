package main

import (
	"CheckFirewall/lib"
	b64 "encoding/base64"
	"fmt"
	"github.com/vmihailenco/msgpack"
	"github.com/zerospam/rsa/rsa"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/check", func(w http.ResponseWriter, req *http.Request) {
		var transportServer lib.TransportServer

		requestBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}

		sDecode, err := b64.StdEncoding.DecodeString(string(requestBody))

		if err != nil {
			panic(err)
		}

		decrypted, err := rsa.PublicDecrypt(sDecode, "./public.pem", rsa.RSA_PKCS1_PADDING)
		if err != nil {
			panic(err)
		}

		err = msgpack.Unmarshal(decrypted, &transportServer)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "%s %d", transportServer.Transport, transportServer.Port)
	})

	http.HandleFunc("/encode", func(w http.ResponseWriter, r *http.Request) {
		transport := lib.TransportServer{
			Port:      25,
			Transport: "mail.example.com",
		}

		msg, err := msgpack.Marshal(transport)
		if err != nil {
			panic(err)
		}
		encrypted, err := rsa.PrivateEncrypt([]byte(msg), "./private.pem", rsa.RSA_PKCS1_PADDING)
		if err != nil {
			panic(err)
		}
		sEnc := b64.StdEncoding.EncodeToString(encrypted)
		fmt.Fprint(w, sEnc)
	})

	http.ListenAndServe(":8082", nil)
}
