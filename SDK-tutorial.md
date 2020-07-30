# Introduction

Pylons SDK can be used for blockchain game developers to build their own games that run on blockchain involving game characters and items.

Important points in pylons SDK are:

- item: Item is an important concept in Pylons sdk. It can simply be considered as the items in original games. Items will have properties in the format of Double, String, Integer etc. Characters can be created as items either.
- recipe: in this sdk, recipes will decide and generate everything like item generation and modification, combination of items. Recipes can be used to get result of some action between items and characters. Every action related to items are taken by recipes.
- trade: Pylons SDK enables accounts to trade their coins themselves. Trade includes items - items trading, coins - coins trading and mixed trading.

## Setup Local Environment

To setup local environment, you need to install Golang first.
The easiest way is using homebrew.

```brew install golang```

After go install is finished, you need to install pylons daemon and pylons cli.
This should be done inside the pylons repo. (https://github.com/Pylons-tech/pylons)

```go clean -i all
go install ./cmd/pylonsd
go install ./cmd/pylonscli
```

Create a genesis block and some test users

```
# Initialize configuration files and genesis file, the name here is "masternode", you can call it anything
pylonsd init masternode --chain-id pylonschain

# Copy the `Address` output here and save it for later use 
# [optional] add "--ledger" at the end to use a Ledger Nano S 
pylonscli keys add jack

# Copy the `Address` output here and save it for later use
pylonscli keys add alice

# Add both accounts, with coins to the genesis file
pylonsd add-genesis-account $(pylonscli keys show jack -a) 100pylon,1000jackcoin
pylonsd add-genesis-account $(pylonscli keys show alice -a) 100pylon,1000alicecoin

# Configure your CLI to eliminate need for chain-id flag
pylonscli config chain-id pylonschain
pylonscli config output json
pylonscli config indent true
pylonscli config trust-node true
```

In the pylons repo, there is init_accounts.local.sh file, which creates several test accounts called michael, eugen, jose and node0.
Node0 account is used for genesis, so you should avoid using this account for test purpose.

You can run `init_accounts.local.sh` file for test accounts setup

```
sh init_accounts.local.sh
```

After that, start pylons daemon

```
pylonsd start
```

## Setup Pylons SDK

```
git clone https://github.com/Pylons-tech/pylons_sdk
```

## Integration Test
Check integration test is working fine.
```
make int_tests
```

## Fixture Test
Check fixture test is working fine.
```
make fixture_tests ARGS="--accounts michael,eugen"
```

- Fixture tests args
Fixture tests have 2 args, accounts and scenarios.

Accounts arg will point the test accounts that will be used in the fixture test.
The pylons fixture test uses placeholder accounts, so the test accounts that are pointed in this arg will be replaced to the placeholder accounts in the fixture tests.

Scenarios arg will point the scenarios that will run in the fixture tests.
Pylons fixture tests have several scenarios - that simulate real gaming cases with json format. And you can point specific scenarios to run in this args

```
make fixture_tests ARGS="--scenarios submarin"
```








