## Installation
``
go get github.com/sysfa/gocurl
``
### Information
Use URL with scheme
This library using Net Package
if you using POST method you need to add Content-Length Manually
### TO DO
add Concurrency
#### Usage
```
func main() {
	c := gocurl.New()
	c.setHeader("Content-Type", "application/json")
	z := map[string]string{"user": "duar"}
	c.Curl("POST", "http://localhost:80/user", z)
	a, err := c.Exec()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a.Status)
	fmt.Println(a.Body)
	fmt.Println(a.Header)
}
```
