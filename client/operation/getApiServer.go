package operation

import "math/rand"

func getApiServer() string {
	serverSlice := []string{"192.168.74.98:10000"}
	return serverSlice[rand.Intn(len(serverSlice))]
}