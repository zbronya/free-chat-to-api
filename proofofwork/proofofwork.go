package proofofwork

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/zbronya/free-chat-to-api/model"
	"golang.org/x/crypto/sha3"
	"math/rand"
	"time"
	_ "time/tzdata"
)

var (
	cores   = []int{8, 12, 16, 24, 32}
	screens = []int{3000, 4000, 6000}
	script  = "https://cdn.oaistatic.com/_next/static/chunks/main-c5c262a33e3f13d2.js?dpl=baf36960d05dde6d8b941194fa4093fb5cb78c6a"

	dpl = "baf36960d05dde6d8b941194fa4093fb5cb78c6a"

	errorPrefix = "gAAAAABwQ8Lk5FbGpA2NcR9dShT6gYjU7VxZ4D"
)

func getParseTime() string {
	loc, _ := time.LoadLocation("America/Los_Angeles")
	now := time.Now().In(loc)
	return now.Format("Mon Jan 02 2006 15:04:05") + " GMT-0800 (Pacific Time)"
}

func GetConfig(ua string) []interface{} {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	core := cores[rand.Intn(4)]
	rand.New(rand.NewSource(time.Now().UnixNano()))
	screen := screens[rand.Intn(3)]
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return []interface{}{core + screen, getParseTime(), int64(4294705152), 0, ua, script, dpl, "en-US", "en-US,en"}
}

func GetChatRequirementReq(config []interface{}) model.ChatRequirementReq {
	j, _ := json.Marshal(config)
	result := base64.StdEncoding.EncodeToString(j)
	return model.ChatRequirementReq{
		P: "gAAAAAC" + result,
	}
}

func CalcProofToken(config []interface{}, seed string, diff string) string {
	diffLen := len(diff) / 2
	hasher := sha3.New512()
	for i := 0; i < 1000000; i++ {
		config[3] = i
		j, _ := json.Marshal(config)
		base := base64.StdEncoding.EncodeToString(j)
		hasher.Write([]byte(seed + base))
		hash := hasher.Sum(nil)
		hasher.Reset()
		if hex.EncodeToString(hash[:diffLen]) <= diff {
			return "gAAAAAB" + base
		}
	}
	return errorPrefix + base64.StdEncoding.EncodeToString([]byte(`"`+seed+`"`))
}
