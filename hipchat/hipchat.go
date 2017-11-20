package hipchat

import(
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"net/http"

	"github.com/taylorzr/hibye/root"
)

const Yellow = "yellow"
const Green  = "green"
const Red    = "red"

const room = "hibye"

var token, hipchatTokenExists = os.LookupEnv("HIPCHAT_TOKEN")

func init() {
	if !hipchatTokenExists {
		log.Fatal("HIPCHAT_TOKEN environment variable not set")
	}
}

type UserResponse struct {
	Items []root.User `json:"items"`
}

func SendMessage(message string, color string) error {
	log.Println(message)

	client := new(http.Client)

	url := fmt.Sprintf("https://api.hipchat.com/v2/room/%s/message?auth_token=%s", room, token)

	json, err := json.Marshal(map[string]string{ "message": message, "color": color })

	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(json))

	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)

	if err != nil {
		return err
	}

	if response.StatusCode != 201 {
		message := fmt.Sprintf("Sending message: Expected 204 response, got %d", response.StatusCode)
		return errors.New(message)
	}

	return nil
}

func GetAllUsers() ([]root.User, error) {
	httpBody, err := getUsers()

	if err != nil {
		return nil, err
	}

	users, err := unmarshalUsers(httpBody)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func getUsers() ([]byte, error) {
	client := new(http.Client)

	url := fmt.Sprintf("https://api.hipchat.com/v2/user?max-results=1000&auth_token=%s", token)

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		message := fmt.Sprintf("Getting users: Expected 200 response, got %d", response.StatusCode)
		return nil, errors.New(message)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func unmarshalUsers(httpBody []byte) ([]root.User, error) {
	userResponse := UserResponse{}

	err := json.Unmarshal(httpBody, &userResponse)

	if err != nil {
		return nil, err
	}

	return userResponse.Items, nil
}
