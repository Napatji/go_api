package account

import "strconv"

var AccountList []Account

type Account struct {
	Account_id string `json:"account_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Balance    int64  `json:"balance"`
}

func (a *Account) SetAccount(accId, name, email string, balance int64) {
	a.Account_id = accId
	a.Name = name
	a.Email = email
	a.Balance = balance
}

func GetAccountList() []Account {
	return AccountList
}

func GetNextId(acl []Account) string {
	s := 0
	for i, acc := range acl {
		if i == len(acl)-1 {
			s, _ = strconv.Atoi(acc.Account_id)
		}
	}
	return strconv.Itoa(s + 1)
}

func CheckId(acl []Account, id string) (hasId bool, index int) {
	hasId = false
	index = -99
	for i, a := range acl {
		if a.Account_id == id {
			hasId = true
			index = i
		}
	}
	return hasId, index
}
