package utils

import (
	"crypto/md5"
	"encoding/hex"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GenerateUUID(id uint) (uuid.UUID, error) {
	h := md5.New()
	io.WriteString(h, string(id)+time.Now().String())
	sum := h.Sum(nil)
	return uuid.FromBytes(sum)
}

func GenerateHash() string {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(strconv.FormatInt(time.Now().Unix(), 10)),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Fatal(err)
	}
	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetIPAddress(r *http.Request) string {
	var ipAddress string
	ipAddress = strings.Split(r.RemoteAddr, ":")[0]
	for _, h := range []string{"X-Forwarded-For", "X-Real-IP"} {
		for _, ip := range strings.Split(r.Header.Get(h), ",") {
			// header can contain spaces too, strip those out.
			ip = strings.TrimSpace(ip)
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() {
				// bad address, go to next
				continue
			} else {
				ipAddress = ip
				goto Done
			}
		}
	}
Done:
	return ipAddress
}
