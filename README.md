# gnark-solidity-checker

`gnark-solidity-checker generate` is a helper to compile gnark solidity verification circuits using `solc`,
generate go bindings using `abigen` and submit a proof running on geth simulated backend using `gnark-solidity-checker verify`.

## Install dependencies

### Install Solidity

`brew install solidity`
or

```bash
sudo add-apt-repository ppa:ethereum/ethereum
sudo apt-get update
sudo apt-get install solc
```

### Install Go

`brew install golang`

### Install `abigen`

`go install github.com/ethereum/go-ethereum/cmd/abigen@v1.12.0`
