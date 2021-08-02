package ara

import (
	"fmt"
)

//go:generate mockery --name=AraService

type AraService interface {
	GetRecordList(filter string) (*RecordList, error)
	GetRecord(recordId int) (*Record, error)
}

type araService struct {
	araAddr string
}

func NewAraService(araAddr string) AraService {
	return &araService{araAddr: araAddr}
}

func (a *araService) composeQuery(handler, filter string) string {
	return fmt.Sprintf("http://%s/api/v1/%s?%s", a.araAddr, handler, filter)
}
