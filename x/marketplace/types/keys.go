package types

var (
	ServiceKey = []byte{0x01}
)

func GetServiceKey(id string) []byte {
	return append(ServiceKey, []byte(id)...)
}
