package config

var address string = ":8080"

func GetConfigValue(key string) string {
	switch key {
	case "address":
		return address
	}
	return ""
}
