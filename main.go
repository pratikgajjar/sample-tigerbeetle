package main

import (
	"fmt"
	"log"
	"reflect"

	. "github.com/tigerbeetle/tigerbeetle-go"
	. "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

// Since we only require Go 1.17 we can't do this as a generic function
// even though that would be fine. So do the dynamic approach for now.
func assert(a, b any, field string) {
	if !reflect.DeepEqual(a, b) {
		log.Fatalf("Expected %s to be [%+v (%T)], got: [%+v (%T)]", field, b, b, a, a)
	}
}

func main() {
	fmt.Println("Import ok!")
	client, err := NewClient(ToUint128(0), []string{"127.0.0.1:3000", "127.0.0.1:3001", "127.0.0.1:3002"})
	if err != nil {
		log.Printf("Error creating client: %s", err)
		return
	}

	// account0 := Account{
	// 	ID:     ToUint128(100),
	// 	Ledger: 1,
	// 	Code:   718,
	// 	Flags: AccountFlags{
	// 		DebitsMustNotExceedCredits: true,
	// 		Linked:                     true,
	// 	}.ToUint16(),
	// }
	// account1 := Account{
	// 	ID:     ToUint128(101),
	// 	Ledger: 1,
	// 	Code:   718,
	// 	Flags: AccountFlags{
	// 		History: true,
	// 	}.ToUint16(),
	// }
	// accountErrors, err := client.CreateAccounts([]Account{account0, account1})
	// log.Println(accountErrors)

	curTid := ID()
	log.Println(curTid)
	// Start a pending transfer
	transferRes, err := client.CreateTransfers([]Transfer{
		{
			ID:              curTid,
			DebitAccountID:  ToUint128(101),
			CreditAccountID: ToUint128(100),
			Amount:          ToUint128(500),
			Ledger:          1,
			Code:            1,
			Flags:           TransferFlags{Pending: false}.ToUint16(),
		},
	})
	if err != nil {
		log.Fatalf("Error creating transfer: %s", err)
	}

	for _, err := range transferRes {
		log.Fatalf("Error creating transfer: %s", err.Result)
	}
	log.Println(transferRes, "<-transferRes")

	// Validate accounts pending and posted debits/credits before finishing the two-phase transfer
	accounts, err := client.LookupAccounts([]Uint128{ToUint128(1), ToUint128(2)})
	if err != nil {
		log.Fatalf("Could not fetch accounts: %s", err)
	}
	assert(len(accounts), 2, "accounts")

	// Create a second transfer simply posting the first transfer
	transferRes, err = client.CreateTransfers([]Transfer{
		{
			ID:              ID(),
			DebitAccountID:  ToUint128(1),
			CreditAccountID: ToUint128(2),
			Amount:          ToUint128(500),
			PendingID:       curTid,
			Ledger:          1,
			Code:            1,
			Flags:           TransferFlags{PostPendingTransfer: true}.ToUint16(),
		},
	})
	log.Println(transferRes, "post")
	if err != nil {
		log.Fatalf("Error creating transfers: %s", err)
	}

	defer client.Close()
}
