package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Helper is the helper struct
type Helper struct {
	StringRequirement string
	AssumeRole        string
	ExternalID        string
}

// GetEnv Gets an environment variable with a default value if not present
func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

// GetSTSToken is a helper function to get data from AWS STS api
func (helper Helper) GetSTSToken(token string, role string, duration int64, externalID string) (STSSession, error) {
	stsSession := STSSession{}
	if token == "" {
		return stsSession, errors.New("No Token In Request")
	}

	if role == "" {
		role = helper.AssumeRole
	}

	if externalID == "" {
		externalID = helper.ExternalID
	}

	userData, err := getUserData(token, helper.StringRequirement)
	if err != nil {
		return stsSession, err
	}
	fmt.Println("USER => " + userData.Email + " is Making STS Request." + " is Making STS Request for Role => " + role + " with Duration => " + strconv.FormatInt(duration, 10))

	sess, err := session.NewSession()

	if err != nil {
		return stsSession, err
	}

	creds := stscreds.NewCredentials(sess, role, func(p *stscreds.AssumeRoleProvider) {
		p.RoleSessionName = userData.Email
		p.Duration = time.Duration(duration) * time.Minute
		p.ExternalID = &externalID
	})

	credentials, err := creds.Get()
	if err != nil {
		return stsSession, err
	}
	expiresAt, err := creds.ExpiresAt()
	if err != nil {
		return stsSession, err
	}

	if credentials.AccessKeyID == "" {
		return stsSession, errors.New("Server Not Able to Create Keys")
	}

	stsSession = STSSession{
		Creds:     credentials,
		ExpiresAt: expiresAt,
	}

	return stsSession, nil

}

func getUserData(token string, stringRequirement string) (IDData, error) {
	// make sure we can validate the token
	// get the users email who just refreshed
	userData := IDData{}
	uri := "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + token
	userData, err := callForUserData(uri, token)
	if err != nil {
		return userData, err
	}

	if userData.Email == "" {
		uri = "https://www.googleapis.com/oauth2/v3/tokeninfo?access_token=" + token
		userData, err = callForUserData(uri, token)
		if err != nil {
			return userData, err
		}
	}

	if strings.Contains(userData.Email, stringRequirement) == false {
		return userData, errors.New("Invalid Request for STS Key. Email => " + userData.Email)
	}

	return userData, nil
}

func callForUserData(uri string, token string) (IDData, error) {
	userData := IDData{}
	resp, err := http.Get(uri)
	if err != nil {
		return userData, err
	}
	defer resp.Body.Close()
	userBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return userData, err
	}
	err = json.Unmarshal(userBody, &userData)
	if err != nil {
		return userData, err
	}

	return userData, nil
}
