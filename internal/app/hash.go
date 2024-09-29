package app

import "golang.org/x/crypto/bcrypt"

func hashpassword(pasword string) (res []byte, err error) {
	arr := []byte(pasword)
	res, err = bcrypt.GenerateFromPassword(arr, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return res, err
}
