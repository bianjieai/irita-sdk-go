package integration_test

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdk "github.com/bianjieai/irita-sdk-go"
	"github.com/bianjieai/irita-sdk-go/types"
	"github.com/bianjieai/irita-sdk-go/types/store"
)

const (
	nodeURI = "tcp://localhost:26657"
	chainID = "test"
	// mode    = types.Commit
	// fee     = "4point"
	// gas     = 200000
	// algo    = "sm2"
	// level   = "info"
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	addr    = "iaa1rgnu8grzt6mwnjg7jss7w0sfyjn67g4et0hzfz"
)

var (
	path = os.ExpandEnv("$HOME/.iritacli")
)

type IntegrationTestSuite struct {
	suite.Suite
	sdk.IRITAClient
	r            *rand.Rand
	rootAccount  MockAccount
	randAccounts []MockAccount
}

type SubTest struct {
	testName string
	testCase func(s IntegrationTestSuite)
}

// MockAccount define a account for test
type MockAccount struct {
	Name, Password string
	Address        types.AccAddress
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	options := []types.Option{
		types.KeyDAOOption(store.NewMemory(nil)),
		types.TimeoutOption(10),
	}
	cfg, err := types.NewClientConfig(nodeURI, chainID, options...)
	if err != nil {
		panic(err)
	}

	s.IRITAClient = sdk.NewIRITAClient(cfg)
	s.r = rand.New(rand.NewSource(time.Now().UnixNano()))
	s.rootAccount = MockAccount{
		Name:     "v1",
		Password: "YQVGsOjegu",
		Address:  types.MustAccAddressFromBech32(addr),
	}
	s.initAccount()
}

func (s *IntegrationTestSuite) TearDownSuite() {
	_ = os.Remove(path)
}

func (s *IntegrationTestSuite) initAccount() {
	_, err := s.Key.Import(s.Account().Name,
		s.Account().Password,
		string(getPrivKeyArmor()))
	if err != nil {
		panic(err)
	}

	//var receipts bank.Receipts
	for i := 0; i < 5; i++ {
		name := s.RandStringOfLength(10)
		pwd := s.RandStringOfLength(16)
		address, _, err := s.Key.Add(name, "11111111")
		if err != nil {
			panic("generate test account failed")
		}

		s.randAccounts = append(s.randAccounts, MockAccount{
			Name:     name,
			Password: pwd,
			Address:  types.MustAccAddressFromBech32(address),
		})
	}
}

// RandStringOfLength return a random string
func (s *IntegrationTestSuite) RandStringOfLength(l int) string {
	var result []byte
	bytes := []byte(charset)
	for i := 0; i < l; i++ {
		result = append(result, bytes[s.r.Intn(len(bytes))])
	}
	return string(result)
}

// GetRandAccount return a random test account
func (s *IntegrationTestSuite) GetRandAccount() MockAccount {
	return s.randAccounts[s.r.Intn(len(s.randAccounts))]
}

// Account return a test account
func (s *IntegrationTestSuite) Account() MockAccount {
	return s.rootAccount
}

func getPrivKeyArmor() []byte {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	path = filepath.Dir(path)
	path = filepath.Join(path, "integration_test/scripts/priv.key")
	bz, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bz
}
