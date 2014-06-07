// another_test project main.go
package bingtranslate

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	AuthAddr      = "https://datamarket.accesscontrol.windows.net/v2/OAuth2-13"
	TransAddr     = "http://api.microsofttranslator.com/v2/Http.svc/Translate?"
	resultPattern = `^<string xmlns=".+?">(.*)</string>$`
)

var ResultRE = regexp.MustCompile(resultPattern)

/*
Supported Language Code
ar : Arabic
bg : Bulgarian
ca : Catalan
zh-CHS : Chinese (Simplified)
zh-CHT : Chinese (Traditional)
cs : Czech
da : Danish
nl : Dutch
en : English
et : Estonian
fi : Finnish
fr : French
de : German
el : Greek
ht : Haitian Creole
he : Hebrew
hi : Hindi
hu : Hungarian
id : Indonesian
it : Italian
ja : Japanese
ko : Korean
lv : Latvian
lt : Lithuanian
mww : Hmong Daw
no : Norwegian
pl : Polish
pt : Portuguese
ro : Romanian
ru : Russian
sk : Slovak
sl : Slovenian
es : Spanish
sv : Swedish
th : Thai
tr : Turkish
uk : Ukrainian
vi : Vietnamese
*/
var SupportedLangs = regexp.MustCompile(`^(ar|bg|ca|zh-CHT|zh-CHS|cs|da|nl|en|et|fi|fr|de|el|ht|he|hi|hu|id|it|ja|ko|lv|lt|mww|no|pl|pt|ro|ru|sk|sl|es|sv|th|tr|uk|vi)$`)

//go to https://datamarket.azure.com/developer/applications
//register your own client id and client secret in order to use
//bing translate API
func MakeAuthURL(client_id, client_secret string) string {
	data, _ := url.ParseQuery("")
	data.Add("client_id", client_id)
	data.Add("client_secret", client_secret)
	data.Add("grant_type", "client_credentials")
	data.Add("scope", "http://api.microsofttranslator.com")
	return data.Encode()
}

func GetToken(authurl string) (string, error) {
	auReader := strings.NewReader(authurl)
	request, _ := http.NewRequest("POST", AuthAddr, auReader)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	resultJson := make(map[string]interface{}, 4)
	err = json.Unmarshal(result, &resultJson)
	if err != nil {
		return "", err
	}

	if token, ok := resultJson["access_token"]; ok {
		return token.(string), nil
	} else {
		return "", nil
	}
	panic("")
}

func Translate(token, text, from, to string) (string, error) {
	if token == "" || text == "" {
		return "", nil
	}
	if !LangCodeCheck(from) {
		err := errors.New(from + " is not supported for Bing Translator.")
		return "", err
	}
	if !LangCodeCheck(to) {
		err := errors.New(to + " is not supported for Bing Translator.")
		return "", err
	}

	texturl, err := url.ParseQuery("")
	texturl.Add("text", text)
	texturl.Add("from", from)
	texturl.Add("to", to)

	request, _ := http.NewRequest("GET", TransAddr+texturl.Encode(), nil)
	request.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func ParseResult(result string) string {
	mat := ResultRE.FindStringSubmatch(result)
	if len(mat) == 2 {
		return mat[1]
	} else {
		return ""
	}
}

func LangCodeCheck(lang string) bool {
	if SupportedLangs.MatchString(lang) {
		return true
	} else {
		return false
	}
	panic("")
}
