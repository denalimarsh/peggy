// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generated

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// GeneratedABI is the input ABI used to generate the binding from.
const GeneratedABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"bridgeBank\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"oracle\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"valset\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"cosmosBridge\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_cosmosBridge\",\"type\":\"address\"},{\"name\":\"_bridgeBank\",\"type\":\"address\"},{\"name\":\"_oracle\",\"type\":\"address\"},{\"name\":\"_valset\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_cosmosBridge\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_bridgeBank\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_oracle\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_valset\",\"type\":\"address\"}],\"name\":\"LogContractsRegistered\",\"type\":\"event\"}]"

// GeneratedBin is the compiled bytecode used for deploying new contracts.
const GeneratedBin = `608060405234801561001057600080fd5b506040516080806105378339810180604052608081101561003057600080fd5b8101908080519060200190929190805190602001909291908051906020019092919080519060200190929190505050836000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555082600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555081600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507f039b733f31259b106f1d278c726870d5b28c7db22957d63df8dbaa70bd3a032a6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16604051808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200194505050505060405180910390a15050505061023c806102fb6000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80630e41f373146100515780637dc0d1d01461009b5780637f54af0c146100e5578063b0e9ef711461012f575b600080fd5b610059610179565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100a361019f565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100ed6101c5565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6101376101eb565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff168156fea165627a7a7230582075bae4621d0253e7fc36fad677762be2ae52714ed078b5bab9a582832874cc8b0029`

