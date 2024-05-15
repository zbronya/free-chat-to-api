package proofofwork

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/allegro/bigcache"
	"github.com/zbronya/free-chat-to-api/logger"
	"github.com/zbronya/free-chat-to-api/model"
	"golang.org/x/crypto/sha3"
	"math/rand"
	"strings"
	"time"
	_ "time/tzdata"
)

var (
	cores   = []int{1, 2, 4, 8}
	screens = []int{3000, 4000, 6000}
	script  = "https://cdn.oaistatic.com/_next/static/chunks/main-c5c262a33e3f13d2.js?dpl=a44a6d28cfe80fc54efc0ce87573ae13d9b8e9bd"

	dpl = "a44a6d28cfe80fc54efc0ce87573ae13d9b8e9bd"

	errorPrefix = "wQ8Lk5FbGpA2NcR9dShT6gYjU7VxZ4D"
)

var cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(1 * time.Hour))

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
	return []interface{}{core + screen, getParseTime(), int64(4294705152), 0, ua, script, dpl, "en-US", "en-US,en", 0, "webdriver−false", "location", "onmouseenter"}
}

func GetChatRequirementReq(config []interface{}) model.ChatRequirementReq {

	result, err := cache.Get("config")

	return model.ChatRequirementReq{
		P: "gAAAAACWzIzNDgsIldlZCBNYXkgMTUgMjAyNCAwOToyMzoxNiBHTVQrMDgwMCAo5Lit5Zu95qCH5YeG5pe26Ze0KSIsNDI5NDcwNTE1MiwyMywiTW96aWxsYS81LjAgKE1hY2ludG9zaDsgSW50ZWwgTWFjIE9TIFggMTBfMTVfNykgQXBwbGVXZWJLaXQvNTM3LjM2IChLSFRNTCwgbGlrZSBHZWNrbykgQ2hyb21lLzEyNC4wLjAuMCBTYWZhcmkvNTM3LjM2IiwiaHR0cHM6Ly9jZG4ub2Fpc3RhdGljLmNvbS9fbmV4dC9zdGF0aWMvY2h1bmtzL21haW4tMGI1NjAxZWMwOWVlYzc4Yi5qcz9kcGw9N2M3NDBjOTgxY2JkMWYyODZkYjhiY2JkMDNjYTBmMDI5OWQ4MTBmYiIsImRwbD03Yzc0MGM5ODFjYmQxZjI4NmRiOGJjYmQwM2NhMGYwMjk5ZDgxMGZiIiwiemgtQ04iLCJ6aC1DTix6aCxlbiIsMjMyLCJ3ZWJraXRQZXJzaXN0ZW50U3RvcmFnZeKIkltvYmplY3QgRGVwcmVjYXRlZFN0b3JhZ2VRdW90YV0iLCJsb2NhdGlvbiIsIm91dGVyV2lkdGgiXQ==",
	}

	if err != nil {
		randomFloat := rand.Float64()
		seed := fmt.Sprintf("%.6f", randomFloat)
		tmp := CalcProofToken(config, seed, "000000")

		if strings.HasPrefix(tmp, errorPrefix) {
			logger.GetLogger().Warn("CalcProofToken config: " + tmp)
			return model.ChatRequirementReq{
				P: tmp,
			}
		}

		cache.Set("config", []byte(tmp))
		return model.ChatRequirementReq{
			P: "gAAAAAC" + tmp,
		}
	} else {
		return model.ChatRequirementReq{
			P: "gAAAAAC" + string(result),
		}
	}
}

func CalcProofToken(config []interface{}, seed string, diff string) string {
	diffLen := len(diff) / 2
	hasher := sha3.New512()
	for i := 0; i < 1000000; i++ {
		config[3] = i
		config[9] = (i + 2) / 2
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
