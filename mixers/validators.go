package mixers

// Static token validator
func StaticTokenValidator(validToken string) TokenValidator {
	return func(token string) (any, bool) {
		if token == validToken {
			return "system_admin", true
		}
		return nil, false
	}
}
