package pokeapi

import (
	"fmt"
	"testing"
)

func TestGetLocationAreasList(t *testing.T) {

	var pokeapi_URL = "https://pokeapi.co/api/v2/location-area/"
	res, err := GetLocationAreasList(pokeapi_URL)
	if err != nil {
		t.Errorf("error in getlocationareas %v", err)
	}
	fmt.Println(res)
}
