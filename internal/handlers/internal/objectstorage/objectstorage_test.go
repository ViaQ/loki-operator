package objectstorage_test

import (
	"testing"

	"github.com/ViaQ/loki-operator/internal/handlers/internal/objectstorage"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func TestExtract(t *testing.T) {
	type test struct {
		name            string
		secret          *corev1.Secret
		configMap       *corev1.ConfigMap
		insecure        bool
		correctEndpoint string
		wantErr         bool
	}
	table := []test{
		{
			name:      "minio: missing endpoint",
			secret:    &corev1.Secret{},
			configMap: &corev1.ConfigMap{},
			wantErr:   true,
		},
		{
			name: "minio: missing bucketnames",
			secret: &corev1.Secret{
				Data: map[string][]byte{
					"endpoint": []byte("here"),
				},
			},
			configMap: &corev1.ConfigMap{},
			wantErr:   true,
		},
		{
			name: "minio: missing access_key_id",
			secret: &corev1.Secret{
				Data: map[string][]byte{
					"endpoint":    []byte("here"),
					"bucketnames": []byte("this,that"),
				},
			},
			configMap: &corev1.ConfigMap{},
			wantErr:   true,
		},
		{
			name: "minio: missing access_key_secret",
			secret: &corev1.Secret{
				Data: map[string][]byte{
					"endpoint":      []byte("here"),
					"bucketnames":   []byte("this,that"),
					"access_key_id": []byte("id"),
				},
			},
			configMap: &corev1.ConfigMap{},
			wantErr:   true,
		},
		{
			name: "minio: all set",
			secret: &corev1.Secret{
				Data: map[string][]byte{
					"endpoint":          []byte("here"),
					"bucketnames":       []byte("this,that"),
					"access_key_id":     []byte("id"),
					"access_key_secret": []byte("secret"),
				},
			},
			configMap:       &corev1.ConfigMap{},
			correctEndpoint: "here",
		},
		{
			name:      "ocs: missing endpoint",
			secret:    &corev1.Secret{},
			configMap: &corev1.ConfigMap{},
			wantErr:   true,
		},
		{
			name:   "ocs: missing port",
			secret: &corev1.Secret{},
			configMap: &corev1.ConfigMap{
				Data: map[string]string{
					"BUCKET_HOST": "here",
				},
			},
			wantErr: true,
		},
		{
			name:   "ocs: missing bucketnames",
			secret: &corev1.Secret{},
			configMap: &corev1.ConfigMap{
				Data: map[string]string{
					"BUCKET_HOST": "here",
					"BUCKET_PORT": "42",
				},
			},
			wantErr: true,
		},
		{
			name:   "ocs: missing access_key_id",
			secret: &corev1.Secret{},
			configMap: &corev1.ConfigMap{
				Data: map[string]string{
					"BUCKET_HOST": "here",
					"BUCKET_PORT": "42",
					"BUCKET_NAME": "this,that",
				},
			},
			wantErr: true,
		},
		{
			name: "ocs: missing access_key_secret",
			secret: &corev1.Secret{
				Data: map[string][]byte{
					"AWS_ACCESS_KEY_ID": []byte("id"),
				},
			},
			configMap: &corev1.ConfigMap{
				Data: map[string]string{
					"BUCKET_HOST": "here",
					"BUCKET_PORT": "42",
					"BUCKET_NAME": "this,that",
				},
			},
			wantErr: true,
		},
		{
			name: "ocs: all set secure",
			secret: &corev1.Secret{
				Data: map[string][]byte{
					"AWS_ACCESS_KEY_ID":     []byte("id"),
					"AWS_SECRET_ACCESS_KEY": []byte("secret"),
				},
			},
			configMap: &corev1.ConfigMap{
				Data: map[string]string{
					"BUCKET_HOST": "here",
					"BUCKET_PORT": "42",
					"BUCKET_NAME": "this,that",
				},
			},
			insecure:        false,
			correctEndpoint: "https://here:42",
		},
		{
			name: "ocs: all set insecure",
			secret: &corev1.Secret{
				Data: map[string][]byte{
					"AWS_ACCESS_KEY_ID":     []byte("id"),
					"AWS_SECRET_ACCESS_KEY": []byte("secret"),
				},
			},
			configMap: &corev1.ConfigMap{
				Data: map[string]string{
					"BUCKET_HOST": "here",
					"BUCKET_PORT": "42",
					"BUCKET_NAME": "this,that",
				},
			},
			insecure:        true,
			correctEndpoint: "http://here:42",
		},
	}
	for _, tst := range table {
		tst := tst
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()

			ret, err := objectstorage.Extract(tst.secret, tst.configMap, tst.insecure)
			if !tst.wantErr {
				require.NoError(t, err)
				require.Equal(t, ret.Endpoint, tst.correctEndpoint)
			}
			if tst.wantErr {
				require.NotNil(t, err)
			}
		})
	}
}
