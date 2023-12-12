package auth

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
)

type Authenticator struct {
	heartbeatUrl string

	name             string
	maxPlayers, port int
	public           bool

	salt string
}

func genSalt() string {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	salt := make([]byte, 16)

	for i := 0; i < 16; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
		if err != nil {
			return ""
		}
		salt[i] = characters[randomIndex.Int64()]
	}

	return string(salt)
}

func NewAuthenticator(url string, name string, maxPlayers, port int, public bool) *Authenticator {
	return &Authenticator{heartbeatUrl: url, name: name, maxPlayers: maxPlayers, port: port, public: public, salt: genSalt()}
}

func (auth *Authenticator) Heartbeat(playerCount int) (string, error) {
	pb := "False"
	if auth.public {
		pb = "True"
	}
	u, err := url.ParseRequestURI(fmt.Sprintf("%s?port=%d&max=%d&name=%s&public=%s&version=7&salt=%s&users=%d&software=Obsidian", auth.heartbeatUrl, auth.port, auth.maxPlayers, url.QueryEscape(auth.name), pb, auth.salt, playerCount))
	if err != nil {
		return "", err
	}

	r, _ := http.Get(u.String())

	d, _ := io.ReadAll(r.Body)

	if _, err := url.ParseRequestURI(string(d)); err == nil {
		return string(d), nil
	}
	return "", fmt.Errorf(string(d))
}

func (auth *Authenticator) Validate(key string, username string) (ok bool) {
	d := md5.Sum([]byte(auth.salt + username))
	calculatedKey := hex.EncodeToString(d[:])
	return key == calculatedKey
}
