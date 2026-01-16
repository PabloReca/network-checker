package backend

type NIC struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}

type Device struct {
	Name string `json:"name"`
	NICs []NIC  `json:"nics"`
}

type Config struct {
	Devices []Device `json:"devices"`
}
