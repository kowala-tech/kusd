﻿# TokenMock.sol

**TokenMock**

## Contract Members
**Constants & Variables**

```js
//public members
uint256 public totalSupply;
//private members
mapping(address => uint256) private balances;
```

**Events**

```js
event Transfer(address indexed from, address indexed to, uint256 value, bytes indexed data);
```

## Functions

- [balanceOf](#balanceof)
- [name](#name)
- [symbol](#symbol)
- [decimals](#decimals)
- [totalSupply](#totalsupply)
- [transfer](#transfer)
- [transfer](#transfer)
- [transfer](#transfer)

### balanceOf

```js
function balanceOf(address who) public view
returns(uint256)
```

**Arguments**

| Name        | Type           | Description  |
| ------------- |------------- | -----|
| who | address |  | 

### name

```js
function name() public view
returns(_name string)
```

### symbol

```js
function symbol() public view
returns(_symbol string)
```

### decimals

```js
function decimals() public view
returns(_decimals uint8)
```

### totalSupply

```js
function totalSupply() public view
returns(_supply uint256)
```

### transfer

```js
function transfer(address to, uint256 value) public
returns(ok bool)
```

**Arguments**

| Name        | Type           | Description  |
| ------------- |------------- | -----|
| to | address |  | 
| value | uint256 |  | 

### transfer

```js
function transfer(address to, uint256 value, bytes data) public
returns(ok bool)
```

**Arguments**

| Name        | Type           | Description  |
| ------------- |------------- | -----|
| to | address |  | 
| value | uint256 |  | 
| data | bytes |  | 

### transfer

```js
function transfer(address to, uint256 value, bytes data, string custom_fallback) public
returns(ok bool)
```

**Arguments**

| Name        | Type           | Description  |
| ------------- |------------- | -----|
| to | address |  | 
| value | uint256 |  | 
| data | bytes |  | 
| custom_fallback | string |  | 

## Contracts

- [BalanceContract](BalanceContract.md)
- [CappedToken](CappedToken.md)
- [Consensus](Consensus.md)
- [ConsensusMock](ConsensusMock.md)
- [DomainResolver](DomainResolver.md)
- [DomainResolverMock](DomainResolverMock.md)
- [ExchangeMgr](ExchangeMgr.md)
- [FIFSRegistrar](FIFSRegistrar.md)
- [Initializable](Initializable.md)
- [KNS](KNS.md)
- [KNSRegistry](KNSRegistry.md)
- [KNSRegistryV1](KNSRegistryV1.md)
- [KRC223](KRC223.md)
- [Math](Math.md)
- [Migrations](Migrations.md)
- [MiningToken](MiningToken.md)
- [MintableToken](MintableToken.md)
- [MultiSigWallet](MultiSigWallet.md)
- [NameHash](NameHash.md)
- [OracleMgr](OracleMgr.md)
- [Ownable](Ownable.md)
- [Pausable](Pausable.md)
- [PublicResolver](PublicResolver.md)
- [SafeMath](SafeMath.md)
- [strings](strings.md)
- [SystemVars](SystemVars.md)
- [Token](Token.md)
- [TokenMock](TokenMock.md)
- [TokenReceiver](TokenReceiver.md)
- [ValidatorMgr](ValidatorMgr.md)
