package monkebase

func ReadTokenStat(token string) (owner string, valid bool, err error) {
    return
}

func CreateToken(ID string) (token string, expires int64, err error) {
    return
}

func CreateSecret(ID string) (secret string, err error) {
    return
}

func RevokeToken(token string) (err error) {
    return
}

func RevokeSecret(secret string) (err error) {
    return
}

func RevokeTokenOf(token string) (err error) {
    return
}

func RevokeSecretOf(secret string) (err error) {
	return
}

func CheckPassword(ID, password string) (ok bool, err error) {
	return
}

func SetPassword(ID, password string) (err error) {
	return
}
