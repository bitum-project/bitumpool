// Copyright (c) 2018 The Bitum developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package dividend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	bolt "github.com/coreos/bbolt"
	"github.com/bitum-project/bitumd/chaincfg"
	"github.com/bitum-project/bitumd/bitumutil"
	"github.com/bitum-project/bitumpool/database"
	"github.com/bitum-project/bitumpool/util"
)

// Miner types lists all known BITUM miners
const (
	CPU           = "cpu"
	InnosiliconD9 = "innosilicond9"
	// ObeliskBITUM1   = "obeliskbitum1"
	AntminerDR3  = "antminerdr3"
	AntminerDR5  = "antminerdr5"
	WhatsminerD1 = "whatsminerd1"
)

// MinerHashes is a map of all known BITUM miners and their coressponding
// hashrates.
var MinerHashes = map[string]*big.Int{
	CPU:           new(big.Int).SetInt64(70E3),
	InnosiliconD9: new(big.Int).SetInt64(2.4E12),
	// ObeliskBITUM1:   new(big.Int).SetInt64(1.1E12),
	AntminerDR3:  new(big.Int).SetInt64(7.8E12),
	AntminerDR5:  new(big.Int).SetInt64(35E12),
	WhatsminerD1: new(big.Int).SetInt64(48E12),
}

// MinerPorts is a map of all known BITUM miners and the coressponding
// ports configured to be connected on.
var MinerPorts = map[string]uint32{
	CPU: 5550,
	// ObeliskBITUM1:   5551,
	InnosiliconD9: 5552,
	AntminerDR3:   5553,
	AntminerDR5:   5554,
	WhatsminerD1:  5555,
}

// Convenience variables.
var (
	zeroRat = new(big.Rat).SetInt64(0)
	zeroInt = new(big.Int).SetInt64(0)
)

var (
	// PoolFeesK is the key used to track pool fee payouts.
	PoolFeesK = "fees"

	// PPS represents the pay per share payment method.
	PPS = "pps"

	// PPLNS represents the pay per last n shares payment method.
	PPLNS = "pplns"
)

// ShareWeights reprsents the associated weights for each known BITUM miner.
// With the share weight of the lowest hash BITUM miner (LHM) being 1, the
// rest were calculated as :
// 				(Hash of Miner X * Weight of LHM)/ Hash of LHM
var ShareWeights = map[string]*big.Rat{
	CPU: new(big.Rat).SetFloat64(1.0), // Reserved for testing.
	// ObeliskBITUM1:   new(big.Rat).SetFloat64(1.0),
	InnosiliconD9: new(big.Rat).SetFloat64(2.182),
	AntminerDR3:   new(big.Rat).SetFloat64(7.091),
	AntminerDR5:   new(big.Rat).SetFloat64(31.181),
	WhatsminerD1:  new(big.Rat).SetFloat64(43.636),
}

// CalculatePoolDifficulty determines the difficulty at which the provided
// hashrate can generate a pool share by the provided target time.
func CalculatePoolDifficulty(net *chaincfg.Params, hashRate *big.Int, targetTimeSecs *big.Int) (*big.Int, error) {
	hashesPerTargetTime := new(big.Int).Mul(hashRate, targetTimeSecs)
	powLimit := net.PowLimit
	powLimitFloat, _ := new(big.Float).SetInt(powLimit).Float64()

	// The number of possible iterations is calculated as:
	//
	//    iterations := 2^(256 - floor(log2(pow_limit)))
	iterations := math.Pow(2, 256-math.Floor(math.Log2(powLimitFloat)))

	// The difficulty at which the provided hashrate can mine a block is
	// calculated as:
	//
	//    difficulty = (hashes_per_sec * target_in_seconds) / iterations
	difficulty := new(big.Rat).Quo(new(big.Rat).SetInt(hashesPerTargetTime),
		new(big.Rat).SetFloat64(iterations))
	diff := new(big.Int).Quo(difficulty.Num(), difficulty.Denom())

	// Clamp the difficulty to 1 if needed.
	if diff.Cmp(zeroInt) == 0 {
		diff = new(big.Int).SetInt64(1)
	}

	return diff, nil
}

// DifficultyToTarget converts the provided difficulty to a target based on the
// active network.
func DifficultyToTarget(net *chaincfg.Params, difficulty *big.Int) (*big.Int, error) {
	powLimit := net.PowLimit

	// The corresponding target is calculated as:
	//
	//    target = pow_limit / difficulty
	//
	// The result is clamped to the pow limit if it exceeds it.
	target := new(big.Int).Div(powLimit, difficulty)

	if target.Cmp(powLimit) > 0 {
		target = powLimit
	}

	return target, nil
}

// CalculatePoolTarget determines the target difficulty at which the provided
// hashrate can generate a pool share by the provided target time.
func CalculatePoolTarget(net *chaincfg.Params, hashRate *big.Int, targetTimeSecs *big.Int) (*big.Int, *big.Int, error) {
	difficulty, err := CalculatePoolDifficulty(net, hashRate, targetTimeSecs)
	if err != nil {
		return nil, nil, err
	}

	target, err := DifficultyToTarget(net, difficulty)
	return target, difficulty, err
}

// Share represents verifiable work performed by a pool client.
type Share struct {
	Account   string   `json:"account"`
	Weight    *big.Rat `json:"weight"`
	CreatedOn int64    `json:"createdOn"`
}

// NewShare creates a shate with the provided account and weight.
func NewShare(account string, weight *big.Rat) *Share {
	return &Share{
		Account:   account,
		Weight:    weight,
		CreatedOn: time.Now().UnixNano(),
	}
}

// ErrNotSupported is returned when an entity does not support an action.
func ErrNotSupported(tp, action string) error {
	return fmt.Errorf("action (%v) not supported for type (%v)",
		action, tp)
}

