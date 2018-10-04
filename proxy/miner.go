package proxy

import (
	"log"
	"math/big"
	"strconv"
	"strings"

	"gitlab.com/aquachain/aquachain/common"
	"gitlab.com/aquachain/aquachain/consensus/lightvalid"
)

var hasher = lightvalid.New()

func (s *ProxyServer) processShare(login, id, ip string, t *BlockTemplate, params []string) (bool, bool) { // exist , validshare
	nonceHex := params[0]
	hashNoNonce := params[1]
	mixDigest := params[2]
	nonce, _ := strconv.ParseUint(strings.Replace(nonceHex, "0x", "", -1), 16, 64)
	shareDiff := s.config.Proxy.Difficulty
	if t == nil {
		log.Printf("nil block template from %v@%v", login, ip)
		return false, false
	}
	h, ok := t.headers[hashNoNonce]
	if !ok {
		log.Printf("Stale share from %v@%v", login, ip)
		return false, false
	}
	if t.version == 0 {
		log.Printf("No version set %v@%v", login, ip)
		return false, false
	}

	// extra check for duplicate shares without db
	s.sharelock.Lock()
	latestwork, exists := s.shares[nonce]
	if exists {
		if latestwork == hashNoNonce {
			s.sharelock.Unlock()
			log.Printf("evil share from %v@%v", login, ip)
			return true, false
		}
	}
	s.shares[nonce] = hashNoNonce
	s.sharelock.Unlock()

	share := Block{
		number:      h.height,
		hashNoNonce: common.HexToHash(hashNoNonce),
		difficulty:  big.NewInt(shareDiff),
		nonce:       nonce,
		mixDigest:   common.HexToHash(mixDigest),
		version:     t.version,
	}

	block := Block{
		number:      h.height,
		hashNoNonce: common.HexToHash(hashNoNonce),
		difficulty:  h.diff,
		nonce:       nonce,
		mixDigest:   common.HexToHash(mixDigest),
		version:     t.version,
	}

	// test at pool difficulty
	if !hasher.Verify(share) {
		log.Println("not ok share, invalid POW")
		return false, false
	}

	// at real difficulty
	if hasher.Verify(block) {
		ok, err := s.rpc().SubmitWork(params)
		if err != nil {
			log.Printf("Block submission failure at height %v for %v: %v", h.height, t.Header, err)
		} else if !ok {
			log.Printf("Block rejected at height %v for %v", h.height, t.Header)
			return false, false
		} else {
			s.fetchBlockTemplate()
			exist, err := s.backend.WriteBlock(login, id, params, shareDiff, h.diff.Int64(), h.height, s.hashrateExpiration)
			if exist {
				return true, false
			}
			if err != nil {
				log.Println("Failed to insert block candidate into backend:", err)
			} else {
				log.Printf("Inserted block %v to backend", h.height)
			}
			log.Printf("Block found by miner %v@%v at height %d", login, ip, h.height)
			s.sharelock.Lock()
			for i := range s.shares {
				delete(s.shares, i)
			}
			s.sharelock.Unlock()
		}
	} else {
		exist, err := s.backend.WriteShare(login, id, params, shareDiff, h.height, s.hashrateExpiration)
		if exist {
			return true, false
		}
		if err != nil {
			log.Println("Failed to insert share data into backend:", err)
		}
	}
	return false, true
}
