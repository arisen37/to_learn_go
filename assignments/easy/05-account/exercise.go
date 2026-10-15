package lab

type Account struct{ cents int64 }

func NewAccount(cents int64) (*Account, error) { panic("TODO: NewAccount") }
func (a *Account) Deposit(cents int64) error   { panic("TODO: Deposit") }
func (a *Account) Withdraw(cents int64) error  { panic("TODO: Withdraw") }
func (a *Account) Balance() int64              { panic("TODO: Balance") }
