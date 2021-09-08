# Numary Ledger [![test](https://github.com/numary/ledger/actions/workflows/main.yml/badge.svg)](https://github.com/numary/ledger/actions/workflows/main.yml)

Numary is a programmable financial ledger that wants to make building financial apps safe, fun and cheap.

Building financial software is both critical and notably hard. The [same bugs](https://medium.com/selency-tech-product/your-balance-is-0-30000000004-b6f7870bd32e) are repeated again and again, paving the highway to catastrophes.

Numary wants to tackle this issue with a ledger that provides atomic multi-postings transactions and is programmable in [Numscript](https://github.com/numary/machine), a built-in language dedicated to money movements. It will shine for apps that require a lot of custom, money-touching code such as:

* E-commerce with complex payments flows, payments splitting, such as marketplaces
* Company-issued currencies systems, e.g. Twitch Bits
* In-game currencies, inventories and trading systems, e.g. Fortnite V-Bucks
* Payment gateways using non-standard assets, e.g. learning credits
* Local currencies and complementary finance

# Getting started

Numary works as a standalone binary, the latest of which can be downloaded from the [releases page](https://github.com/numary/ledger/releases). You can move the binary to any executable path, such as to `/usr/local/bin`.

```SHELL

numary server start

# Submit a first transaction
echo "
send [USD/2 599] (
  source = @world
  destination = @payments:001
)

send [USD/2 599] (
  source = @payments:001
  destination = @rides:0234
)

send [USD/2 599] (
  source = @rides:0234
  destination = {
    85/100 to @drivers:042
    15/100 to @platform:fees
  }
)
" > example.num

numary exec quickstart example.num

# Get the balances of drivers:042
curl -X GET http://localhost:3068/quickstart/accounts/drivers:042

# List transactions
curl -X GET http://localhost:3068/quickstart/transactions
```

# Documentation

You can find the complete Numary documentation at [docs.numary.com](https://docs.numary.com)

# Dashboard

<img width="400" alt="control-screenshot" src="https://user-images.githubusercontent.com/1770991/126158742-393ac0d0-1048-4b57-a7fd-7381f3da2ca8.png">

A simple [dashboard](https://github.com/numary/control) is built in the ledger binary, to make it easier to visualize transactions. It can be started with:

```SHELL
numary ui
```

# Roadmap

We keep an open roadmap of the upcoming releases and features [here](https://numary.notion.site/OSS-Roadmap-4535fa5716fb4f618027201afcc6f204).
