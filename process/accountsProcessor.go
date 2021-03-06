package process

import (
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/elrond-accounts-manager/core"
	"github.com/ElrondNetwork/elrond-accounts-manager/data"
	"github.com/tidwall/gjson"
)

const (
	pathNodeStatusMeta = "/network/status/4294967295"
)

type accountsProcessor struct {
	AccountsGetterHandler
	restClient RestClientHandler
}

// NewAccountsProcessor will create a new instance of accountsProcessor
func NewAccountsProcessor(restClient RestClientHandler, acctsGetter AccountsGetterHandler) (*accountsProcessor, error) {
	return &accountsProcessor{
		restClient:            restClient,
		AccountsGetterHandler: acctsGetter,
	}, nil
}

// GetAllAccountsWithStake will return all accounts with stake
func (ap *accountsProcessor) GetAllAccountsWithStake() (map[string]*data.AccountInfoWithStakeValues, []string, error) {
	legacyDelegators, err := ap.GetLegacyDelegatorsAccounts()
	if err != nil {
		return nil, nil, err
	}

	validators, err := ap.GetValidatorsAccounts()
	if err != nil {
		return nil, nil, err
	}

	delegators, err := ap.GetDelegatorsAccounts()
	if err != nil {
		return nil, nil, err
	}

	lkMexAccountsWithStake, err := ap.GetLKMEXStakeAccounts()
	if err != nil {
		return nil, nil, err
	}

	allAccounts, allAddresses := ap.mergeAccounts(legacyDelegators, validators, delegators, lkMexAccountsWithStake)

	calculateTotalStakeForAccounts(allAccounts)

	return allAccounts, allAddresses, nil
}

func calculateTotalStakeForAccounts(accounts map[string]*data.AccountInfoWithStakeValues) {
	for _, account := range accounts {
		totalStake, totalStakeNum := computeTotalBalance(
			account.DelegationLegacyWaiting,
			account.DelegationLegacyActive,
			account.ValidatorsActive,
			account.ValidatorTopUp,
			account.Delegation,
		)

		account.TotalStake = totalStake
		account.TotalStakeNum = totalStakeNum
	}
}

func (ap *accountsProcessor) mergeAccounts(
	legacyDelegators, validators, delegators, lkMexAccountsWithStake map[string]*data.AccountInfoWithStakeValues,
) (map[string]*data.AccountInfoWithStakeValues, []string) {
	allAddresses := make([]string, 0)
	mergedAccounts := make(map[string]*data.AccountInfoWithStakeValues)

	for address, legacyDelegator := range legacyDelegators {
		mergedAccounts[address] = legacyDelegator

		allAddresses = append(allAddresses, address)
	}

	for address, stakedValidators := range validators {
		_, ok := mergedAccounts[address]
		if !ok {
			mergedAccounts[address] = stakedValidators

			allAddresses = append(allAddresses, address)
			continue
		}

		mergedAccounts[address].ValidatorsActive = stakedValidators.ValidatorsActive
		mergedAccounts[address].ValidatorsActiveNum = stakedValidators.ValidatorsActiveNum
		mergedAccounts[address].ValidatorTopUp = stakedValidators.ValidatorTopUp
		mergedAccounts[address].ValidatorTopUpNum = stakedValidators.ValidatorTopUpNum
	}

	for address, stakedDelegators := range delegators {
		_, ok := mergedAccounts[address]
		if !ok {
			mergedAccounts[address] = stakedDelegators

			allAddresses = append(allAddresses, address)
			continue
		}

		mergedAccounts[address].Delegation = stakedDelegators.Delegation
		mergedAccounts[address].DelegationNum = stakedDelegators.DelegationNum
	}

	for address, lkMexAccount := range lkMexAccountsWithStake {
		_, ok := mergedAccounts[address]
		if !ok {
			mergedAccounts[address] = lkMexAccount

			allAddresses = append(allAddresses, address)
			continue
		}

		mergedAccounts[address].LKMEXStake = lkMexAccount.LKMEXStake
		mergedAccounts[address].LKMEXStakeNum = lkMexAccount.LKMEXStakeNum
	}

	return mergedAccounts, allAddresses
}

// ComputeClonedAccountsIndex will compute cloned accounts index based on current epoch
func (ap *accountsProcessor) ComputeClonedAccountsIndex() (string, error) {
	log.Info("Compute name of the new index...")

	genericAPIResponse := &data.GenericAPIResponse{}
	err := ap.restClient.CallGetRestEndPoint(pathNodeStatusMeta, genericAPIResponse, core.GetEmptyApiCredentials())
	if err != nil {
		return "", err
	}
	if genericAPIResponse.Error != "" {
		return "", fmt.Errorf("cannot compute accounts index %s", genericAPIResponse.Error)
	}

	epoch := gjson.Get(string(genericAPIResponse.Data), "status.erd_epoch_number")

	return fmt.Sprintf("%s_%s", accountsIndex, epoch.String()), nil
}

func computeTotalBalance(balances ...string) (string, float64) {
	totalBalance := big.NewInt(0)
	totalBalanceFloat := float64(0)

	if len(balances) == 0 {
		return "0", 0
	}

	for _, balance := range balances {
		balanceBig, ok := big.NewInt(0).SetString(balance, 10)
		if !ok {
			continue
		}

		totalBalance = totalBalance.Add(totalBalance, balanceBig)
		totalBalanceFloat += core.ComputeBalanceAsFloat(balance)
	}

	return totalBalance.String(), totalBalanceFloat
}

// IsInterfaceNil returns true if the value under the interface is nil
func (ap *accountsProcessor) IsInterfaceNil() bool {
	return ap == nil
}
