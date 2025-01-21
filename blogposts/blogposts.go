package blogposts

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"strings"
)

type Post struct {
	Title       string
	Description string
	Tags        []string
	Body        string
}

func getPost(fileSystem fs.ReadDirFS, filename string) (Post, error) {
	postFile, err := fileSystem.Open(filename)
	if err != nil {
		return Post{}, err
	}
	defer postFile.Close()
	return newPost(postFile)
}

func readMetadataLine(scanner *bufio.Scanner, tagName string) (string, error) {
	scanner.Scan()
	line := scanner.Text()
	if !strings.HasPrefix(line, tagName) {
		return "", fmt.Errorf("line does not start with %s", tagName)
	}
	return strings.TrimPrefix(line, tagName), nil
}

func readBody(scanner *bufio.Scanner) string {
	scanner.Scan() // ignore a line with '---'
	buf := bytes.Buffer{}
	for scanner.Scan() {
		fmt.Fprintln(&buf, scanner.Text())
	}
	return strings.TrimSuffix(buf.String(), "\n")
}

func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)

	postTitle, err := readMetadataLine(scanner, "Title: ")
	if err != nil {
		return Post{}, err
	}
	postDescription, err := readMetadataLine(scanner, "Description: ")
	if err != nil {
		return Post{}, err
	}
	rawPostTags, err := readMetadataLine(scanner, "Tags: ")
	if err != nil {
		return Post{}, err
	}
	postTags := strings.Split(rawPostTags, ", ")
	postBody := readBody(scanner)
	post := Post{
		Title:       postTitle,
		Description: postDescription,
		Tags:        postTags,
		Body:        postBody,
	}
	return post, nil
}

func NewPostsFromFS(fileSystem fs.ReadDirFS) ([]Post, error) {
	dir, err := fileSystem.ReadDir(".")
	if err != nil {
		return nil, err
	}
	var posts []Post
	for _, entry := range dir {
		post, err := getPost(fileSystem, entry.Name())
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
