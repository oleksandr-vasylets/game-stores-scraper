package settings

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

const MaxCount = 100

type Profile struct {
	CountryCode string
	Locale      string
}

var UserProfile Profile

const userProfileFilename = "userProfile.bin"
const defaultCountryCode = "ua"
const defaultLocale = "uk"

func init() {
	UserProfile = Profile{CountryCode: defaultCountryCode, Locale: defaultLocale}
	if _, err := os.Stat(userProfileFilename); os.IsNotExist(err) {
		fmt.Println(userProfileFilename, "not found, loading fallback values")
		file, err := os.Create(userProfileFilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		file.Close()

		Save()
		return
	}
	file, err := os.Open(userProfileFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer file.Close()

	buffer := new(bytes.Buffer)
	reader := bufio.NewReader(file)
	_, err = reader.WriteTo(buffer)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	dec := gob.NewDecoder(buffer)
	err = dec.Decode(&UserProfile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// TODO: Here it is assumed that the user have not tinkered with the file (deleted, modified etc.)
}

func Save() {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(UserProfile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// 0644 means that we have permissions to read and write, but others are only permitted to read
	err = ioutil.WriteFile(userProfileFilename, buffer.Bytes(), 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
