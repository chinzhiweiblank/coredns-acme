package acme

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/libdns/libdns"
)

type RecordStore struct {
	entries []libdns.Record
}

type Provider struct {
	sync.Mutex
	recordMap map[string]*RecordStore
	// m              sync.Mutex
}

func (p *Provider) getZoneRecords(ctx context.Context, zoneName string) *RecordStore {
	records, found := p.recordMap[zoneName]
	if !found {
		return nil
	}
	return records
}

func (p *Provider) AppendRecords(ctx context.Context, zoneName string, recs []libdns.Record) ([]libdns.Record, error) {
	p.Lock()
	defer p.Unlock()
	zoneRecordStore := p.getZoneRecords(ctx, zoneName)
	if zoneRecordStore == nil {
		zoneRecordStore = new(RecordStore)
		p.recordMap[zoneName] = zoneRecordStore
	}
	zoneRecordStore.entries = append(zoneRecordStore.entries, recs...)
	return zoneRecordStore.entries, nil
}

func (p *Provider) DeleteRecords(ctx context.Context, zoneName string, recs []libdns.Record) ([]libdns.Record, error) {
	p.Lock()
	defer p.Unlock()
	var r libdns.Record
	zoneRecordStore := p.getZoneRecords(ctx, zoneName)
	if zoneRecordStore == nil {
		return nil, nil
	}
	zoneRecordStore.deleteRecords(recs)
	deletedRecords := []libdns.Record{}
	remainingRecords := []libdns.Record{}
	records := p.listAllRecords(ctx, zone)
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
	records := p.listAllRecords(ctx, zone)
	if len(records) == 0 {
		return records, fmt.Errorf("no records were found for %v", zone)
	}
	return records, nil
}

var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
