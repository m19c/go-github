// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestGitService_GetTree(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/u/r/git/trees/s", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{
			  "sha": "s",
			  "tree": [ { "type": "blob" } ]
			}`)
	})

	tree, err := client.Git.GetTree("u", "r", "s", true)
	if err != nil {
		t.Errorf("Git.GetTree returned error: %v", err)
	}

	want := Tree{
		SHA: "s",
		Entries: []TreeEntry{
			TreeEntry{
				Type: "blob",
			},
		},
	}
	if !reflect.DeepEqual(*tree, want) {
		t.Errorf("Tree.Get returned %+v, want %+v", *tree, want)
	}
}

func TestGitService_CreateTree(t *testing.T) {
	setup()
	defer teardown()

	input := []TreeEntry{
		TreeEntry{
			Path: "file.rb",
			Mode: "100644",
			Type: "blob",
			SHA:  "7c258a9869f33c1e1e1f74fbb32f07c86cb5a75b",
		},
	}

	mux.HandleFunc("/repos/u/r/git/trees/s", func(w http.ResponseWriter, r *http.Request) {
		v := new(createTree)
		json.NewDecoder(r.Body).Decode(v)

		if m := "POST"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}

		want := &createTree{
			BaseTree: "b",
			Entries:  input,
		}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Git.CreateTree request body: %+v, want %+v", v, want)
		}

		fmt.Fprint(w, `{
		 "sha": "cd8274d15fa3ae2ab983129fb037999f264ba9a7",
		  "tree": [
		    {
		      "path": "file.rb",
		      "mode": "100644",
		      "type": "blob",
		      "size": 132,
		      "sha": "7c258a9869f33c1e1e1f74fbb32f07c86cb5a75b"
		    }
		  ]
		}`)
	})

	tree, err := client.Git.CreateTree("u", "r", "s", "b", input)
	if err != nil {
		t.Errorf("Git.CreateTree returned error: %v", err)
	}

	want := Tree{
		"cd8274d15fa3ae2ab983129fb037999f264ba9a7",
		[]TreeEntry{
			TreeEntry{
				Path: "file.rb",
				Mode: "100644",
				Type: "blob",
				Size: 132,
				SHA:  "7c258a9869f33c1e1e1f74fbb32f07c86cb5a75b",
			},
		},
	}

	if !reflect.DeepEqual(*tree, want) {
		t.Errorf("Git.CreateTree returned %+v, want %+v", *tree, want)
	}
}
