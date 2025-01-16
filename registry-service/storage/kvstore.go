package storage

import (
	"log"
	"encoding/json"
	"strings"
	"fmt"

	"github.com/nats-io/nats.go"
	"registry-service/models"
)

type KVStore struct {
	bucket nats.KeyValue
}

var nc *nats.Conn

func InitNATSConnection(natsURL string) int {
	log.Printf("Connecting to NATS at %s", natsURL)
	var err error
	nc, err = nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
		return -1
	}
	log.Println("Connected to NATS")
	return 0
}

func NewKVStore(bucketName string) (*KVStore, error) {
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

func (kv *KVStore) StoreFunctionMetadata(UUID string, metadata *models.FunctionMetadata) error {
	data, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to serialize metadata: %w", err)
	}
	key := fmt.Sprintf("user_%s/%s", metadata.UUID, metadata.Name)

	_, err = kv.bucket.Put(key, data)
	if err != nil {
		return fmt.Errorf("failed to store metadata: %w", err)
	}
	return nil
}

func (kv *KVStore) GetFunctionMetadata(UUID, functionName string) (*models.FunctionMetadata, error) {
	key := fmt.Sprintf("user_%s/%s", UUID, functionName)
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

func (kv *KVStore) ListFunctions(UUID string) ([]models.FunctionMetadata, error) {
    // Retrieve all keys in the bucket
    keys, err := kv.bucket.Keys()
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve keys from the bucket: %w", err)
    }

    var metadataList []models.FunctionMetadata
    prefix := fmt.Sprintf("user_%s/", UUID)

    // Iterate over keys and filter based on the prefix
    for _, key := range keys {
        if !strings.HasPrefix(key, prefix) {
            continue
        }

        // Fetch the entry for the key
        entry, err := kv.bucket.Get(key)
        if err != nil {
            return nil, fmt.Errorf("failed to retrieve metadata for key %s: %w", key, err)
        }

        var metadata models.FunctionMetadata
        if err := json.Unmarshal(entry.Value(), &metadata); err != nil {
            return nil, fmt.Errorf("failed to deserialize metadata for key %s: %w", key, err)
        }

        metadataList = append(metadataList, metadata)
    }

    return metadataList, nil
}

func (kv *KVStore) DeleteFunctionMetadata(UUID, functionName string) error {
	key := fmt.Sprintf("user_%s/%s", UUID, functionName)
	if err := kv.bucket.Delete(key); err != nil {
		return fmt.Errorf("failed to delete metadata: %w", err)
	}
	return nil
}