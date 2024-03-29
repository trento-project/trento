// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	web "github.com/trento-project/trento/web"
)

// TrentoApiService is an autogenerated mock type for the TrentoApiService type
type TrentoApiService struct {
	mock.Mock
}

// GetClustersSettings provides a mock function with given fields:
func (_m *TrentoApiService) GetClustersSettings() (web.ClustersSettingsResponse, error) {
	ret := _m.Called()

	var r0 web.ClustersSettingsResponse
	if rf, ok := ret.Get(0).(func() web.ClustersSettingsResponse); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(web.ClustersSettingsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsWebServerUp provides a mock function with given fields:
func (_m *TrentoApiService) IsWebServerUp() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
