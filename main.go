package main

import (
	"bsc/bsc"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"

	"github.com/ethereum/go-ethereum/common"
)

type Message struct {
	Address string
	AssetId string
}

type ApiResponse struct {
	Data    Data   `json:"data"`
	Network string `json:"network"`
}

type Data struct {
	Addresses []string `json:"addresses"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//ctx := context.Background()
	bnbRpc := "wss://admittedly-next-kitten.ngrok-free.app/"

	contracts := []bsc.Contract{
		{Address: common.HexToAddress("0x55d398326f99059fF775485246999027B3197955"), AssetID: 1},
		{Address: common.HexToAddress("0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d"), AssetID: 2},
		{Address: common.HexToAddress("0xe9e7cea3dedca5984780bafc599bd69add087d56"), AssetID: 3},
	} // 0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599

	ethClient := bsc.NewEthClient(bnbRpc, contracts...)

	//apiResponse, err := callApiForAddresses("http://internal-diam-stake-processor-147803691.us-west-1.elb.amazonaws.com:8000/get-addresses")
	// apiResponse, err := callApiForAddresses("http://localhost:8000/get-addresses")
	// if err != nil {
	// 	log.Printf("Failed to fetch addresses from API: %v", err)
	// }

	// for _, addr := range apiResponse.Data.Addresses {
	// 	//log.Println(addr, " loaded from db")
	// 	ethClient.AddWatchedAddress(addr)
	// }

	// Hardcoded test addresses
	testAddresses := []string{
		"0x4a34d3870C0496a4Be4Bd32F44B29A71c9b7f7F3",
		"0x4c9630397E08C8b54375ad8608025Db7c7026840",
	}

	for _, addr := range testAddresses {
		ethClient.AddWatchedAddress(addr)
		log.Printf("Added test address: %s", addr)
	}
	// //cfg, err := config.LoadDefaultConfig(ctx)
	// cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-1"))
	// if err != nil {
	// 	log.Printf("Failed to load AWS config: %v", err)
	// }

	go func() {
		if err := ethClient.Start(); err != nil {
			log.Printf("Failed to start BNB-Client: %v", err)
		}
	}()

	log.Println("BNB-Client has started listening for new blocks. Watched addresses loaded.")

	app := fiber.New()
	app.Get("/metrics", monitor.New(monitor.Config{Title: "BSC listner Metrics Page"}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, from BSC listner")
	})
	go func() {
		log.Fatal(app.Listen(":3002"))
	}()

	// Keep running until interrupted
	select {}
}