// DeployGenerated deploys a new Ethereum contract, binding an instance of Generated to it.
func DeployGenerated(auth *bind.TransactOpts, backend bind.ContractBackend, _cosmosBridge common.Address, _bridgeBank common.Address, _oracle common.Address, _valset common.Address) (common.Address, *types.Transaction, *Generated, error) {
	parsed, err := abi.JSON(strings.NewReader(GeneratedABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(GeneratedBin), backend, _cosmosBridge, _bridgeBank, _oracle, _valset)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Generated{GeneratedCaller: GeneratedCaller{contract: contract}, GeneratedTransactor: GeneratedTransactor{contract: contract}, GeneratedFilterer: GeneratedFilterer{contract: contract}}, nil
}

// Generated is an auto generated Go binding around an Ethereum contract.
type Generated struct {
	GeneratedCaller     // Read-only binding to the contract
	GeneratedTransactor // Write-only binding to the contract
	GeneratedFilterer   // Log filterer for contract events
}

// GeneratedCaller is an auto generated read-only Go binding around an Ethereum contract.
type GeneratedCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeneratedTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GeneratedTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeneratedFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GeneratedFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeneratedSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GeneratedSession struct {
	Contract     *Generated        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GeneratedCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GeneratedCallerSession struct {
	Contract *GeneratedCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// GeneratedTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GeneratedTransactorSession struct {
	Contract     *GeneratedTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// GeneratedRaw is an auto generated low-level Go binding around an Ethereum contract.
type GeneratedRaw struct {
	Contract *Generated // Generic contract binding to access the raw methods on
}

// GeneratedCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GeneratedCallerRaw struct {
	Contract *GeneratedCaller // Generic read-only contract binding to access the raw methods on
}

// GeneratedTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GeneratedTransactorRaw struct {
	Contract *GeneratedTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGenerated creates a new instance of Generated, bound to a specific deployed contract.
func NewGenerated(address common.Address, backend bind.ContractBackend) (*Generated, error) {
	contract, err := bindGenerated(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Generated{GeneratedCaller: GeneratedCaller{contract: contract}, GeneratedTransactor: GeneratedTransactor{contract: contract}, GeneratedFilterer: GeneratedFilterer{contract: contract}}, nil
}

// NewGeneratedCaller creates a new read-only instance of Generated, bound to a specific deployed contract.
func NewGeneratedCaller(address common.Address, caller bind.ContractCaller) (*GeneratedCaller, error) {
	contract, err := bindGenerated(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GeneratedCaller{contract: contract}, nil
}

// NewGeneratedTransactor creates a new write-only instance of Generated, bound to a specific deployed contract.
func NewGeneratedTransactor(address common.Address, transactor bind.ContractTransactor) (*GeneratedTransactor, error) {
	contract, err := bindGenerated(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GeneratedTransactor{contract: contract}, nil
}

// NewGeneratedFilterer creates a new log filterer instance of Generated, bound to a specific deployed contract.
func NewGeneratedFilterer(address common.Address, filterer bind.ContractFilterer) (*GeneratedFilterer, error) {
	contract, err := bindGenerated(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GeneratedFilterer{contract: contract}, nil
}

// bindGenerated binds a generic wrapper to an already deployed contract.
func bindGenerated(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GeneratedABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Generated *GeneratedRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Generated.Contract.GeneratedCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Generated *GeneratedRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Generated.Contract.GeneratedTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Generated *GeneratedRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Generated.Contract.GeneratedTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Generated *GeneratedCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Generated.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Generated *GeneratedTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Generated.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Generated *GeneratedTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Generated.Contract.contract.Transact(opts, method, params...)
}

// BridgeBank is a free data retrieval call binding the contract method 0x0e41f373.
//
// Solidity: function bridgeBank() constant returns(address)
func (_Generated *GeneratedCaller) BridgeBank(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Generated.contract.Call(opts, out, "bridgeBank")
	return *ret0, err
}

// BridgeBank is a free data retrieval call binding the contract method 0x0e41f373.
//
// Solidity: function bridgeBank() constant returns(address)
func (_Generated *GeneratedSession) BridgeBank() (common.Address, error) {
	return _Generated.Contract.BridgeBank(&_Generated.CallOpts)
}

// BridgeBank is a free data retrieval call binding the contract method 0x0e41f373.
//
// Solidity: function bridgeBank() constant returns(address)
func (_Generated *GeneratedCallerSession) BridgeBank() (common.Address, error) {
	return _Generated.Contract.BridgeBank(&_Generated.CallOpts)
}

// CosmosBridge is a free data retrieval call binding the contract method 0xb0e9ef71.
//
// Solidity: function cosmosBridge() constant returns(address)
func (_Generated *GeneratedCaller) CosmosBridge(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Generated.contract.Call(opts, out, "cosmosBridge")
	return *ret0, err
}

// CosmosBridge is a free data retrieval call binding the contract method 0xb0e9ef71.
//
// Solidity: function cosmosBridge() constant returns(address)
func (_Generated *GeneratedSession) CosmosBridge() (common.Address, error) {
	return _Generated.Contract.CosmosBridge(&_Generated.CallOpts)
}

// CosmosBridge is a free data retrieval call binding the contract method 0xb0e9ef71.
//
// Solidity: function cosmosBridge() constant returns(address)
func (_Generated *GeneratedCallerSession) CosmosBridge() (common.Address, error) {
	return _Generated.Contract.CosmosBridge(&_Generated.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() constant returns(address)
func (_Generated *GeneratedCaller) Oracle(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Generated.contract.Call(opts, out, "oracle")
	return *ret0, err
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() constant returns(address)
func (_Generated *GeneratedSession) Oracle() (common.Address, error) {
	return _Generated.Contract.Oracle(&_Generated.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() constant returns(address)
func (_Generated *GeneratedCallerSession) Oracle() (common.Address, error) {
	return _Generated.Contract.Oracle(&_Generated.CallOpts)
}

// Valset is a free data retrieval call binding the contract method 0x7f54af0c.
//
// Solidity: function valset() constant returns(address)
func (_Generated *GeneratedCaller) Valset(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Generated.contract.Call(opts, out, "valset")
	return *ret0, err
}

// Valset is a free data retrieval call binding the contract method 0x7f54af0c.
//
// Solidity: function valset() constant returns(address)
func (_Generated *GeneratedSession) Valset() (common.Address, error) {
	return _Generated.Contract.Valset(&_Generated.CallOpts)
}

// Valset is a free data retrieval call binding the contract method 0x7f54af0c.
//
// Solidity: function valset() constant returns(address)
func (_Generated *GeneratedCallerSession) Valset() (common.Address, error) {
	return _Generated.Contract.Valset(&_Generated.CallOpts)
}

// GeneratedLogContractsRegisteredIterator is returned from FilterLogContractsRegistered and is used to iterate over the raw logs and unpacked data for LogContractsRegistered events raised by the Generated contract.
type GeneratedLogContractsRegisteredIterator struct {
	Event *GeneratedLogContractsRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GeneratedLogContractsRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GeneratedLogContractsRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GeneratedLogContractsRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GeneratedLogContractsRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GeneratedLogContractsRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GeneratedLogContractsRegistered represents a LogContractsRegistered event raised by the Generated contract.
type GeneratedLogContractsRegistered struct {
	CosmosBridge common.Address
	BridgeBank   common.Address
	Oracle       common.Address
	Valset       common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterLogContractsRegistered is a free log retrieval operation binding the contract event 0x039b733f31259b106f1d278c726870d5b28c7db22957d63df8dbaa70bd3a032a.
//
// Solidity: event LogContractsRegistered(address _cosmosBridge, address _bridgeBank, address _oracle, address _valset)
func (_Generated *GeneratedFilterer) FilterLogContractsRegistered(opts *bind.FilterOpts) (*GeneratedLogContractsRegisteredIterator, error) {

	logs, sub, err := _Generated.contract.FilterLogs(opts, "LogContractsRegistered")
	if err != nil {
		return nil, err
	}
	return &GeneratedLogContractsRegisteredIterator{contract: _Generated.contract, event: "LogContractsRegistered", logs: logs, sub: sub}, nil
}

// WatchLogContractsRegistered is a free log subscription operation binding the contract event 0x039b733f31259b106f1d278c726870d5b28c7db22957d63df8dbaa70bd3a032a.
//
// Solidity: event LogContractsRegistered(address _cosmosBridge, address _bridgeBank, address _oracle, address _valset)
func (_Generated *GeneratedFilterer) WatchLogContractsRegistered(opts *bind.WatchOpts, sink chan<- *GeneratedLogContractsRegistered) (event.Subscription, error) {

	logs, sub, err := _Generated.contract.WatchLogs(opts, "LogContractsRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GeneratedLogContractsRegistered)
				if err := _Generated.contract.UnpackLog(event, "LogContractsRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
