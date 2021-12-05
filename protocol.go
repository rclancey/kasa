package kasa

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	//"log"
	"net"
	"time"
)

const (
	INITIALIZATION_VECTOR = 171
	DEFAULT_PORT = 9999
	DEFAULT_TIMEOUT = 5
	BLOCK_SIZE = 4
)

func query(host string, req interface{}, dst interface{}) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}
	//log.Println("DEBUG:", string(payload))
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, DEFAULT_PORT), time.Duration(DEFAULT_TIMEOUT) * time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.SetWriteDeadline(time.Now().Add(time.Duration(DEFAULT_TIMEOUT) * time.Second))
	err = binary.Write(conn, binary.BigEndian, int32(len(payload)))
	if err != nil {
		return err
	}
	_, err = conn.Write(encrypt(payload))
	if err != nil {
		return err
	}
	conn.SetReadDeadline(time.Now().Add(time.Duration(DEFAULT_TIMEOUT) * time.Second))
	var respLen int32
	err = binary.Read(conn, binary.BigEndian, &respLen)
	if err != nil {
		return err
	}
	cipher := make([]byte, int(respLen))
	_, err = io.ReadFull(conn, cipher)
	if err != nil {
		return err
	}
	plain := decrypt(cipher)
	return json.Unmarshal(plain, dst)
}

func encrypt(plain []byte) []byte {
	key := byte(INITIALIZATION_VECTOR)
	cipher := make([]byte, len(plain))
	for i, b := range plain {
		key = key ^ b
		cipher[i] = key
	}
	return cipher
}

func decrypt(cipher []byte) []byte {
	key := byte(INITIALIZATION_VECTOR)
	plain := make([]byte, len(cipher))
	for i, b := range cipher {
		plain[i] = key ^ b
		key = b
	}
	return plain
}
