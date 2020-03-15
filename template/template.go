package template

import (
	"bytes"
	"github.com/Streamlet/NoteIsSite/config"
	"io/ioutil"
	"sync"
	"text/template"
)

type BasicItem struct {
	Uri        string
	Name       string
	IsAncestor bool
	Children   []*BasicItem
	Parent     *BasicItem
}

type PageData struct {
	*BasicItem
	Content string
}

func (item BasicItem) HasChildren() bool {
	return item.Children != nil && len(item.Children) > 0
}

func (item BasicItem) Root() *BasicItem {
	root := &item
	for root.Parent != nil {
		root = root.Parent
	}
	return root
}

func (item BasicItem) Ancestors() []*BasicItem {
	ancestors := make([]*BasicItem, 0)
	for p := &item; p != nil; p = p.Parent {
		ancestors = append([]*BasicItem{p}, ancestors...)
	}
	return ancestors
}

type Executor interface {
	Update(templateRoot string) error

	GetIndex(data PageData) ([]byte, error)
	GetCategory(data PageData) ([]byte, error)
	GetContent(data PageData) ([]byte, error)

	Get404() []byte
	Get500() []byte
}

func NewExecutor(templateRoot string) (Executor, error) {
	td := new(templateData)
	err := td.Update(templateRoot)
	if err != nil {
		return nil, err
	}
	return td, nil
}

type templateData struct {
	lock             sync.RWMutex
	indexTemplate    string
	categoryTemplate string
	contentTemplate  string
	err404           []byte
	err500           []byte
}

func (td *templateData) Update(templateRoot string) error {
	c := config.GetSiteConfig().Template
	index, err := ioutil.ReadFile(templateRoot + "/" + c.IndexTemplate)
	if err != nil {
		return err
	}
	category, err := ioutil.ReadFile(templateRoot + "/" + c.CategoryTemplate)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadFile(templateRoot + "/" + c.ContentTemplate)
	if err != nil {
		return err
	}
	err404, _ := ioutil.ReadFile(templateRoot + "/" + c.ErrorPage404)
	err500, _ := ioutil.ReadFile(templateRoot + "/" + c.ErrorPage500)

	defer td.lock.Unlock()
	td.lock.Lock()

	td.indexTemplate = string(index)
	td.categoryTemplate = string(category)
	td.contentTemplate = string(content)
	td.err404 = err404
	td.err500 = err500

	return nil
}

func (td templateData) GetIndex(data PageData) ([]byte, error) {
	defer td.lock.RUnlock()
	td.lock.RLock()

	return td.execute(td.indexTemplate, data)
}

func (td templateData) GetCategory(data PageData) ([]byte, error) {
	defer td.lock.RUnlock()
	td.lock.RLock()

	return td.execute(td.categoryTemplate, data)
}

func (td templateData) GetContent(data PageData) ([]byte, error) {
	defer td.lock.RUnlock()
	td.lock.RLock()

	return td.execute(td.contentTemplate, data)
}

func (td templateData) execute(tmpl string, data interface{}) ([]byte, error) {
	tt := template.New("")
	_, err := tt.Parse(tmpl)
	if err != nil {
		return td.err500, err
	}

	var buffer []byte
	w := bytes.NewBuffer(buffer)
	err = tt.Execute(w, data)
	if err != nil {
		return td.err500, err
	}
	return w.Bytes(), nil
}

func (td templateData) Get404() []byte {
	defer td.lock.RUnlock()
	td.lock.RLock()

	return td.err404
}

func (td templateData) Get500() []byte {
	defer td.lock.RUnlock()
	td.lock.RLock()

	return td.err500
}
