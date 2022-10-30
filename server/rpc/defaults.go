package rpc

type Defaults struct {
	Name string

	listenAddr string
}

func getDefaults() Defaults {
	d := Defaults{
		Name:       "name",
		listenAddr: "0.0.0.0:8000",
	}

	return d
}
