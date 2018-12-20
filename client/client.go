package client

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/aki237/nscjar"
	"golang.org/x/net/publicsuffix"
)

type lyndaClient struct {
	net *http.Client

	cookiePath string
}

var client *lyndaClient

func Init(cookiePath string) {
	client = &lyndaClient{cookiePath: cookiePath, net: &http.Client{}}
}

func GetInstance() *lyndaClient {

	if client == nil {
		return &lyndaClient{net: &http.Client{}}
	}

	return client
	// return &http.Client{Jar: bakeCookies(cookieFile)}
}

func (c *lyndaClient) Get(url string) (*http.Response, error) {
	return c.net.Get(url)
}

func (c *lyndaClient) GetVideoJsonData(url string) ([]byte, error) {
	return exec.Command("curl", "-L", url, "-b", c.cookiePath).Output()
}

func bakeCookies(cookieFile string) http.CookieJar {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(cookieFile)
	if err != nil {
		panic(err)
	}

	nscjar := nscjar.Parser{}
	cookies, err := nscjar.Unmarshal(f)
	if err != nil {
		panic(err)
	}

	var wwwLinda []*http.Cookie
	var dotLinda []*http.Cookie
	for _, cookie := range cookies {
		if strings.Contains(cookie.Domain, "www.lynda") {
			wwwLinda = append(wwwLinda, cookie)
		} else if strings.Contains(cookie.Domain, "lynda") {
			dotLinda = append(dotLinda, cookie)
		}
	}

	// fmt.Printf("url: %s == cookie: %+v\n", url, cookie)
	// fmt.Printf("%+q\n", wwwLinda)
	// fmt.Printf("%+q\n", dotLinda)

	wwwu, err := url.Parse("https://www.lynda.com")
	jar.SetCookies(wwwu, wwwLinda)

	// dotu, err := url.Parse("https://lynda.com")
	jar.SetCookies(wwwu, append(jar.Cookies(wwwu), dotLinda...))

	return jar
}
