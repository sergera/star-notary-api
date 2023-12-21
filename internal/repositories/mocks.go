package repositories

import "github.com/sergera/star-notary-backend/internal/domain"

type MockDBConnection struct {
	OpenCalled  bool
	CloseCalled bool
}

func (m *MockDBConnection) Open() {
	m.OpenCalled = true
}

func (m *MockDBConnection) Close() {
	m.CloseCalled = true
}

type MockStarRepository struct {
	CreateStarCalled     bool
	SetPriceCalled       bool
	RemoveFromSaleCalled bool
	PurchaseCalled       bool
	SetNameCalled        bool
	GetStarRangeCalled   bool
}

func (m *MockStarRepository) CreateStar(star domain.StarModel) error {
	m.CreateStarCalled = true
	return nil
}

func (m *MockStarRepository) SetPrice(star domain.StarModel) error {
	m.SetPriceCalled = true
	return nil
}

func (m *MockStarRepository) RemoveFromSale(star domain.StarModel) error {
	m.RemoveFromSaleCalled = true
	return nil
}

func (m *MockStarRepository) Purchase(star domain.StarModel) error {
	m.PurchaseCalled = true
	return nil
}

func (m *MockStarRepository) SetName(star domain.StarModel) error {
	m.SetNameCalled = true
	return nil
}

func (m *MockStarRepository) GetStarRange(rangeModel domain.StarRangeModel) ([]domain.StarModel, error) {
	m.GetStarRangeCalled = true
	return []domain.StarModel{}, nil
}

type MockWalletRepository struct {
	CreateWalletCalled bool
}

func (m *MockWalletRepository) CreateWallet(wallet *domain.WalletModel) error {
	m.CreateWalletCalled = true
	return nil
}
