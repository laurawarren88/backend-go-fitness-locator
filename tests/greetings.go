package greetings

func Hello(name string) string {
	if name == "" {
		return "Hello, World"
	}
	return "Hello, " + name
}
