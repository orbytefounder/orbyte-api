package internal
import (
	"errors"
	"sync"
)
type Wallet struct {
	mu sync.Mutex
	Balance float64
}
func NewWallet(initialDeposit float64) *Wallet {
	return &Wallet{Balance: initialDeposit}
}
func (w *Wallet) DeductServerBill(amount float64) (float64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.Balance < amount {
		return 0, errors.New("insufficient wallet balance for server usage")
	}
	platformFee := amount * 0.02
	w.Balance -= platformFee
	return w.Balance, nil
}