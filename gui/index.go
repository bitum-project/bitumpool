package gui

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/bitum-project/bitumd/bitumutil"
	"github.com/bitum-project/bitumpool/dividend"
	"github.com/bitum-project/bitumpool/network"
)

type indexData struct {
	PoolStats    *network.PoolStats
	PoolHashRate *big.Rat
	SoloPool     bool
	AccountStats *AccountStats
	Address      string
	Admin        bool
	Error        string
}

// AccountStats is a snapshot of an accounts contribution to the pool. This
// comprises of blocks mined by the pool and payments made to the account.
type AccountStats struct {
	MinedWork []*network.AcceptedWork
	Payments  []*dividend.Payment
	Clients   []*network.ClientInfo
}

func (ui *GUI) GetIndex(w http.ResponseWriter, r *http.Request) {
	poolStats, err := ui.hub.FetchPoolStats()
	if err != nil {
		log.Error(err)
		http.Error(w, "FetchPoolStats error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	clientInfo := ui.hub.FetchClientInfo()

	poolHashRate := new(big.Rat).SetInt64(0)
	for _, clients := range clientInfo {
		for _, client := range clients {
			poolHashRate = poolHashRate.Add(poolHashRate, client.HashRate)
		}
	}

	data := indexData{
		PoolStats:    poolStats,
		PoolHashRate: poolHashRate,
		SoloPool:     ui.cfg.SoloPool,
		Admin:        false,
	}

	address := r.FormValue("address")

	if address == "" {
		ui.renderTemplate(w, r, "index", data)
		return
	}

	// Add address to template data so the input field can be repopulated
	// with the users input.
	data.Address = address

	// Ensure address is correct length
	if len(address) != 35 {
		data.Error = fmt.Sprintf("Address should be 35 characters")
		ui.renderTemplate(w, r, "index", data)
		return
	}

	// Ensure address can be decoded
	addr, err := bitumutil.DecodeAddress(address)
	if err != nil {
		data.Error = fmt.Sprintf("Failed to decode address")
		ui.renderTemplate(w, r, "index", data)
		return
	}

	// Ensure address is on the active network.
	if !addr.IsForNet(ui.cfg.ActiveNet) {
		data.Error = fmt.Sprintf("This is not a %s address", ui.cfg.ActiveNet.Name)
		ui.renderTemplate(w, r, "index", data)
		return
	}

	accountID := *dividend.AccountID(address)

	work, err := ui.hub.FetchMinedWorkByAddress(accountID)
	if err != nil {
		log.Error(err)
		http.Error(w, "FetchMinedWorkByAddress error: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	payments, err := ui.hub.FetchPaymentsForAddress(accountID)
	if err != nil {
		log.Error(err)
		http.Error(w, "FetchPaymentsForAddress error: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	if len(work) == 0 && len(payments) == 0 {
		log.Infof("Nothing found for address: %s", address)
		data.Error = fmt.Sprintf("Nothing found for address: %s", address)
		ui.renderTemplate(w, r, "index", data)
		return
	}

	data.AccountStats = &AccountStats{
		MinedWork: work,
		Payments:  payments,
		Clients:   clientInfo[accountID],
	}

	ui.renderTemplate(w, r, "index", data)
}