// Create persists a share to the database.
func (s *Share) Create(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		pbkt := tx.Bucket(database.PoolBkt)
		if pbkt == nil {
			return database.ErrBucketNotFound(database.PoolBkt)
		}
		bkt := pbkt.Bucket(database.ShareBkt)
		if bkt == nil {
			return database.ErrBucketNotFound(database.ShareBkt)
		}
		sBytes, err := json.Marshal(s)
		if err != nil {
			return err
		}
		err = bkt.Put(util.NanoToBigEndianBytes(s.CreatedOn), sBytes)
		return err
	})
	return err
}

// Update is not supported for shares.
func (s *Share) Update(db *bolt.DB) error {
	return ErrNotSupported("share", "update")
}

// Delete is not supported for shares.
func (s *Share) Delete(db *bolt.DB) error {
	return ErrNotSupported("share", "delete")
}

// ErrDivideByZero is returned the divisor of division operation is zero.
func ErrDivideByZero() error {
	return fmt.Errorf("divide by zero: divisor is zero")
}

// PPSEligibleShares fetches all shares within the provided inclusive bounds.
func PPSEligibleShares(db *bolt.DB, min []byte, max []byte) ([]*Share, error) {
	eligibleShares := make([]*Share, 0)
	err := db.View(func(tx *bolt.Tx) error {
		pbkt := tx.Bucket(database.PoolBkt)
		if pbkt == nil {
			return database.ErrBucketNotFound(database.PoolBkt)
		}
		bkt := pbkt.Bucket(database.ShareBkt)
		if bkt == nil {
			return database.ErrBucketNotFound(database.ShareBkt)
		}

		c := bkt.Cursor()
		if min == nil {
			for k, v := c.First(); k != nil; k, v = c.Next() {
				var share Share
				err := json.Unmarshal(v, &share)
				if err != nil {
					return err
				}

				eligibleShares = append(eligibleShares, &share)
			}
		}

		if min != nil {
			for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
				var share Share
				err := json.Unmarshal(v, &share)
				if err != nil {
					return err
				}

				eligibleShares = append(eligibleShares, &share)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return eligibleShares, err
}

// PPLNSEligibleShares fetches all shares keyed greater than the provided
// minimum.
func PPLNSEligibleShares(db *bolt.DB, min []byte) ([]*Share, error) {
	eligibleShares := make([]*Share, 0)
	err := db.View(func(tx *bolt.Tx) error {
		pbkt := tx.Bucket(database.PoolBkt)
		if pbkt == nil {
			return database.ErrBucketNotFound(database.PoolBkt)
		}
		bkt := pbkt.Bucket(database.ShareBkt)
		if bkt == nil {
			return database.ErrBucketNotFound(database.ShareBkt)
		}

		c := bkt.Cursor()
		for k, v := c.Last(); k != nil && bytes.Compare(k, min) > 0; k, v = c.Prev() {
			var share Share
			err := json.Unmarshal(v, &share)
			if err != nil {
				return err
			}

			eligibleShares = append(eligibleShares, &share)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return eligibleShares, err
}

// CalculateSharePercentages calculates the percentages due each account according
// to their weighted shares.
func CalculateSharePercentages(shares []*Share) (map[string]*big.Rat, error) {
	totalShares := new(big.Rat)
	tally := make(map[string]*big.Rat)
	dividends := make(map[string]*big.Rat)

	// Tally all share weights for each participation account.
	for _, share := range shares {
		totalShares = totalShares.Add(totalShares, share.Weight)
		if _, ok := tally[share.Account]; ok {
			tally[share.Account] = tally[share.Account].
				Add(tally[share.Account], share.Weight)
			continue
		}

		tally[share.Account] = share.Weight
	}

	// Calculate each participating account to be claimed.
	for account, shareCount := range tally {
		if tally[account].Cmp(zeroRat) == 0 {
			return nil, ErrDivideByZero()
		}

		dividend := new(big.Rat).Quo(shareCount, totalShares)
		dividends[account] = dividend
	}

	return dividends, nil
}

// CalculatePayments calculates the payments due participating accounts.
func CalculatePayments(percentages map[string]*big.Rat, total bitumutil.Amount, poolFee float64, height uint32, estMaturity uint32) ([]*Payment, error) {
	// Deduct pool fee from the amount to be shared.
	fee := total.MulF64(poolFee)
	amtSansFees := total - fee

	// Calculate each participating account's portion of the amount after fees.
	payments := make([]*Payment, 0)
	for account, percentage := range percentages {
		percent, _ := percentage.Float64()
		amt := amtSansFees.MulF64(percent)
		payments = append(payments, NewPayment(account, amt, height, estMaturity))
	}

	// Add a payout entry for pool fees.
	payments = append(payments, NewPayment(PoolFeesK, fee, height, estMaturity))

	return payments, nil
}

// PruneShares removes invalidated shares from the db.
func PruneShares(db *bolt.DB, minNano int64) error {
	minBytes := util.NanoToBigEndianBytes(minNano)
	err := db.Update(func(tx *bolt.Tx) error {
		pbkt := tx.Bucket(database.PoolBkt)
		if pbkt == nil {
			return database.ErrBucketNotFound(database.PoolBkt)
		}
		bkt := pbkt.Bucket(database.ShareBkt)
		if bkt == nil {
			return database.ErrBucketNotFound(database.ShareBkt)
		}

		toDelete := [][]byte{}
		cursor := bkt.Cursor()
		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			if bytes.Compare(minBytes, k) > 0 {
				toDelete = append(toDelete, k)
			}
		}

		for _, entry := range toDelete {
			err := bkt.Delete(entry)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
