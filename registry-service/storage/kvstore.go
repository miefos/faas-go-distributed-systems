package storage

import (
	"log"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"registry-service/models"
)

type KVStore struct {
	bucket nats.KeyValue
}

func NewKVStore(natsUrl, bucketName string) (*KVStore, error) {
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	} else {
		log.Printf("Connected to NATS")
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize JetStream: %w", err)
	} else {
		log.Printf("JetStream initialized")
	}

	bucket, err := js.KeyValue(bucketName)
	if err != nil {
		if bucket, err = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: bucketName,
		}); err != nil {
			return nil, fmt.Errorf("failed to create KeyValue bucket: %w", err)
		} else {
			log.Printf("KeyValue bucket created")
		}
	} else {
		log.Printf("KeyValue bucket found")
	}

	return &KVStore{bucket: bucket}, nil
}

func (kv *KVStore) StoreFunctionMetadata(userID string, metadata *models.FunctionMetadata) error {
	data, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to serialize metadata: %w", err)
	}
	key := fmt.Sprintf("user_%s/%s", userID, metadata.ID)

	_, err = kv.bucket.Put(key, data)
	if err != nil {
		return fmt.Errorf("failed to store metadata: %w", err)
	}
	return nil
}

func (kv *KVStore) GetFunctionMetadata(userID, functionID string) (*models.FunctionMetadata, error) {
	key := fmt.Sprintf("user_%s/%s", userID, functionID)
	entry, err := kv.bucket.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve metadata: %w", err)
	}

	var metadata models.FunctionMetadata
	if err := json.Unmarshal(entry.Value(), &metadata); err != nil {
		return nil, fmt.Errorf("failed to deserialize metadata: %w", err)
	}
	return &metadata, nil
}
