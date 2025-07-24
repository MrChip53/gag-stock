package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/mrchip53/go-a-garden/pkg/gag"
)

var items = []string{
	"Mushroom",
	"Pepper",
	"Dragon Fruit",
	"Coconut",
	"Grape",
	"Cacao",
	"Beanstalk",
	"Ember Lily",
	"Sugar Apple",
	"Burning Bud",
	"Giant Pinecone",

	"Basic Sprinkler",
	"Advanced Sprinkler",
	"Godly Sprinkler",
	"Master Sprinkler",
	"Medium Toy",
	"Medium Treat",
	"Tanning Mirror",
	"Friendship Pot",

	"Zen Seed Pack",
	"Zen Egg",
	"Hot Spring",
	"Soft Sunshine",
	"Koi",
	"Spiked Mango",
	"Pet Shard Tranquil",

	"Common Summer Egg",
	"Rare Summer Egg",
	"Mythical Egg",
	"Paradise Egg",
	"Bug Egg",
}

var lastStockId int64

var lastSeedRestockTime int64
var lastGearRestockTime int64
var lastEggRestockTime int64
var lastMerchantRestockTime int64
var lastCosmeticRestockTime int64
var lastEventRestockTime int64

var firstRun = true

func notify(message string) {
	err := beeep.Alert("Grow a Garden", message, "")
	if err != nil {
		log.Printf("Failed to notify: %v", err)
	}
}

func checkStock(wantedItems []string) error {
	stock, err := gag.GetGagStock()
	if err != nil {
		return err
	}

	if stock.LastApiFetch == lastStockId {
		return nil
	}

	lastStockId = stock.LastApiFetch

	timeStr := time.Unix(stock.LastApiFetch/1000, 0).Format(time.RFC822Z)

	fmt.Printf("Refreshed stock at %s\n", timeStr)

	foundItems, err := stock.CheckStock(wantedItems, gag.LastRestockTimes{
		Seeds:     lastSeedRestockTime,
		Gears:     lastGearRestockTime,
		Eggs:      lastEggRestockTime,
		Merchants: lastMerchantRestockTime,
		Cosmetics: lastCosmeticRestockTime,
		Event:     lastEventRestockTime,
	}, firstRun)

	if err != nil {
		return err
	}

	firstRun = false

	lastSeedRestockTime = stock.CategoryRefreshStatus.Seeds.LastRefresh
	lastGearRestockTime = stock.CategoryRefreshStatus.Gears.LastRefresh
	lastEggRestockTime = stock.CategoryRefreshStatus.Eggs.LastRefresh
	lastMerchantRestockTime = stock.CategoryRefreshStatus.Merchants.LastRefresh
	lastCosmeticRestockTime = stock.CategoryRefreshStatus.Cosmetics.LastRefresh
	lastEventRestockTime = stock.CategoryRefreshStatus.Event.LastRefresh

	if len(foundItems) > 0 {
		sb := strings.Builder{}
		sb.WriteString("The following items are in stock:\n\n")

		for _, item := range foundItems {
			sb.WriteString(fmt.Sprintf("%s\n", item))
		}

		sb.WriteString(fmt.Sprintf("Updated at %s\n", timeStr))

		notify(sb.String())
	}

	return nil
}

func routine(ctx context.Context) {
	timer := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			fmt.Println("Checking stock, next check at ", time.Now().Add(1*time.Minute).Format(time.RFC822Z))
			err := checkStock(items)
			if err != nil {
				log.Printf("Failed to check stock: %v", err)
			}
		}
	}
}

func init() {
	beeep.AppName = "Grow a Garden Stock Checker"
}

func main() {
	fmt.Println("Starting Grow a Garden stock checker")

	err := checkStock(items)
	if err != nil {
		log.Fatalf("Failed to check stock: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go routine(ctx)

	select {}
}
