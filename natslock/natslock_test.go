package natslock

import (
	"reflect"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var nc *nats.Conn

var jetstream nats.JetStreamContext

func TestMain(m *testing.M) {
	natsSrv, err := natsserver.NewServer(&natsserver.Options{
		Host:      "127.0.0.1",
		Port:      natsserver.RANDOM_PORT,
		JetStream: true,
	})
	if err != nil {
		panic(err)
	}

	defer natsSrv.Shutdown()

	if err := natsserver.Run(natsSrv); err != nil {
		panic(err)
	}

	nc, err = nats.Connect(natsSrv.ClientURL(), nats.Timeout(time.Second))
	if err != nil {
		panic(err)
	}

	jetstream, err = nc.JetStream(nats.MaxWait(time.Second))
	if err != nil {
		panic(err)
	}

	m.Run()
}

func TestNew(t *testing.T) {
	locker := New()

	lockerType := reflect.TypeOf(locker).String()
	assert.Equal(t, "*natslock.Locker", lockerType)

	locker = New(WithLogger(zap.NewExample()))
	assert.Equal(t, true, locker.Logger.Core().Enabled(zap.DebugLevel))

	const bucketName = "test-bucket-1"

	kvStore, err := NewKeyValue(jetstream, bucketName, time.Minute)
	assert.NoError(t, err)

	defer func() {
		err := jetstream.DeleteKeyValue(kvStore.Bucket())
		assert.NoError(t, err)
	}()

	locker = New(WithKeyValueStore(kvStore))
	assert.Equal(t, kvStore, locker.KVStore)
	assert.Equal(t, bucketName, locker.KVStore.Bucket())
}

func TestNewKeyValue(t *testing.T) {
	const testName = "test-bucket-2"

	const testTTL = time.Minute

	_, err := NewKeyValue(jetstream, testName, 0)
	assert.Error(t, err)

	_, err = NewKeyValue(jetstream, "", testTTL)
	assert.Error(t, err)

	// first call should create a new bucket
	got, err := NewKeyValue(jetstream, testName, testTTL)
	assert.NoError(t, err)

	defer func() {
		err := jetstream.DeleteKeyValue(got.Bucket())
		assert.NoError(t, err)
	}()

	assert.Equal(t, testName, got.Bucket())

	status, err := got.Status()
	assert.NoError(t, err)
	assert.Equal(t, testTTL, status.TTL())

	// second call should return the existing bucket
	got, err = NewKeyValue(jetstream, testName, testTTL)
	assert.NoError(t, err)
	assert.Equal(t, testName, got.Bucket())

	status, err = got.Status()
	assert.NoError(t, err)
	assert.Equal(t, testTTL, status.TTL())
}

func TestLocker_AquireLead(t *testing.T) {
	const testName = "test-bucket-acquire-lead"

	const testTTL = time.Minute

	kv, err := NewKeyValue(jetstream, testName, testTTL)
	assert.NoError(t, err)

	defer func() {
		_ = jetstream.DeleteKeyValue(kv.Bucket())
	}()

	l := New(WithKeyValueStore(kv))

	_, err = l.AcquireLead(uuid.FromStringOrNil(""))
	assert.Error(t, err)

	testID, err := uuid.DefaultGenerator.NewV4()
	assert.NoError(t, err)

	testID2, err := uuid.DefaultGenerator.NewV4()
	assert.NoError(t, err)

	// when no key exists, we should be able to acquire the lead
	isLead, err := l.AcquireLead(testID)
	assert.NoError(t, err)
	assert.Equal(t, true, isLead)

	// if the key exists with our id, we should maintain the lead
	isLead, err = l.AcquireLead(testID)
	assert.NoError(t, err)
	assert.Equal(t, true, isLead)

	// if the key exists with a different id, we should not get the lead
	isLead, err = l.AcquireLead(testID2)
	assert.NoError(t, err)
	assert.Equal(t, false, isLead)

	// if the bucket doesn't exist, we should get an error
	err = jetstream.DeleteKeyValue(kv.Bucket())
	assert.NoError(t, err)

	_, err = l.AcquireLead(testID)
	assert.Error(t, err)
}

func TestLocker_ReleaseLead(t *testing.T) {
	const testName = "test-bucket-release-lead"

	const testTTL = time.Minute

	kv, err := NewKeyValue(jetstream, testName, testTTL)
	assert.NoError(t, err)

	defer func() {
		_ = jetstream.DeleteKeyValue(kv.Bucket())
	}()

	l := New(WithKeyValueStore(kv))

	err = l.ReleaseLead(uuid.FromStringOrNil(""))
	assert.Error(t, err)

	testID, err := uuid.DefaultGenerator.NewV4()
	assert.NoError(t, err)

	testID2, err := uuid.DefaultGenerator.NewV4()
	assert.NoError(t, err)

	// ok to release lead when key does not exist
	err = l.ReleaseLead(testID)
	assert.NoError(t, err)

	_, err = l.KVStore.PutString(l.KVKey, testID.String())
	assert.NoError(t, err)

	// when key exists and the value doesn't match the given id, it should be left
	err = l.ReleaseLead(testID2)
	assert.NoError(t, err)

	_, err = l.KVStore.Get(l.KVKey)
	assert.NoError(t, err)

	// when key exists and the value mathes the given id, it should be deleted
	err = l.ReleaseLead(testID)
	assert.NoError(t, err)

	_, err = l.KVStore.Get(l.KVKey)
	assert.EqualError(t, err, nats.ErrKeyNotFound.Error())

	// if the bucket doesn't exist, we should get an error
	err = jetstream.DeleteKeyValue(kv.Bucket())
	assert.NoError(t, err)

	err = l.ReleaseLead(testID)
	assert.Error(t, err)
}

func TestLocker_Name(t *testing.T) {
	kv, err := jetstream.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "test-bucket-name",
		TTL:    time.Minute,
	})
	if err != nil {
		panic(err)
	}

	defer func() {
		err := jetstream.DeleteKeyValue(kv.Bucket())
		assert.NoError(t, err)
	}()

	l := &Locker{KVStore: kv}
	assert.Equal(t, "test-bucket-name", l.Name())
}

func TestLocker_TTL(t *testing.T) {
	kv, err := jetstream.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "test-bucket-ttl",
		TTL:    time.Minute,
	})
	if err != nil {
		panic(err)
	}

	defer func() {
		err := jetstream.DeleteKeyValue(kv.Bucket())
		assert.NoError(t, err)
	}()

	l := &Locker{KVStore: kv}
	assert.Equal(t, time.Minute, l.TTL())
}
