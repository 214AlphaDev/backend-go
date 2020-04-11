package wishlist

func toPublicGraphqlError(actualError error, allowedErrors ...string) (*string, error) {

	for _, allowedError := range allowedErrors {
		if allowedError == actualError.Error() {
			e := actualError.Error()
			return &e, nil
		}
	}

	return nil, actualError

}
