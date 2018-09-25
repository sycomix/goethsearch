package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"
)

//MultiSearch - Многопоточный поиск адресов
func MultiSearch(ctx *cli.Context) error {
	pattern := regexp.MustCompile(ctx.String("pattern"))
	threads := ctx.Int("threads")
	count := ctx.Int("count")
	check := ctx.Bool("balance")

	found := 0
	rpc, err := ethclient.Dial("https://mainnet.infura.io/UoiDPLlaAoM5GmpK8aR3")
	if err != nil {
		panic("Connection error")
	}
	for i := 0; i < threads; i++ {
		go func(num int) {
			ctx := context.Background()
			for count > found {
				key, err := crypto.GenerateKey()

				if err != nil {
					panic("Key generation error")
				}

				address := crypto.PubkeyToAddress(key.PublicKey).Hex()
				privateKey := hex.EncodeToString(key.D.Bytes())

				if pattern.MatchString(string(address)) {
					if check {
						balance, _ := rpc.BalanceAt(ctx, common.HexToAddress(address), nil)
						if balance.String() != "0" {
							found++
						}
						fmt.Printf("[%d][%d - %d] %s => %s [%d]\n", num, found, count, string(address), string(privateKey), balance)
					} else {
						found++
						fmt.Printf("[%d][%d - %d] %s => %s \n", num, found, count, string(address), string(privateKey))
					}
				}
			}
		}(i)
	}

	fmt.Scanln()

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "goethsearch"
	app.Usage = "search private keys by address pattern"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "pattern, p",
			Value: "(.*)",
			Usage: "Pattern for search",
		},
		cli.IntFlag{
			Name:  "threads, t",
			Value: 1,
			Usage: "Number of threads",
		},
		cli.IntFlag{
			Name:  "count, c",
			Value: 1,
			Usage: "Number of keys",
		},
		cli.BoolFlag{
			Name:  "balance, b",
			Usage: "Check balance",
		},
	}

	app.Action = MultiSearch

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
