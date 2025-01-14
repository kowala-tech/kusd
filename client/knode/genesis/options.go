package genesis

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/kowala-tech/kcoin/client/common"
	"github.com/kowala-tech/kcoin/client/params"
	"github.com/pkg/errors"
)

const (
	MainNetwork  = "main"
	TestNetwork  = "test"
	OtherNetwork = "other"

	KonsensusConsensus = "konsensus"
)

var (
	availableNetworks = map[string]bool{
		MainNetwork:  true,
		TestNetwork:  true,
		OtherNetwork: true,
	}

	availableConsensusEngines = map[string]bool{
		KonsensusConsensus: true,
	}

	ErrEmptyMaxNumValidators             = errors.New("max number of validators is mandatory")
	ErrInvalidMaxNumValidators           = errors.New("invalid max num of validators")
	ErrEmptyFreezePeriod                 = errors.New("freeze period in days is mandatory")
	ErrInvalidFreezePeriod               = errors.New("freeze period is invalid")
	ErrEmptyWalletAddressValidator       = errors.New("wallet address of genesis validator is mandatory")
	ErrInvalidWalletAddressValidator     = errors.New("wallet address of genesis validator is invalid")
	ErrInvalidAddressInPrefundedAccounts = errors.New("address in prefunded accounts is invalid")
	ErrInvalidContractsOwnerAddress      = errors.New("address used for smart contracts is invalid")
	ErrInvalidNetwork                    = errors.New("invalid Network, use main, test or other")
	ErrInvalidConsensusEngine            = errors.New("invalid consensus engine")
	ErrInvalidAddress                    = errors.New("Invalid address")
)

type Options struct {
	Network           string
	BlockNumber       uint64
	SystemVars        *SystemVarsOpts
	Governance        *GovernanceOpts
	Consensus         *ConsensusOpts
	StabilityContract *StabilityContractOpts
	DataFeedSystem    *DataFeedSystemOpts
	PrefundedAccounts []PrefundedAccount
	ExtraData         string
}

type StabilityContractOpts struct {
	MinDeposit uint64
}

type SystemVarsOpts struct {
	InitialPrice float64
}

type TokenHolder struct {
	Address   string
	NumTokens uint64
}

type MiningTokenOpts struct {
	Name     string
	Symbol   string
	Cap      uint64
	Decimals uint64
	Holders  []TokenHolder
}

type ConsensusOpts struct {
	Engine           string
	MaxNumValidators uint64
	FreezePeriod     uint64
	BaseDeposit      uint64
	SuperNodeAmount  uint64
	Validators       []Validator
	MiningToken      *MiningTokenOpts
}

type GovernanceOpts struct {
	Origin           string
	Governors        []string
	NumConfirmations uint64
}

type PriceOpts struct {
	SyncFrequency uint64
	UpdatePeriod  uint64
}

type DataFeedSystemOpts struct {
	MaxNumOracles uint64
	Price         PriceOpts
}

type PrefundedAccount struct {
	Address string
	Balance uint64
}

type Validator struct {
	Address string
	Deposit uint64
}

type validValidator struct {
	address common.Address
	deposit *big.Int
}

type validValidatorMgrOpts struct {
	maxNumValidators *big.Int
	freezePeriod     *big.Int
	baseDeposit      *big.Int
	superNodeAmount  *big.Int
	validators       []*validValidator
	miningTokenAddr  common.Address
	owner            common.Address
}

type validPriceOpts struct {
	syncFrequency *big.Int
	updatePeriod  *big.Int
}

type validOracleMgrOpts struct {
	maxNumOracles    *big.Int
	price            validPriceOpts
	validatorMgrAddr common.Address
	owner            common.Address
}

type validSystemVarsOpts struct {
	initialPrice  *big.Int
	initialSupply *big.Int
	owner         common.Address
}

type validStabilityContractOpts struct {
	minDeposit     *big.Int
	systemVarsAddr common.Address
}

type validTokenHolder struct {
	address common.Address
	balance *big.Int
}

