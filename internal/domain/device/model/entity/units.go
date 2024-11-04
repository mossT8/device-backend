package entity

type Units struct {
	ID mysqlRecordId `json:"id"`

	Name   mysqlText `json:"name"`
	Symbol mysqlText `json:"symbol"`
}

func NewUnits(name, symbol string) Units {
	return Units{
		Name:   mysqlText(name),
		Symbol: mysqlText(symbol),
	}
}

func (u *Units) GetID() int64 {
	return int64(u.ID)
}

func (u *Units) GetName() string {
	return string(u.Name)
}

func (u *Units) GetSymbol() string {
	return string(u.Symbol)
}

func (u *Units) SetID(id int64) {
	u.ID = mysqlRecordId(id)
}

func (u *Units) SetName(name string) {
	u.Name = mysqlText(name)
}

func (u *Units) SetSymbol(symbol string) {
	u.Symbol = mysqlText(symbol)
}
