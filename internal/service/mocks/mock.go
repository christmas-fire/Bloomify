// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"

	models "github.com/christmas-fire/Bloomify/internal/models"
	service "github.com/christmas-fire/Bloomify/internal/service"
	gomock "github.com/golang/mock/gomock"
)

// MockAuth is a mock of Auth interface.
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
}

// MockAuthMockRecorder is the mock recorder for MockAuth.
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance.
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockAuth) CreateUser(user models.User) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAuthMockRecorder) CreateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAuth)(nil).CreateUser), user)
}

// GenerateToken mocks base method.
func (m *MockAuth) GenerateToken(username, password string) (service.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", username, password)
	ret0, _ := ret[0].(service.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockAuthMockRecorder) GenerateToken(username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockAuth)(nil).GenerateToken), username, password)
}

// ParseToken mocks base method.
func (m *MockAuth) ParseToken(accessToken string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", accessToken)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockAuthMockRecorder) ParseToken(accessToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockAuth)(nil).ParseToken), accessToken)
}

// RefreshToken mocks base method.
func (m *MockAuth) RefreshToken(refreshToken string) (service.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshToken", refreshToken)
	ret0, _ := ret[0].(service.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshToken indicates an expected call of RefreshToken.
func (mr *MockAuthMockRecorder) RefreshToken(refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshToken", reflect.TypeOf((*MockAuth)(nil).RefreshToken), refreshToken)
}

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockUser) Delete(userId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUserMockRecorder) Delete(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUser)(nil).Delete), userId)
}

// GetAll mocks base method.
func (m *MockUser) GetAll() ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockUserMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockUser)(nil).GetAll))
}

// GetById mocks base method.
func (m *MockUser) GetById(userId int) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", userId)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockUserMockRecorder) GetById(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockUser)(nil).GetById), userId)
}

// UpdatePassword mocks base method.
func (m *MockUser) UpdatePassword(userId int, input models.UpdatePasswordInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePassword", userId, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePassword indicates an expected call of UpdatePassword.
func (mr *MockUserMockRecorder) UpdatePassword(userId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePassword", reflect.TypeOf((*MockUser)(nil).UpdatePassword), userId, input)
}

// UpdateUsername mocks base method.
func (m *MockUser) UpdateUsername(userId int, input models.UpdateUsernameInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUsername", userId, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUsername indicates an expected call of UpdateUsername.
func (mr *MockUserMockRecorder) UpdateUsername(userId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUsername", reflect.TypeOf((*MockUser)(nil).UpdateUsername), userId, input)
}

// MockFlower is a mock of Flower interface.
type MockFlower struct {
	ctrl     *gomock.Controller
	recorder *MockFlowerMockRecorder
}

// MockFlowerMockRecorder is the mock recorder for MockFlower.
type MockFlowerMockRecorder struct {
	mock *MockFlower
}

// NewMockFlower creates a new mock instance.
func NewMockFlower(ctrl *gomock.Controller) *MockFlower {
	mock := &MockFlower{ctrl: ctrl}
	mock.recorder = &MockFlowerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFlower) EXPECT() *MockFlowerMockRecorder {
	return m.recorder
}

// CreateFlower mocks base method.
func (m *MockFlower) CreateFlower(flower models.Flower) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFlower", flower)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFlower indicates an expected call of CreateFlower.
func (mr *MockFlowerMockRecorder) CreateFlower(flower interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFlower", reflect.TypeOf((*MockFlower)(nil).CreateFlower), flower)
}

// Delete mocks base method.
func (m *MockFlower) Delete(flowerId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", flowerId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockFlowerMockRecorder) Delete(flowerId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockFlower)(nil).Delete), flowerId)
}

// GetAll mocks base method.
func (m *MockFlower) GetAll() ([]models.Flower, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]models.Flower)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockFlowerMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockFlower)(nil).GetAll))
}

