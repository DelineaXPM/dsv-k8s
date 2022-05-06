package config

import "testing"

func TestMakeCredentialsValid(t *testing.T) {
	if _, err := MakeCredentials([]byte(`
	{
		"a": {
			"credentials": {
				"clientId": "x",
				"clientSecret": "y"
			},
			"tenant": "i"
		},
		"b": {
			"credentials": {
				"clientId": "x",
				"clientSecret": "y"
			},
			"tenant": "j"
		}
	}`)); err != nil {
		t.Errorf("MakeCredentials should not have failed")
	}
}

func TestMakeCredentialsInvalid(t *testing.T) {
	if _, err := MakeCredentials([]byte(`
	{
		"a": {
			"credentials": {
				"clientId": "x",
			},
			"tenant": "i"
		}
	}`)); err == nil {
		t.Errorf("MakeCredentials should have failed")
	}
}

func TestCredentials(t *testing.T) {
	if credentials, err := MakeCredentials([]byte(`
	{
		"a": {
			"credentials": {
				"clientId": "x",
				"clientSecret": "y"
			},
			"tenant": "i"
		},
		"b": {
			"credentials": {
				"clientId": "x",
				"clientSecret": "y"
			},
			"tenant": "j"
		},
		"c": {
			"credentials": {
				"clientId": "x",
				"clientSecret": "y"
			},
			"tenant": "k"
		}
	}
`)); err != nil {
		t.Error(err)
	} else {
		if len(*credentials) != 3 {
			t.Errorf("expected 3 credential, got %d", len(*credentials))
		}
		if names := credentials.Names(); len(names) != 3 {
			t.Errorf("expected 3 name, got %d", len(names))
		} else {
			switch names[0] {
			case "a":
			case "b":
			case "c":
				break
			default:
				t.Errorf("unexpected name '%s'", names[0])
			}
		}
	}
}
