package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/xxandev/vrm"
)

type User struct {
	Name string `json:"username"`
	Pass string `json:"password"`
}

func (u *User) GetName() string { return u.Name }
func (u *User) GetPass() string { return u.Pass }

var user User

func init() {
	flag.StringVar(&user.Name, "n", "", "vrm user name")
	flag.StringVar(&user.Pass, "p", "", "vrm user password")
	flag.Parse()
}

func main() {
	client := vrm.New()
	if err := client.SetUser(user.Name, user.Pass); err != nil {
		log.Fatalf("error set user: %v\n", err)
	}
	// ========= OR
	// ub, _ := json.MarshalIndent(user, "", "\t")
	// if err := client.SetUserJson(ub); err != nil {
	// 	log.Fatalf("error set access: %v", err)
	// }
	// ========= OR
	// if err := client.SetUserInterface(&user); err != nil {
	// 	log.Fatalf("error set access: %v", err)
	// }

	if err := client.Connect(); err != nil {
		log.Fatalf("error client connect: %v\n", err)
	}
	fmt.Println(client.GetLogonJson())

	var accessList vrm.AccessTokensList
	if err := client.GetAccessTokensList(&accessList); err != nil {
		log.Printf("error get object: %v\n", err)
	} else {
		fmt.Println(formation(accessList))
	}

	var installations vrm.Installations
	if err := client.GetInstallations(&installations); err != nil {
		log.Printf("error get object: %v\n", err)
	} else {
		fmt.Println(formation(installations))
	}

	for n := range installations.Records {
		var Stats interface{}
		if err := client.GetObject(&Stats, installations.Records[n].IDSite, "stats", "type=venus"); err != nil {
			log.Printf("error get object: %v\n", err)
		} else {
			fmt.Println(formation(Stats))
		}
	}

	if err := client.Close(); err != nil {
		log.Fatalf("error client close: %v\n", err)
	}
	fmt.Println("client close, app exit")
}

func formation(v any) string {
	res, _ := json.MarshalIndent(v, "", "\t")
	return string(res)
}
