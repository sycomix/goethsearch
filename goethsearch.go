package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"regexp"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli"
)

func MultiSearch(ctx *cli.Context) error {
	pattern := regexp.MustCompile(ctx.String("pattern"))
	threads := ctx.Int("threads")
	count := ctx.Int("count")

	found := 0

	for i := 0; i < threads; i++ {
		go func(num int) {
			for count > found {
				key, err := crypto.GenerateKey()

				if err != nil {
					panic("Key generation error")
				}

				address := crypto.PubkeyToAddress(key.PublicKey).Hex()
				privateKey := hex.EncodeToString(key.D.Bytes())

				if pattern.MatchString(string(address)) {
					fmt.Printf("[%d][%d - %d] %s => %s \n", num, found+1, count, string(address), string(privateKey))
					found++
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
	}

	app.Action = MultiSearch

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
