# Introduction

Pylons SDK can be used for blockchain game developers to build their own games that run on blockchain involving game characters and items.

Important points in pylons SDK are:

- item: Item is an important concept in Pylons sdk. It can simply be considered as the items in original games. Items will have properties in the format of Double, String, Integer etc. Characters can be created as items either.
- recipe: in this sdk, recipes will decide and generate everything like item generation and modification, combination of items. Recipes can be used to get result of some action between items and characters. Every action related to items are taken by recipes.
- trade: Pylons SDK enables accounts to trade their coins themselves. Trade includes items - items trading, coins - coins trading and mixed trading.

You need to know some basic knowledge about `cosmos-sdk` as our Pylons sdk is based on `cosmos-sdk`.
You can check more details in the following link: https://cosmos.network/sdk

## Setup Local Environment

### Install Golang
  
Pylons sdk is based on Golang. So you need to install Go to setup local environment.
The easiest way is to use homebrew. You can follow this link:
https://ahmadawais.com/install-go-lang-on-macos-with-homebrew/

### Start pylons daemon

You need to start pylons daemon to make the sdk work. Download binary files from Pylons and run the following command.

```
sh init_accounts.local.sh
pylonsd start
```

`init_accounts.local.sh` will create initial test accounts named `node0`, `michael`, `eugen` and `jose`. These accounts will be used in the pylonscli tests.

# Pylonscli

Now you are ready to use pylonscli. Open a new terminal window other than the pylons daemon terminal window. Try to run the pylonscli command.

```
pylonscli
```
This command will give you the overall help for all available commands in pylonscli.
The result should be like the following.

```
The Pylons Client

Usage:
  pylonscli [command]

Available Commands:
  status      Query remote node for status
  config      Create or query an application CLI configuration file
  query       Querying subcommands
  tx          Transactions subcommands
              
  rest-server Start LCD (light-client daemon), a local REST server
              
  keys        Add or view local private keys
              
  help        Help about any command

Flags:
      --chain-id string          Chain ID of tendermint node
  -e, --encoding string          Binary encoding (hex|b64|btc) (default "hex")
  -h, --help                     help for pylonscli
      --home string              directory for config and data (default "/Users/ghostprince/.pylonscli")
      --keyring-backend string   keyring backend of tendermint node
  -o, --output string            Output format (text|json) (default "text")
      --trace                    print out full stack trace on errors

Use "pylonscli [command] --help" for more information about a command.
```

Now let's move forward with some important cli commands one by one.

As we are on the local environment, `--keyring-backend=test` flag should be set for every command we do from now.


### Add private key to the chain

`pylonscli keys` command will let you add or view local private keys in the chain.

- Add your first private key
  ```
  pylonscli keys add test --keyring-backend=test
  ```

  You can replace the private key name that is set `test` in the command with any name you want.

  This command will create a private key with the name given as the argument. The result will be like following.

  ```
  {
    "name": "test",
    "type": "local",
    "address": "cosmos1fun8le2dxrclr633psv7gke6wtlycunnm8dlm7",
    "pubkey": "cosmospub1addwnpepqv6ppfkfu7cm62a2n3qfjr6r2h5dcssrywly59n699v9g6gykccq55hzewf",
    "mnemonic": "choose crater shift until worth wasp win pilot again piece canyon habit mercy come crisp next captain street horn inmate word vapor cake pledge"
  }
  ```

  The `address` value `cosmos1fun8le2dxrclr633psv7gke6wtlycunnm8dlm7` will be used for the other cli commands.

- Show the created private key
  ```
  pylonscli keys show test --keyring-backend=test
  ```
  This will show the info of the private key `test` that is just added.
  The result will be like following.
  ```
  {
    "name": "test",
    "type": "local",
    "address": "cosmos1fun8le2dxrclr633psv7gke6wtlycunnm8dlm7",
    "pubkey": "cosmospub1addwnpepqv6ppfkfu7cm62a2n3qfjr6r2h5dcssrywly59n699v9g6gykccq55hzewf"
  }
  ```

- List all the keys available
  ```
  pylonscli keys list --keyring-backend=test
  ```
  This command will list all the private keys that are available in the chain.
  The result will be like the following.
  ```
  [
    {
      "name": "eugen",
      "type": "local",
      "address": "cosmos1g5w79thfvt86m6cpa0a7jezfv0sjt0u7y09ldm",
      "pubkey": "cosmospub1addwnpepqgmz48l2urqa69zd9djrah2qlc4u7n6gp98ch3ydkyfpm78ktn6xuf7w2rh"
    },
    {
      "name": "michael",
      "type": "local",
      "address": "cosmos1k6qm04kxkz7q69lhy80jf562y8d5rj66y8g8t2",
      "pubkey": "cosmospub1addwnpepqdap2dhf7aj98mgewqjn94t9gfujt4k8u3ztdyeaqsaslfd98eqr2ed8z55"
    },
    {
      "name": "node0",
      "type": "local",
      "address": "cosmos13p8890funv54hflk82ju0zv47tspglpk373453",
      "pubkey": "cosmospub1addwnpepq2ht7s5t3kp7058w2kntx9ha8av396xv78nhs6lszqtcwtf6kwdm20axerv"
    },
    {
      "name": "test",
      "type": "local",
      "address": "cosmos1fun8le2dxrclr633psv7gke6wtlycunnm8dlm7",
      "pubkey": "cosmospub1addwnpepqv6ppfkfu7cm62a2n3qfjr6r2h5dcssrywly59n699v9g6gykccq55hzewf"
    }
  ]
  ```

### Create and update cookbook

All the transactions in pylons sdk will be done based on cookbook. You should create cookbook first to do any transaction in pylons sdk.

