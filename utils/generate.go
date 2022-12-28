package utils

import gonanoid "github.com/matoous/go-nanoid"

// GenRoomId generate a random room id which length is 10
func GenRoomId() string {
	return gonanoid.MustGenerate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 9)
}
