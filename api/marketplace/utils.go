package marketplace

func toPublicGraphqlError(actualError error, allowedErrors ...string) (*string, error) {

	if actualError == nil {
		return nil, nil
	}

	for _, allowedError := range allowedErrors {
		if allowedError == actualError.Error() {
			e := actualError.Error()
			return &e, nil
		}
	}

	return nil, actualError

}
