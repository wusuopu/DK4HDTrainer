package utils

func Try (f func()) error {
	var err error
	func ()  {
		defer func ()  {
			if r := recover(); r!= nil {
				err = r.(error)
			}
		}()

		f()
	}()
	return err
}

func ThrowIfError (e error) {
	if e != nil {
		panic(e)
	}
}
