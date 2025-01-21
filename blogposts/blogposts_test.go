package blogposts_test

import (
	"blogposts"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestNewBlogPosts(t *testing.T) {

	const (
		firstBody  = "Title: Post 1" + "\nDescription: Description 1" + "\nTags: rust, tdd" + "\n---" + "\nHello world!"
		secondBody = "Title: Post 2\n" + "Description: Description 2\n" + "Tags: python, web" + "\n---" + "\nmy\nnem\njeff"
	)

	fs := fstest.MapFS{
		"hello world.md":  {Data: []byte(firstBody)},
		"hello-world2.md": {Data: []byte(secondBody)},
	}

	got, err := blogposts.NewPostsFromFS(fs)

	if err != nil {
		t.Fatal(err)
	}

	want := []blogposts.Post{
		{
			Title:       "Post 1",
			Description: "Description 1",
			Tags:        []string{"rust", "tdd"},
			Body:        "Hello world!",
		},
		{
			Title:       "Post 2",
			Description: "Description 2",
			Tags:        []string{"python", "web"},
			Body:        "my\nnem\njeff",
		},
	}

	assertPost(t, got[0], want[0])
	assertPost(t, got[1], want[1])
}

func assertPost(t *testing.T, got blogposts.Post, want blogposts.Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
