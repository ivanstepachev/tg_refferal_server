package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomTgID() string {
	num := 10000 + rand.Intn(90000)
	return strconv.Itoa(num)
}

func RandomStartMessage() string {
	withRefLink := fmt.Sprintf("/start %v", RandomTgID())
	startMessages := []string{withRefLink, "/start landing", "/start"}
	n := len(startMessages)
	return startMessages[rand.Intn(n)]
}