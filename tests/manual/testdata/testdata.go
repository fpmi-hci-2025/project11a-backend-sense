package testdata

import (
	"encoding/json"
	"os"
	"sync"
)

// TestData stores test data
type TestData struct {
	BaseURL string `json:"base_url"`
	Users   struct {
		User1 struct {
			Username    string `json:"username"`
			Email       string `json:"email"`
			Password    string `json:"password"`
			Phone       string `json:"phone"`
			Description string `json:"description"`
		} `json:"user1"`
		User2 struct {
			Username    string `json:"username"`
			Email       string `json:"email"`
			Password    string `json:"password"`
			Phone       string `json:"phone"`
			Description string `json:"description"`
		} `json:"user2"`
	} `json:"users"`
	Tokens struct {
		User1 string `json:"user1"`
		User2 string `json:"user2"`
	} `json:"tokens"`
	IDs struct {
		User1ID       string `json:"user1_id"`
		User2ID       string `json:"user2_id"`
		PublicationID string `json:"publication_id"`
		CommentID     string `json:"comment_id"`
		MediaID       string `json:"media_id"`
	} `json:"ids"`
}

var (
	data     *TestData
	dataLock sync.RWMutex
)

// LoadTestData loads test data from file
func LoadTestData(filePath string) error {
	dataLock.Lock()
	defer dataLock.Unlock()

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data = &TestData{}
	if err := json.NewDecoder(file).Decode(data); err != nil {
		return err
	}

	return nil
}

// SaveTestData saves test data to file
func SaveTestData(filePath string) error {
	dataLock.RLock()
	defer dataLock.RUnlock()

	if data == nil {
		return nil
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// GetTestData returns test data
func GetTestData() *TestData {
	dataLock.RLock()
	defer dataLock.RUnlock()
	return data
}

// SetUser1Token sets token for user1
func SetUser1Token(token string) {
	dataLock.Lock()
	defer dataLock.Unlock()
	if data != nil {
		data.Tokens.User1 = token
	}
}

// SetUser2Token sets token for user2
func SetUser2Token(token string) {
	dataLock.Lock()
	defer dataLock.Unlock()
	if data != nil {
		data.Tokens.User2 = token
	}
}

// SetUser1ID sets user1 ID
func SetUser1ID(id string) {
	dataLock.Lock()
	defer dataLock.Unlock()
	if data != nil {
		data.IDs.User1ID = id
	}
}

// SetUser2ID sets user2 ID
func SetUser2ID(id string) {
	dataLock.Lock()
	defer dataLock.Unlock()
	if data != nil {
		data.IDs.User2ID = id
	}
}

// SetPublicationID sets publication ID
func SetPublicationID(id string) {
	dataLock.Lock()
	defer dataLock.Unlock()
	if data != nil {
		data.IDs.PublicationID = id
	}
}

// SetCommentID sets comment ID
func SetCommentID(id string) {
	dataLock.Lock()
	defer dataLock.Unlock()
	if data != nil {
		data.IDs.CommentID = id
	}
}

// SetMediaID sets media ID
func SetMediaID(id string) {
	dataLock.Lock()
	defer dataLock.Unlock()
	if data != nil {
		data.IDs.MediaID = id
	}
}
