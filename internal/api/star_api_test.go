package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sergera/star-notary-backend/internal/notifier"
	"github.com/sergera/star-notary-backend/internal/repositories"
)

func TestCreateStar(t *testing.T) {
	mockStarRepo := &repositories.MockStarRepository{}
	mockWalletRepo := &repositories.MockWalletRepository{}
	mockNotifier := &notifier.MockStarNotifier{}
	mockConn := &repositories.MockDBConnection{}
	api := StarAPI{mockConn, mockStarRepo, mockWalletRepo, mockNotifier}

	body := `{
		"token_id": "123", 
		"coordinates": "103015.00+253015.00", 
		"name": "star name", 
		"date": "2023-01-01T00:00:00Z", 
		"owner": "0xa0Ee7A142d267C1f36714E4a8F75612F20a79720"
	}`
	req := httptest.NewRequest(http.MethodPost, "/createstar", strings.NewReader(body))
	rr := httptest.NewRecorder()

	api.CreateStar(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}

	if !mockWalletRepo.CreateWalletCalled {
		t.Error("Expected CreateWallet to be called, but it wasn't")
	}

	if !mockStarRepo.CreateStarCalled {
		t.Error("Expected CreateStar to be called, but it wasn't")
	}
}

func TestSetPrice(t *testing.T) {
	mockStarRepo := &repositories.MockStarRepository{}
	mockNotifier := &notifier.MockStarNotifier{}
	api := StarAPI{starRepo: mockStarRepo, notifier: mockNotifier}

	body := `{"token_id": "123", "price": "100"}`
	req := httptest.NewRequest(http.MethodPost, "/setprice", strings.NewReader(body))
	rr := httptest.NewRecorder()

	api.SetPrice(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
		respBody, _ := ioutil.ReadAll(rr.Body)
		t.Log("Response body:", string(respBody))
	}

	if !mockStarRepo.SetPriceCalled {
		t.Error("Expected SetPrice to be called, but it wasn't")
	}

	if !mockNotifier.PublishCalled {
		t.Error("Expected Publish to be called, but it wasn't")
	}
}

func TestRemoveFromSale(t *testing.T) {
	mockStarRepo := &repositories.MockStarRepository{}
	mockNotifier := &notifier.MockStarNotifier{}
	api := StarAPI{starRepo: mockStarRepo, notifier: mockNotifier}

	body := `{"token_id": "123"}`
	req := httptest.NewRequest(http.MethodPost, "/removefromsale", strings.NewReader(body))
	rr := httptest.NewRecorder()

	api.RemoveFromSale(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}

	if !mockStarRepo.RemoveFromSaleCalled {
		t.Error("Expected RemoveFromSale to be called, but it wasn't")
	}

	if !mockNotifier.PublishCalled {
		t.Error("Expected Publish to be called, but it wasn't")
	}
}

func TestPurchase(t *testing.T) {
	mockStarRepo := &repositories.MockStarRepository{}
	mockWalletRepo := &repositories.MockWalletRepository{}
	mockNotifier := &notifier.MockStarNotifier{}
	api := StarAPI{starRepo: mockStarRepo, walletRepo: mockWalletRepo, notifier: mockNotifier}

	body := `{
		"token_id": "123",
		"owner": "0xa0Ee7A142d267C1f36714E4a8F75612F20a79720",
		"price": "100",
		"date": "2023-01-01T00:00:00Z"
	}`
	req := httptest.NewRequest(http.MethodPost, "/purchase", strings.NewReader(body))
	rr := httptest.NewRecorder()

	api.Purchase(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}

	if !mockWalletRepo.CreateWalletCalled {
		t.Error("Expected CreateWallet to be called, but it wasn't")
	}

	if !mockStarRepo.PurchaseCalled {
		t.Error("Expected Purchase to be called, but it wasn't")
	}

	if !mockNotifier.PublishCalled {
		t.Error("Expected Publish to be called, but it wasn't")
	}
}

func TestSetName(t *testing.T) {
	mockStarRepo := &repositories.MockStarRepository{}
	mockNotifier := &notifier.MockStarNotifier{}
	api := StarAPI{starRepo: mockStarRepo, notifier: mockNotifier}

	body := `{"token_id": "123", "name": "new star name", "date": "2023-01-01T00:00:00Z"}`
	req := httptest.NewRequest(http.MethodPost, "/setname", strings.NewReader(body))
	rr := httptest.NewRecorder()

	api.SetName(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}

	if !mockStarRepo.SetNameCalled {
		t.Error("Expected SetName to be called, but it wasn't")
	}

	if !mockNotifier.PublishCalled {
		t.Error("Expected Publish to be called, but it wasn't")
	}
}

func TestGetStarRange(t *testing.T) {
	mockStarRepo := &repositories.MockStarRepository{}
	api := StarAPI{starRepo: mockStarRepo}

	req := httptest.NewRequest(http.MethodGet, "/getstarrange?start=1&end=5&oldest-first=false", nil)
	rr := httptest.NewRecorder()

	api.GetStarRange(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}

	if !mockStarRepo.GetStarRangeCalled {
		t.Error("Expected GetStarRange to be called, but it wasn't")
	}
}