// GetById mocks base method.
func (m *MockFlower) GetById(flowerId int) (models.Flower, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", flowerId)
	ret0, _ := ret[0].(models.Flower)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockFlowerMockRecorder) GetById(flowerId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockFlower)(nil).GetById), flowerId)
}

// GetFlowersByDescription mocks base method.
func (m *MockFlower) GetFlowersByDescription(description string) ([]models.Flower, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlowersByDescription", description)
	ret0, _ := ret[0].([]models.Flower)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFlowersByDescription indicates an expected call of GetFlowersByDescription.
func (mr *MockFlowerMockRecorder) GetFlowersByDescription(description interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlowersByDescription", reflect.TypeOf((*MockFlower)(nil).GetFlowersByDescription), description)
}

// GetFlowersByName mocks base method.
func (m *MockFlower) GetFlowersByName(name string) ([]models.Flower, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlowersByName", name)
	ret0, _ := ret[0].([]models.Flower)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFlowersByName indicates an expected call of GetFlowersByName.
func (mr *MockFlowerMockRecorder) GetFlowersByName(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlowersByName", reflect.TypeOf((*MockFlower)(nil).GetFlowersByName), name)
}

// GetFlowersByPrice mocks base method.
func (m *MockFlower) GetFlowersByPrice(price string) ([]models.Flower, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlowersByPrice", price)
	ret0, _ := ret[0].([]models.Flower)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFlowersByPrice indicates an expected call of GetFlowersByPrice.
func (mr *MockFlowerMockRecorder) GetFlowersByPrice(price interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlowersByPrice", reflect.TypeOf((*MockFlower)(nil).GetFlowersByPrice), price)
}

// GetFlowersByStock mocks base method.
func (m *MockFlower) GetFlowersByStock(stock string) ([]models.Flower, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlowersByStock", stock)
	ret0, _ := ret[0].([]models.Flower)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFlowersByStock indicates an expected call of GetFlowersByStock.
func (mr *MockFlowerMockRecorder) GetFlowersByStock(stock interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlowersByStock", reflect.TypeOf((*MockFlower)(nil).GetFlowersByStock), stock)
}

// UpdateDescription mocks base method.
func (m *MockFlower) UpdateDescription(flowerId int, input models.UpdateDescriptionInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDescription", flowerId, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDescription indicates an expected call of UpdateDescription.
func (mr *MockFlowerMockRecorder) UpdateDescription(flowerId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDescription", reflect.TypeOf((*MockFlower)(nil).UpdateDescription), flowerId, input)
}

// UpdateName mocks base method.
func (m *MockFlower) UpdateName(flowerId int, input models.UpdateNameInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateName", flowerId, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateName indicates an expected call of UpdateName.
func (mr *MockFlowerMockRecorder) UpdateName(flowerId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateName", reflect.TypeOf((*MockFlower)(nil).UpdateName), flowerId, input)
}

// UpdatePrice mocks base method.
func (m *MockFlower) UpdatePrice(flowerId int, input models.UpdatePriceInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePrice", flowerId, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePrice indicates an expected call of UpdatePrice.
func (mr *MockFlowerMockRecorder) UpdatePrice(flowerId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePrice", reflect.TypeOf((*MockFlower)(nil).UpdatePrice), flowerId, input)
}

// UpdateStock mocks base method.
func (m *MockFlower) UpdateStock(flowerId int, input models.UpdateStockInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStock", flowerId, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStock indicates an expected call of UpdateStock.
func (mr *MockFlowerMockRecorder) UpdateStock(flowerId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStock", reflect.TypeOf((*MockFlower)(nil).UpdateStock), flowerId, input)
}

// MockOrder is a mock of Order interface.
type MockOrder struct {
	ctrl     *gomock.Controller
	recorder *MockOrderMockRecorder
}

// MockOrderMockRecorder is the mock recorder for MockOrder.
type MockOrderMockRecorder struct {
	mock *MockOrder
}

// NewMockOrder creates a new mock instance.
func NewMockOrder(ctrl *gomock.Controller) *MockOrder {
	mock := &MockOrder{ctrl: ctrl}
	mock.recorder = &MockOrderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrder) EXPECT() *MockOrderMockRecorder {
	return m.recorder
}

// CreateOrder mocks base method.
func (m *MockOrder) CreateOrder(userId int, order_flowers models.OrderFlowers) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", userId, order_flowers)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockOrderMockRecorder) CreateOrder(userId, order_flowers interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockOrder)(nil).CreateOrder), userId, order_flowers)
}

// Delete mocks base method.
func (m *MockOrder) Delete(orderId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", orderId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockOrderMockRecorder) Delete(orderId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockOrder)(nil).Delete), orderId)
}

// GetAll mocks base method.
func (m *MockOrder) GetAll() ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockOrderMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockOrder)(nil).GetAll))
}

// GetAllOrderFlowers mocks base method.
func (m *MockOrder) GetAllOrderFlowers() ([]models.OrderFlowers, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllOrderFlowers")
	ret0, _ := ret[0].([]models.OrderFlowers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllOrderFlowers indicates an expected call of GetAllOrderFlowers.
func (mr *MockOrderMockRecorder) GetAllOrderFlowers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllOrderFlowers", reflect.TypeOf((*MockOrder)(nil).GetAllOrderFlowers))
}

// GetById mocks base method.
func (m *MockOrder) GetById(orderId int) (models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", orderId)
	ret0, _ := ret[0].(models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockOrderMockRecorder) GetById(orderId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockOrder)(nil).GetById), orderId)
}

// GetOrderFlowersByOrderId mocks base method.
func (m *MockOrder) GetOrderFlowersByOrderId(orderFlowersId int) ([]models.OrderFlowers, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderFlowersByOrderId", orderFlowersId)
	ret0, _ := ret[0].([]models.OrderFlowers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderFlowersByOrderId indicates an expected call of GetOrderFlowersByOrderId.
func (mr *MockOrderMockRecorder) GetOrderFlowersByOrderId(orderFlowersId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderFlowersByOrderId", reflect.TypeOf((*MockOrder)(nil).GetOrderFlowersByOrderId), orderFlowersId)
}

// GetOrdersByUserId mocks base method.
func (m *MockOrder) GetOrdersByUserId(userId string) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersByUserId", userId)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersByUserId indicates an expected call of GetOrdersByUserId.
func (mr *MockOrderMockRecorder) GetOrdersByUserId(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersByUserId", reflect.TypeOf((*MockOrder)(nil).GetOrdersByUserId), userId)
}

// UpdateOrder mocks base method.
func (m *MockOrder) UpdateOrder(orderId int, input models.UpdateOrderInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrder", orderId, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrder indicates an expected call of UpdateOrder.
func (mr *MockOrderMockRecorder) UpdateOrder(orderId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrder", reflect.TypeOf((*MockOrder)(nil).UpdateOrder), orderId, input)
}

// UpdateOrderFlowerId mocks base method.
func (m *MockOrder) UpdateOrderFlowerId(orderId int, input models.UpdateOrderFlowerIdInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrderFlowerId", orderId, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrderFlowerId indicates an expected call of UpdateOrderFlowerId.
func (mr *MockOrderMockRecorder) UpdateOrderFlowerId(orderId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrderFlowerId", reflect.TypeOf((*MockOrder)(nil).UpdateOrderFlowerId), orderId, input)
}

// UpdateOrderQuantity mocks base method.
func (m *MockOrder) UpdateOrderQuantity(orderId int, input models.UpdateOrderQuantityInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrderQuantity", orderId, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrderQuantity indicates an expected call of UpdateOrderQuantity.
func (mr *MockOrderMockRecorder) UpdateOrderQuantity(orderId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrderQuantity", reflect.TypeOf((*MockOrder)(nil).UpdateOrderQuantity), orderId, input)
}
