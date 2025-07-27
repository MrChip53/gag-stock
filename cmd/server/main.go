package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gen2brain/beeep"
	"github.com/mrchip53/go-a-garden/pkg/gag"
)

var items = []string{
	// Seed Shop
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
	"Elder Strawberry",

	// Gear Shop
	"Basic Sprinkler",
	"Advanced Sprinkler",
	"Godly Sprinkler",
	"Master Sprinkler",
	"Medium Toy",
	"Medium Treat",
	"Tanning Mirror",
	"Friendship Pot",

	// Zen Shop
	"Zen Seed Pack",
	"Zen Egg",
	"Hot Spring",
	"Soft Sunshine",
	"Koi",
	"Spiked Mango",
	"Pet Shard Tranquil",
	"Pet Shard Corrupted",
	"Raiju",

	// Egg Shop
	"Common Summer Egg",
	"Rare Summer Egg",
	"Mythical Egg",
	"Paradise Egg",
	"Bug Egg",
}

var stockManager *gag.StockManager

func stockUpdateCallback(sc gag.ShopContainer) {
	foundItems := sc.GetWantedStock(items)

	notifyDesktop(foundItems, sc.GetTimeString())
}

func notify(message string) {
	err := beeep.Alert("Grow a Garden", message, "")
	if err != nil {
		log.Printf("Failed to notify: %v", err)
	}
}

func notifyDesktop(foundItems []gag.Item, timeStr string) {
	if len(foundItems) > 0 {
		sb := strings.Builder{}
		sb.WriteString("The following items are in stock:\n\n")

		for _, item := range foundItems {
			sb.WriteString(fmt.Sprintf("%s x %d\n", item.Name, item.Count))
		}

		sb.WriteString(fmt.Sprintf("Updated at %s\n", timeStr))

		notify(sb.String())
	}
}

func init() {
	beeep.AppName = "Grow a Garden Stock Checker"
}

func main() {
	log.Println("Starting Grow a Garden stock checker")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Starting stock manager")
	stockManager = gag.NewStockManager(ctx, stockUpdateCallback)

	log.Println("Starting server")
	server := newServer()

	server.GET("/last-seen", func(w http.ResponseWriter, r *http.Request) {
		sc := stockManager.GetShopContainer()
		lastSeen := sc.GetLastSeenItems()

		jsonBytes, err := json.Marshal(lastSeen)
		if err != nil {
			log.Printf("Failed to marshal last seen: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(jsonBytes)))
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	})

	server.GET("/all", func(w http.ResponseWriter, r *http.Request) {
		sc := stockManager.GetShopContainer()
		items := sc.GetAllItems()

		jsonBytes, err := json.Marshal(items)
		if err != nil {
			log.Printf("Failed to marshal items: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(jsonBytes)))
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	})

	server.GET("/wanted", func(w http.ResponseWriter, r *http.Request) {
		sc := stockManager.GetShopContainer()
		items := sc.GetWantedStock(items)

		jsonBytes, err := json.Marshal(items)
		if err != nil {
			log.Printf("Failed to marshal items: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(jsonBytes)))
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	})

	server.GET("/images", func(w http.ResponseWriter, r *http.Request) {
		images := stockManager.GetImageData()

		jsonBytes, err := json.Marshal(images)
		if err != nil {
			log.Printf("Failed to marshal images: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(jsonBytes)))
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	})

	go func() {
		log.Println("Starting server on port 8001")
		err := server.server.ListenAndServe()
		if err != nil {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	go func() {
		<-ctx.Done()
		server.server.Shutdown(ctx)
	}()

	select {}
}
