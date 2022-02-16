package util

func ReformatPath(path string) string {
	if path == "" {
		path = "./"
	} else if path[len(path)-1:] != "/" {
		path += "/"
	}
	return path
}
