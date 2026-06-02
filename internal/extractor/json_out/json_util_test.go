package json_out

import (
	"reflect"
	"testing"
)

func TestParseJSONPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    []pathPart
		wantErr bool
	}{
		{
			name: "nested object and array",
			path: "items[0].name",
			want: []pathPart{
				{Key: "items", IsArray: true, Index: intPtr(0)},
				{Key: "name"},
			},
		},
		{
			name: "root array part unsupported by setter but parsable",
			path: "[1]",
			want: []pathPart{{IsArray: true, Index: intPtr(1)}},
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
		{
			name:    "empty part",
			path:    "a..b",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseJSONPath(tt.path)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseJSONPath(%q) expected error", tt.path)
				}
				return
			}

			if err != nil {
				t.Fatalf("parseJSONPath(%q) unexpected error: %v", tt.path, err)
			}

			if !equalPathParts(got, tt.want) {
				t.Fatalf("parseJSONPath(%q) = %#v, want %#v", tt.path, got, tt.want)
			}
		})
	}
}

func TestParsePathPart(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    pathPart
		wantErr bool
	}{
		{
			name: "plain key",
			raw:  "field",
			want: pathPart{Key: "field"},
		},
		{
			name: "array only",
			raw:  "[2]",
			want: pathPart{IsArray: true, Index: intPtr(2)},
		},
		{
			name: "key with array index",
			raw:  "items[3]",
			want: pathPart{Key: "items", IsArray: true, Index: intPtr(3)},
		},
		{
			name:    "invalid array index",
			raw:     "items[x]",
			wantErr: true,
		},
		{
			name:    "missing array close bracket",
			raw:     "items[1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePathPart(tt.raw)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("parsePathPart(%q) expected error", tt.raw)
				}
				return
			}

			if err != nil {
				t.Fatalf("parsePathPart(%q) unexpected error: %v", tt.raw, err)
			}

			if !equalPathPart(got, tt.want) {
				t.Fatalf("parsePathPart(%q) = %#v, want %#v", tt.raw, got, tt.want)
			}
		})
	}
}

func TestSetJSONPath(t *testing.T) {
	t.Run("create nested objects and arrays", func(t *testing.T) {
		root := map[string]any{}

		if err := setJSONPath(&root, "users[1].name", "alice"); err != nil {
			t.Fatalf("setJSONPath unexpected error: %v", err)
		}

		usersAny, ok := root["users"]
		if !ok {
			t.Fatalf("root missing users key")
		}

		users, ok := usersAny.([]any)
		if !ok {
			t.Fatalf("users type = %T, want []any", usersAny)
		}

		if len(users) != 2 {
			t.Fatalf("len(users) = %d, want 2", len(users))
		}

		if users[0] != nil {
			t.Fatalf("users[0] = %v, want nil", users[0])
		}

		user1, ok := users[1].(map[string]any)
		if !ok {
			t.Fatalf("users[1] type = %T, want map[string]any", users[1])
		}

		if user1["name"] != "alice" {
			t.Fatalf("users[1].name = %v, want alice", user1["name"])
		}
	})

	t.Run("duplicate path should error", func(t *testing.T) {
		root := map[string]any{}

		if err := setJSONPath(&root, "a.b", 1); err != nil {
			t.Fatalf("first setJSONPath unexpected error: %v", err)
		}

		if err := setJSONPath(&root, "a.b", 2); err == nil {
			t.Fatalf("second setJSONPath expected duplicate error")
		}
	})

	t.Run("path conflict non object should error", func(t *testing.T) {
		root := map[string]any{"a": 1}

		if err := setJSONPath(&root, "a.b", 2); err == nil {
			t.Fatalf("setJSONPath expected conflict error")
		}
	})

	t.Run("root array path should error", func(t *testing.T) {
		root := map[string]any{}

		if err := setJSONPath(&root, "[0]", "x"); err == nil {
			t.Fatalf("setJSONPath expected root array path error")
		}
	})
}

func TestEnsureArrayLen(t *testing.T) {
	orig := []any{"x"}
	got := ensureArrayLen(orig, 3)

	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}

	if got[0] != "x" {
		t.Fatalf("got[0] = %v, want x", got[0])
	}

	if got[1] != nil || got[2] != nil {
		t.Fatalf("new elements = %#v, want nil, nil", got[1:])
	}

	unchanged := ensureArrayLen([]any{"a", "b"}, 1)
	if !reflect.DeepEqual(unchanged, []any{"a", "b"}) {
		t.Fatalf("ensureArrayLen should not shrink, got %#v", unchanged)
	}
}

func intPtr(v int) *int {
	return &v
}

func equalPathParts(a, b []pathPart) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !equalPathPart(a[i], b[i]) {
			return false
		}
	}

	return true
}

func equalPathPart(a, b pathPart) bool {
	if a.Key != b.Key || a.IsArray != b.IsArray {
		return false
	}

	if (a.Index == nil) != (b.Index == nil) {
		return false
	}

	if a.Index != nil && *a.Index != *b.Index {
		return false
	}

	return true
}