type validMiningTokenOpts struct {
	name     string
	symbol   string
	cap      *big.Int
	decimals *big.Int
	holders  []*validTokenHolder
	owner    common.Address
}

type validMultiSigOpts struct {
	multiSigCreator  *common.Address
	multiSigOwners   []common.Address
	numConfirmations *big.Int
}

type validPrefundedAccount struct {
	accountAddress *common.Address
	balance        *big.Int
}

type validGenesisOptions struct {
	network           string
	blockNumber       uint64
	consensusEngine   string
	prefundedAccounts []*validPrefundedAccount
	multiSig          *validMultiSigOpts
	validatorMgr      *validValidatorMgrOpts
	oracleMgr         *validOracleMgrOpts
	miningToken       *validMiningTokenOpts
	sysvars           *validSystemVarsOpts
	stability         *validStabilityContractOpts
	ExtraData         string
}

func validateOptions(options Options) (*validGenesisOptions, error) {
	network, err := mapNetwork(options.Network)
	if err != nil {
		return nil, err
	}

	consensusEngine := KonsensusConsensus
	if options.Consensus.Engine != "" {
		consensusEngine, err = mapConsensusEngine(options.Consensus.Engine)
		if err != nil {
			return nil, err
		}
	}

	// sysvars
	initialPrice := new(big.Int)
	new(big.Float).Mul(new(big.Float).SetFloat64(options.SystemVars.InitialPrice), big.NewFloat(params.Kcoin)).Int(initialPrice)

	// governance
	multiSigCreator, err := getAddress(options.Governance.Origin)
	if err != nil {
		return nil, err
	}

	multiSigOwners := make([]common.Address, 0, len(options.Governance.Governors))
	for _, governor := range options.Governance.Governors {
		owner, err := getAddress(governor)
		if err != nil {
			return nil, err
		}
		multiSigOwners = append(multiSigOwners, *owner)
	}

	numConfirmations := new(big.Int).SetUint64(options.Governance.NumConfirmations)

	// consensus
	maxNumValidators := new(big.Int).SetUint64(options.Consensus.MaxNumValidators)
	consensusBaseDeposit := new(big.Int).Mul(new(big.Int).SetUint64(options.Consensus.BaseDeposit), big.NewInt(params.Kcoin))
	consensusFreezePeriod := new(big.Int).SetUint64(options.Consensus.FreezePeriod)
	superNodeAmount := new(big.Int).Mul(new(big.Int).SetUint64(options.Consensus.SuperNodeAmount), big.NewInt(params.Kcoin))

	validators := make([]*validValidator, 0, len(options.Consensus.Validators))
	for _, validator := range options.Consensus.Validators {
		addr, err := getAddress(validator.Address)
		if err != nil {
			return nil, err
		}
		validators = append(validators, &validValidator{
			address: *addr,
			deposit: new(big.Int).Mul(new(big.Int).SetUint64(validator.Deposit), big.NewInt(params.Kcoin)),
		})
	}

	// stability contract
	minDeposit := new(big.Int).Mul(new(big.Int).SetUint64(options.StabilityContract.MinDeposit), big.NewInt(params.Kcoin))

	// data feed system
	maxNumOracles := new(big.Int).SetUint64(options.DataFeedSystem.MaxNumOracles)
	syncFrequency := new(big.Int).SetUint64(options.DataFeedSystem.Price.SyncFrequency)
	updatePeriod := new(big.Int).SetUint64(options.DataFeedSystem.Price.UpdatePeriod)

	// mining tokens
	decimals := new(big.Int).Exp(common.Big1, new(big.Int).SetUint64(options.Consensus.MiningToken.Decimals), nil)
	cap := new(big.Int).Mul(new(big.Int).SetUint64(options.Consensus.MiningToken.Cap), big.NewInt(params.Kcoin))

	holders := make([]*validTokenHolder, 0, len(options.Consensus.MiningToken.Holders))
	for _, holder := range options.Consensus.MiningToken.Holders {
		addr, err := getAddress(holder.Address)
		if err != nil {
			return nil, err
		}
		holders = append(holders, &validTokenHolder{
			address: *addr,
			balance: new(big.Int).Mul(new(big.Int).SetUint64(holder.NumTokens), big.NewInt(params.Kcoin)),
		})
	}

	// prefund accounts
	mintedAmount, validPrefundedAccounts, err := mapPrefundedAccounts(options.PrefundedAccounts)
	if err != nil {
		return nil, err
	}

	return &validGenesisOptions{
		network:         network,
		blockNumber:     options.BlockNumber,
		consensusEngine: consensusEngine,
		sysvars: &validSystemVarsOpts{
			initialPrice:  initialPrice,
			initialSupply: mintedAmount,
		},
		multiSig: &validMultiSigOpts{
			multiSigCreator:  multiSigCreator,
			multiSigOwners:   multiSigOwners,
			numConfirmations: numConfirmations,
		},
		validatorMgr: &validValidatorMgrOpts{
			maxNumValidators: maxNumValidators,
			freezePeriod:     consensusFreezePeriod,
			baseDeposit:      consensusBaseDeposit,
			superNodeAmount:  superNodeAmount,
			validators:       validators,
		},
		oracleMgr: &validOracleMgrOpts{
			maxNumOracles: maxNumOracles,
			price: validPriceOpts{
				syncFrequency: syncFrequency,
				updatePeriod:  updatePeriod,
			},
		},
		miningToken: &validMiningTokenOpts{
			name:     options.Consensus.MiningToken.Name,
			symbol:   options.Consensus.MiningToken.Symbol,
			cap:      cap,
			decimals: decimals,
			holders:  holders,
		},
		stability: &validStabilityContractOpts{
			minDeposit: minDeposit,
		},
		prefundedAccounts: validPrefundedAccounts,
		ExtraData:         options.ExtraData,
	}, nil
}

