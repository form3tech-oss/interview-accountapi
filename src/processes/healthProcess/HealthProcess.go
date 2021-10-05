package healthProcess

import (
	"fmt"

	"../../../src/endpoint"
)

func GetHealth() {
	fmt.Println("\nGetting Health information:")
	endpoint.GetHealth()
}
