package kolide

import (
	"testing"
	"time"

	"github.com/WatchBeam/clock"
	"github.com/stretchr/testify/assert"
)

func TestResetHosts(t *testing.T) {
	host := Host{}
	result := host.ResetPrimaryNetwork()
	assert.False(t, result)

	host.NetworkInterfaces = []*NetworkInterface{
		&NetworkInterface{
			ID:        1,
			IPAddress: "192.168.1.2",
		},
		&NetworkInterface{
			ID:        2,
			IPAddress: "192.168.1.3",
		},
	}

	result = host.ResetPrimaryNetwork()
	assert.True(t, result)
	assert.Equal(t, uint(1), *host.PrimaryNetworkInterfaceID)

	host.PrimaryNetworkInterfaceID = &host.NetworkInterfaces[1].ID
	result = host.ResetPrimaryNetwork()
	assert.False(t, result)
	assert.Equal(t, uint(2), *host.PrimaryNetworkInterfaceID)

	host.NetworkInterfaces = host.NetworkInterfaces[:1]
	result = host.ResetPrimaryNetwork()
	assert.True(t, result)
	assert.Equal(t, uint(1), *host.PrimaryNetworkInterfaceID)

	host.NetworkInterfaces = []*NetworkInterface{}
	result = host.ResetPrimaryNetwork()
	assert.True(t, result)
	assert.Nil(t, host.PrimaryNetworkInterfaceID)
}

func TestHostStatus(t *testing.T) {
	mockClock := clock.NewMockClock()

	var testCases = []struct {
		seenTime            time.Time
		distributedInterval uint
		configTLSRefresh    uint
		status              string
	}{
		{mockClock.Now().Add(-30 * time.Second), 10, 3600, StatusOnline},
		{mockClock.Now().Add(-45 * time.Second), 10, 3600, StatusOffline},
		{mockClock.Now().Add(-30 * time.Second), 3600, 10, StatusOnline},
		{mockClock.Now().Add(-45 * time.Second), 3600, 10, StatusOffline},

		{mockClock.Now().Add(-70 * time.Second), 60, 60, StatusOnline},
		{mockClock.Now().Add(-91 * time.Second), 60, 60, StatusOffline},

		{mockClock.Now().Add(-1 * time.Second), 10, 10, StatusOnline},
		{mockClock.Now().Add(-1 * time.Minute), 10, 10, StatusOffline},
		{mockClock.Now().Add(-31 * 24 * time.Hour), 10, 10, StatusMIA},

		// Ensure behavior is reasonable if we don't have the values
		{mockClock.Now().Add(-1 * time.Second), 0, 0, StatusOnline},
		{mockClock.Now().Add(-1 * time.Minute), 0, 0, StatusOffline},
		{mockClock.Now().Add(-31 * 24 * time.Hour), 0, 0, StatusMIA},
	}

	for _, tt := range testCases {
		t.Run("", func(t *testing.T) {
			// Save interval values
			h := Host{
				DistributedInterval: tt.distributedInterval,
				ConfigTLSRefresh:    tt.configTLSRefresh,
				SeenTime:            tt.seenTime,
			}

			assert.Equal(t, tt.status, h.Status(mockClock.Now()))
		})
	}

}

func TestHostIsNew(t *testing.T) {
	mockClock := clock.NewMockClock()

	host := Host{}

	host.CreatedAt = mockClock.Now().AddDate(0, 0, -1)
	assert.True(t, host.IsNew(mockClock.Now()))

	host.CreatedAt = mockClock.Now().AddDate(0, 0, -2)
	assert.False(t, host.IsNew(mockClock.Now()))
}
