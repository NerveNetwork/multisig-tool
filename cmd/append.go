package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/niels1286/multisig-tool/cfg"
	"github.com/niels1286/multisig-tool/i18n"
	"github.com/niels1286/multisig-tool/utils"
	txprotocal "github.com/niels1286/nerve-go-sdk/protocal"
	"github.com/niels1286/nerve-go-sdk/protocal/txdata"
	"github.com/spf13/cobra"
	"math/big"
	"strings"
)

// appendCmd represents the deposit command
var appendCmd = &cobra.Command{
	Use:   "append",
	Short: i18n.GetText("0002"),
	Long:  i18n.GetText("0002"),
	Run: func(cmd *cobra.Command, args []string) {
		if "" == nodeHash || strings.TrimSpace(nodeHash) == "" {
			fmt.Println(i18n.GetText("0017"))
			return
		}
		sdk := utils.GetOfficalSdk()
		msAccount, err := sdk.CreateMultiAccount(m, pks)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		cId := cfg.MainChainId
		aId := cfg.MainAssetsId
		tx := utils.AssembleTransferTx(m, pks, cId, aId, amount, "", msAccount.Address, 0, cfg.POCLockValue, nil, true)
		if tx == nil {
			fmt.Println(i18n.GetText("10001"))
			return
		}
		tx.TxType = txprotocal.APPEND_AGENT_DEPOSIT
		value := big.NewFloat(amount)
		value = value.Mul(value, big.NewFloat(100000000))
		x, _ := value.Uint64()

		depositData := txdata.ChangeNodeDeposit{
			Amount:   big.NewInt(int64(x)),
			Address:  msAccount.AddressBytes,
			NodeHash: txprotocal.ImportNulsHash(nodeHash),
		}
		tx.Extend, err = depositData.Serialize()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		txBytes, err := tx.Serialize()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		txHex := hex.EncodeToString(txBytes)

		fmt.Println("Successed:\ntxHex : " + txHex)
		fmt.Println("txHash : " + tx.GetHash().String())
	},
}

func init() {
	rootCmd.AddCommand(appendCmd)
	appendCmd.Flags().IntVarP(&m, "m", "m", 0, i18n.GetText("0014"))
	appendCmd.MarkFlagRequired("m")
	appendCmd.Flags().StringVarP(&pks, "publickeys", "p", "", i18n.GetText("0015"))
	appendCmd.MarkFlagRequired("publickeys")
	appendCmd.Flags().Float64VarP(&amount, "amount", "a", 0, i18n.GetText("0019"))
	appendCmd.MarkFlagRequired("amount")
	appendCmd.Flags().StringVarP(&nodeHash, "nodeHash", "n", "", i18n.GetText("0018"))
	appendCmd.MarkFlagRequired("nodeHash")
}
