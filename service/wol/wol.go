package wol

type BudafokiWol struct {
	mac string
}

func NewBudafokiWol(mac string) *BudafokiWol {
	return &BudafokiWol{
		mac: mac,
	}
}

func (w *BudafokiWol) Wake() error {
	return nil
}