func getAddress(s string) (*common.Address, error) {
	if !common.IsHexAddress(s) {
		return nil, fmt.Errorf("%s:%s", ErrInvalidAddress, s)
	}
	address := common.HexToAddress(s)

	return &address, nil
}

func mapNetwork(network string) (string, error) {
	if !availableNetworks[network] {
		return "", fmt.Errorf("%v:%s", ErrInvalidNetwork, network)
	}

	return network, nil
}

func mapConsensusEngine(consensus string) (string, error) {
	if !availableConsensusEngines[consensus] {
		return "", ErrInvalidConsensusEngine
	}

	return consensus, nil
}

func mapWalletAddress(a string) (*common.Address, error) {
	stringAddr := a

	if text := strings.TrimSpace(a); text == "" {
		return nil, ErrEmptyWalletAddressValidator
	}

	if strings.HasPrefix(stringAddr, "0x") {
		stringAddr = strings.TrimPrefix(stringAddr, "0x")
	}

	if len(stringAddr) != 40 {
		return nil, ErrInvalidWalletAddressValidator
	}

	bigaddr, _ := new(big.Int).SetString(stringAddr, 16)
	address := common.BigToAddress(bigaddr)

	return &address, nil
}

func mapPrefundedAccounts(accounts []PrefundedAccount) (*big.Int, []*validPrefundedAccount, error) {
	var validAccounts []*validPrefundedAccount

	mintedAmount := new(big.Int)
	for _, a := range accounts {
		address, err := mapWalletAddress(a.Address)
		if err != nil {
			return nil, nil, ErrInvalidAddressInPrefundedAccounts
		}

		balance := new(big.Int).Mul(new(big.Int).SetUint64(a.Balance), new(big.Int).SetUint64(params.Kcoin))
		mintedAmount.Add(mintedAmount, balance)

		validAccount := &validPrefundedAccount{
			accountAddress: address,
			balance:        balance,
		}

		validAccounts = append(validAccounts, validAccount)
	}

	return mintedAmount, validAccounts, nil
}
