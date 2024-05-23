package proofofwork

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/zbronya/free-chat-to-api/model"
	"golang.org/x/crypto/sha3"
	"math/rand"
	"time"
	_ "time/tzdata"
)

var (
	cores   = []int{1, 2, 4, 8}
	screens = []int{3000, 4000, 6000}
	script  = "https://cdn.oaistatic.com/_next/static/chunks/2565-9cf19ba0b7d24a5d.js?dpl=4811fd1c94b550c8f03fcc863ee6c1a99940efc5"

	dpl = "4811fd1c94b550c8f03fcc863ee6c1a99940efc5"

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
	return []interface{}{core + screen,
		getParseTime(),
		int64(4294705152),
		0,
		ua,
		script,
		dpl,
		"en-US",
		"en-US,en",
		0,
		"updateAdInterestGroups−function updateAdInterestGroups() { [native code] }",
		"location",
		"__NEXT_PRELOADREADY",
		885.6999999880791,
	}
}

func GetChatRequirementReq(config []interface{}) model.ChatRequirementReq {
	randomFloat := rand.Float64()
	seed := fmt.Sprintf("%.6f", randomFloat)
	token := CalcProofToken(config, seed, "0")
	return model.ChatRequirementReq{
		P: token,
	}

}

func CalcProofToken(config []interface{}, seed string, diff string) string {

	diffLen := len(diff) / 2
	hasher := sha3.New512()
	startTime := time.Now()
	for i := 0; i < 1000000; i++ {
		config[3] = i
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		config[9] = elapsed.Milliseconds()
		j, _ := json.Marshal(config)
		base := base64.StdEncoding.EncodeToString(j)
		hasher.Write([]byte(seed + base))
		hash := hasher.Sum(nil)
		hasher.Reset()
		if hex.EncodeToString(hash[:diffLen]) <= diff {
			return base
		}
	}
	return errorPrefix + base64.StdEncoding.EncodeToString([]byte(`"`+seed+`"`))
}
