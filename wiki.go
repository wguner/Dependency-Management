package main

import (
	"fmt"
	"io/ioutil"
)

// creating a Page struct here
type Page struct {
	Title string
	// represent the context of the page
	Body []byte
}

// saving pages
// creating a function with a pointer to Page, returns type error, no parameters
func (p *Page) save() error {
	// will save Body to a text file, title as the filename
	filename := p.Title + ".txt"
	// 0600 represents read-only
	return ioutil.WriteFile(filename, p.Body, 0600)
	// save returns error to let user know if anything goes wrong.
}

//functions can return multiple values???
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	// reads file content into body
	body, err := ioutil.ReadFile(filename)
	// io.ReadFile returns []byte and error
	// we handle the error in here
	if err != nil {
		return nil, err
	}
	// if it is nil then it has successfully loaded a Page. If not, it will be an error that can be handled by the caller
	return &Page{Title: title, Body: body}, nil
}

// creates an .exe, when run it creates a .txt file with test value
func main() {
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
}
