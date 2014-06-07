#A Go library using Bing Translator API

example code:

```
package main

import "github.com/meoow/bingtranslate"
import "fmt"

const (
	client_id     = "your_id"
	client_secret = "your_secret"
)

func main() {
	authurl   := bingtranslate.MakeAuthURL(client_id, client_secret)
	token, _  := bingtranslate.GetToken(authurl)
	result, _ := bingtranslate.Translate(token, "Hello World", "en", "zh-CHS")
	
	fmt.Println(bingtranslate.ParseResult(result))
}
```
> 世界你好