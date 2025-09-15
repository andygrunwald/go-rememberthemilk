package rememberthemilk

import (
	"net/url"
	"testing"
)

func TestClient_SignRequest(t *testing.T) {
	tests := []struct {
		name         string
		sharedSecret string
		params       url.Values
		expected     string
	}{
		{
			name:         "RTM documentation example",
			sharedSecret: "BANANAS",
			params: url.Values{
				"yxz": []string{"foo"},
				"feg": []string{"bar"},
				"abc": []string{"baz"},
			},
			expected: "82044aae4dd676094f23f1ec152159ba",
		},
		{
			name:         "empty parameters",
			sharedSecret: "BANANAS",
			params:       url.Values{},
			expected:     "96616b070abae0ea857ee4ae67c39b8f", // MD5 hash of "BANANAS"
		},
		{
			name:         "single parameter",
			sharedSecret: "SECRET",
			params: url.Values{
				"key": []string{"value"},
			},
			expected: "77703a75c535fe087f2814d7473b9a45", // MD5 hash of "SECRETkeyvalue"
		},
		{
			name:         "multiple parameters with special characters",
			sharedSecret: "TEST",
			params: url.Values{
				"api_key": []string{"123456"},
				"method":  []string{"rtm.test.echo"},
				"format":  []string{"json"},
			},
			expected: "71dd4d64b84a082be2949a31c9b93803", // MD5 hash of "TESTapi_key123456formatjsonmethodrtm.test.echo"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				sharedSecret: tt.sharedSecret,
			}
			result := client.SignRequest(tt.params)

			if result != tt.expected {
				t.Errorf("SignRequest() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
