package user_contracts

import fc "common/contracts/file-contracts"

type Response struct {
	User  *User    `json:"user"`
	File  *fc.File `json:"file"`
	Error *string  `json:"error"`
}
