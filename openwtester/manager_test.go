package openwtester

import (
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openw"
	"github.com/blocktree/openwallet/openwallet"
	"path/filepath"
	"testing"
)

var (
	testApp        = "fiiicoin-adapter"
	configFilePath = filepath.Join("conf")
)

func testInitWalletManager() *openw.WalletManager {
	log.SetLogFuncCall(true)
	tc := openw.NewConfig()

	tc.ConfigDir = configFilePath
	tc.EnableBlockScan = false
	tc.SupportAssets = []string{
		"FIII",
	}
	return openw.NewWalletManager(tc)
	//tm.Init()
}

func TestWalletManager_CreateWallet(t *testing.T) {
	tm := testInitWalletManager()
	w := &openwallet.Wallet{Alias: "HELLO FIII", IsTrust: true, Password: "12345678"}
	nw, key, err := tm.CreateWallet(testApp, w)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("wallet:", nw)
	log.Info("key:", key)

}

func TestWalletManager_GetWalletInfo(t *testing.T) {

	tm := testInitWalletManager()

	wallet, err := tm.GetWalletInfo(testApp, "W7tue6SDce38fPwerdKqyebUh6yo2nTQLC")
	if err != nil {
		log.Error("unexpected error:", err)
		return
	}
	log.Info("wallet:", wallet)
}

func TestWalletManager_GetWalletList(t *testing.T) {

	tm := testInitWalletManager()

	list, err := tm.GetWalletList(testApp, 0, 10000000)
	if err != nil {
		log.Error("unexpected error:", err)
		return
	}
	for i, w := range list {
		log.Info("wallet[", i, "] :", w)
	}
	log.Info("wallet count:", len(list))

	tm.CloseDB(testApp)
}

func TestWalletManager_CreateAssetsAccount(t *testing.T) {

	tm := testInitWalletManager()

	walletID := "WKFkmvsSFz5mC1cAX3edJC2e6hH6ow3X9E"
	account := &openwallet.AssetsAccount{Alias: "mainnetFIII", WalletID: walletID, Required: 1, Symbol: "FIII", IsTrust: true}
	account, address, err := tm.CreateAssetsAccount(testApp, walletID, "12345678", account, nil)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("account:", account)
	log.Info("address:", address)

	tm.CloseDB(testApp)
}

func TestWalletManager_GetAssetsAccountList(t *testing.T) {

	tm := testInitWalletManager()

	walletID := "WKFkmvsSFz5mC1cAX3edJC2e6hH6ow3X9E"
	list, err := tm.GetAssetsAccountList(testApp, walletID, 0, 10000000)
	if err != nil {
		log.Error("unexpected error:", err)
		return
	}
	for i, w := range list {
		log.Infof("account[%d] : %+v", i, w)
	}
	log.Info("account count:", len(list))

	tm.CloseDB(testApp)

}

func TestWalletManager_CreateAddress(t *testing.T) {

	tm := testInitWalletManager()

	walletID := "WKFkmvsSFz5mC1cAX3edJC2e6hH6ow3X9E"
	accountID := "4h4wnCmpzgy3ZTeoMHs3gjDCuWyXQcxDsk9dcwbNGhmR"
	address, err := tm.CreateAddress(testApp, walletID, accountID, 300)
	if err != nil {
		log.Error(err)
		return
	}

	for _, addr := range address {
		log.Info(addr.Address)
	}


	tm.CloseDB(testApp)
}

func TestWalletManager_GetAddressList(t *testing.T) {

	tm := testInitWalletManager()

	walletID := "WKFkmvsSFz5mC1cAX3edJC2e6hH6ow3X9E"
	//accountID := "HX4tUVg5eETb6SvZeGeAFwk4PQ1CWS6dQeyjj3CqfYyK"
	accountID := "4h4wnCmpzgy3ZTeoMHs3gjDCuWyXQcxDsk9dcwbNGhmR"
	list, err := tm.GetAddressList(testApp, walletID, accountID, 0, -1, false)
	if err != nil {
		log.Error("unexpected error:", err)
		return
	}
	for i, w := range list {
		log.Info("address[", i, "] :", w.Address)
		//log.Info("address[", i, "] :", w.PublicKey)
	}
	log.Info("address count:", len(list))

	tm.CloseDB(testApp)
}


