package objectstorage

import (
	"github.com/ViaQ/logerr/kverrors"
	"github.com/ViaQ/loki-operator/internal/manifests"

	corev1 "k8s.io/api/core/v1"
)

// Extract reads a k8s secret and config map into a manifest
// object storage struct if valid. Try to read the secret
// first, then the config map and exit with error if neither
// one has the required values
func Extract(s *corev1.Secret, c *corev1.ConfigMap, insecure bool) (*manifests.ObjectStorage, error) {
	// Extract and validate mandatory fields
	var endpoint string
	var buckets string
	// Check if we have a correct config map
	if _, ok := c.Data["BUCKET_NAME"]; ok {
		buckets, ok = c.Data["BUCKET_NAME"]
		if !ok {
			return nil, kverrors.New("missing config map field", "field", "BUCKET_NAME")
		}
		host, ok := c.Data["BUCKET_HOST"]
		if !ok {
			return nil, kverrors.New("missing config map field", "field", "BUCKET_HOST")
		}
		port, ok := c.Data["BUCKET_PORT"]
		if !ok {
			return nil, kverrors.New("missing config map field", "field", "BUCKET_PORT")
		}

		if insecure {
			endpoint = "http://" + host + ":" + port
		} else {
			endpoint = "https://" + host + ":" + port
		}
	} else {
		e, ok := s.Data["endpoint"]
		if !ok {
			return nil, kverrors.New("missing secret field", "field", "endpoint")
		}
		endpoint = string(e)

		b, ok := s.Data["bucketnames"]
		if !ok {
			return nil, kverrors.New("missing secret field", "field", "bucketnames")
		}
		buckets = string(b)
	}
	// TODO buckets are comma-separated list
	id, ok := s.Data["access_key_id"]
	if !ok {
		id, ok = s.Data["AWS_ACCESS_KEY_ID"]
		if !ok {
			return nil, kverrors.New("missing secret field", "field", "access_key_id or AWS_ACCESS_KEY_ID")
		}
	}
	secret, ok := s.Data["access_key_secret"]
	if !ok {
		secret, ok = s.Data["AWS_SECRET_ACCESS_KEY"]
		if !ok {
			return nil, kverrors.New("missing secret field", "field", "access_key_secret or AWS_SECRET_ACCESS_KEY")
		}
	}

	// Extract and validate optional fields
	region, ok := s.Data["region"]
	if !ok {
		region = []byte("")
	}

	return &manifests.ObjectStorage{
		Endpoint:        string(endpoint),
		Buckets:         string(buckets),
		AccessKeyID:     string(id),
		AccessKeySecret: string(secret),
		Region:          string(region),
	}, nil
}
