package tools

type mockDB struct{}

var mockLoginDetails = map[string]LoginDetails{
	"tiberiu": {
		AuthToken: "QZWXEC123",
		Username:  "tiberiu",
	},
}

func (db *mockDB) GetUserLoginDetails(username string) *LoginDetails {

	clientData, ok := mockLoginDetails[username]

	if !ok {
		return nil
	}

	return &clientData
}

func (db *mockDB) SetupDatabase() error {
	return nil
}
