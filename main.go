package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func main() {
	baseUrlFlag := flag.String("url", "https://www.pathofexile.com", "base API URL")
	flag.Parse()

	acc := flag.Arg(0)
	if acc == "" {
		fmt.Println("account name is required")
	}

	poeURL := *baseUrlFlag
	u, err := url.Parse(poeURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	u.Path = "character-window/get-stash-items"

	v := u.Query()
	v.Set("accountName", acc)
	v.Set("realm", "pc")
	v.Set("league", "Expedition")

	u.RawQuery = v.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	stashResp := stashResp{}
	err = json.NewDecoder(resp.Body).Decode(&stashResp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%d tabs\n", len(stashResp.Tabs))
	for _, t := range stashResp.Tabs {
		if t.Selected {
			fmt.Printf("Current tab: %s\n", t.N)
			fmt.Println("Items:")
			for _, i := range stashResp.Items {
				s := "  "
				if i.Name != "" {
					s += i.Name + ", "
				}
				s += i.TypeLine
				fmt.Println(s)
			}
		}
	}
	fmt.Printf("%d items\n", len(stashResp.Tabs))
}

type stashResp struct {
	Tabs []struct {
		N        string `json:"n"`
		Type     string `json:"type"`
		Selected bool   `json:"selected"`
	} `json:"tabs"`
	Items []struct {
		Name     string `json:"name"`
		TypeLine string `json:"typeLine"`
	} `json:"items"`
}
