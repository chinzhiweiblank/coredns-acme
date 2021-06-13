package acme

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/libdns/libdns"
)

type Provider struct {
	recordMap map[string][]libdns.Record
	m         sync.Mutex
}

var provider = Provider{
	recordMap: make(map[string][]libdns.Record),
}

func (p *Provider) AppendRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	p.m.Lock()
	defer p.m.Unlock()
	records := []libdns.Record{}
	if _, found := p.recordMap[zone]; found {
		records = p.recordMap[zone]
	}
	records = append(records, recs...)
	p.recordMap[zone] = records
	return records, nil
}

func (p *Provider) DeleteRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	p.m.Lock()
	defer p.m.Unlock()
	deletedRecords := []libdns.Record{}
	remainingRecords := []libdns.Record{}
	records, found := p.recordMap[zone]
	if !found {
		return deletedRecords, nil
	}
	shouldDelete := make([]bool, len(records))
	for _, deleteRecord := range recs {
		for i, record := range records {
			if reflect.DeepEqual(deleteRecord, record) {
				shouldDelete[i] = true
			}
		}
	}
	for i, record := range records {
		if !shouldDelete[i] {
			remainingRecords = append(remainingRecords, record)
		} else {
			deletedRecords = append(deletedRecords, record)
		}
	}
	p.recordMap[zone] = remainingRecords
	return deletedRecords, nil
}

func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	p.m.Lock()
	defer p.m.Unlock()

	records, found := p.recordMap[zone]
	if !found {
		return nil, fmt.Errorf("records for zone %s not found", zone)
	}

	return records, nil
}
