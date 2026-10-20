package main

import (
	"fmt"
	"sync"
)

type Account struct {
	id      int
	balance int
	mu      sync.Mutex
}

// ❌ 錯誤做法（會死鎖）：如果同時 A 轉給 B，且 B 轉給 A，雙方鎖定順序交叉就會卡死。
// 💡 正確做法：動態調整鎖的順序，確保任何時候都是「小的 ID 先被 Lock」
func transferMoney(from *Account, to *Account, amount int) {
	if from == to {
		return
	} // 同一帳戶只需視為無操作，避免重複鎖定自己。
	firstLock := from
	secondLock := to
	if from.id > to.id {
		firstLock = to
		secondLock = from
	}

	firstLock.mu.Lock()
	defer firstLock.mu.Unlock()

	secondLock.mu.Lock()
	defer secondLock.mu.Unlock()

	from.balance -= amount
	to.balance += amount
	fmt.Printf("成功從帳戶 %d 轉帳給 帳戶 %d 共 $%d\n", from.id, to.id, amount)
}

func main() {
	acc1 := &Account{id: 1, balance: 1000}
	acc2 := &Account{id: 2, balance: 1000}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		transferMoney(acc1, acc2, 100) // 鎖 1 -> 鎖 2
	}()

	go func() {
		defer wg.Done()
		transferMoney(acc2, acc1, 50) // 鎖 1 -> 鎖 2 (排序後依然是 1 先 2 後，不會死鎖)
	}()

	wg.Wait()
	transferMoney(acc1, acc1, 100)
	fmt.Printf("餘額: %d %d；總額: %d\n", acc1.balance, acc2.balance, acc1.balance+acc2.balance)
}
