package tests

type Expect func(any) bool

func ExpectNil(r any) bool {
	return r != nil
}

func ExpectNotNil(r any) bool {
	return r == nil
}
