# CieloGo

CieloGo is a Go client library for interacting with the Cielo Finance API, allowing developers to easily access wallet analytics, pnl stats, related wallets, and wallet tags through a simple and intuitive interface. Built for efficiency and ease of use, CieloGo abstracts the complexity of direct API calls into straightforward Go functions.

Please note, CieloGo is an unofficial client library and is not endorsed or maintained by Cielo Finance.

## Features

- Get feed information with customizable queries.
- Retrieve Non-Fungible Token (NFT) Profit and Loss (PnL) statistics for specific wallets.
- Fetch token PnL data for given wallets.
- Obtain aggregated token PnL stats.
- List related wallets to a given wallet.
- Access tags associated with a wallet.

## Installation

To use CieloGo in your Go project, install it as a module:

```bash
go get github.com/sealtv/cielogo
```

## Usage

To start using CieloGo, import it into your Go project:

```go
import "github.com/sealtv/cielogo"
```

### Initializing the Client

```go
client := cielogo.NewClient()
```

### Making Requests

Here are some examples of how you might call various methods of the CieloGo client.

**Get Feed:**

```go
ctx := context.Background()
req := &apiv1.FeedRequest{...} // setup your request
feed, err := client.GetFeedV1(ctx, req)
```

**Get NFTs PnL:**

```go
nftsPnLReq := apiv1.NftsPnLRequest{Wallet: "0xWALLET_ADDRESS"}
nftsPnl, err := client.GetNftsPnlV1(ctx, &nftsPnLReq)
```

**Get Tokens PnL:**

```go
tokensPnLReq := apiv1.TokensPnLRequest{Wallet: "0xWALLET_ADDRESS"}
tokensPnl, err := client.GetTokensPnlV1(ctx, &tokensPnLReq)
```

For more details on each request and response structure, refer to the [Cielo Finance API documentation](https://developer.cielo.finance).

## Contributing

Contributions to the CieloGo project are welcome. Please feel free to report any bugs, suggest features, or open pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
