package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("(keybase) Usage: %s [username]\n", os.Args[0])
		return
	}

	username := os.Args[1]

	fmt.Printf("(keybase) Attemtping to get E-Mail for: \"%s\"\n", username)

	pgp, err := get_pgp(username)

	if err != nil {
		fmt.Printf("(keybase) Error: %s", err)
		return
	}

	decoded, err := decode(pgp)

	if err != nil {
		fmt.Printf("(keybase) Error: %s", err)
		return
	}

	for _, email := range get_email(decoded) {
		fmt.Printf("(keybase) E-Mail: %s\n", email[1:len(email)-1])
	}
}

func get_pgp(username string) (string, error) {
	resp, err := http.Get("https://keybase.io/" + username + "/pgp_keys.asc")

	if err != nil {
		return "", errors.New("unable to fetch pgp key")
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", errors.New("unable to read body of pgp key")
	}

	sbody := string(body)

	if strings.Contains(sbody, "SELF-SIGNED PUBLIC KEY NOT FOUND") {
		return "", errors.New("user doesn't have a self-signed key")
	}

	split := strings.Split(sbody, "\n")
	return strings.Join(split[3:len(split)-3], ""), nil
}

func decode(pgp string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(pgp)

	if err != nil {
		return "", errors.New("unable to decode pgp")
	}

	return string(decoded), nil
}

func get_email(pgp string) []string {
	return regexp.MustCompile(`<[\w+|\d+]{0,64}\@[\w+|\d+.]{1,300}>`).FindAllString(pgp, -1)
}
