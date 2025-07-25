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

var stockManager *gag.StockManager

func stockUpdateCallback(sm *gag.StockManager) {
	foundItems, err := sm.GetWantedStock(items)
	if err != nil {
		log.Printf("Failed to get wanted stock: %v", err)
		return
	}

	log.Printf("Found %d items in stock", len(foundItems))

	notifyDesktop(foundItems, sm.GetStockTimeString())
}

func notify(message string) {
	err := beeep.Alert("Grow a Garden", message, "")
	if err != nil {
		log.Printf("Failed to notify: %v", err)
	}
}

func notifyDesktop(foundItems []string, timeStr string) {
	if len(foundItems) > 0 {
		sb := strings.Builder{}
		sb.WriteString("The following items are in stock:\n\n")

		for _, item := range foundItems {
			sb.WriteString(fmt.Sprintf("%s\n", item))
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

	server.GET("/wanted", func(w http.ResponseWriter, r *http.Request) {
		items, err := stockManager.GetWantedStock(items)
		if err != nil {
			log.Printf("Failed to get wanted stock: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Printf("Found %d items in stock", len(items))

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
		stock := stockManager.GetStock()
		images := stock.ImageData

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
