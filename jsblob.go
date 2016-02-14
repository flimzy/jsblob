/*
Package jsblob provides GopherJS bindings for the JavaScript Blob objects.

Read more about JavaScript Blobs here: https://developer.mozilla.org/en-US/docs/Web/API/Blob

A js.Object containing an existing Blob can be cast to a Blob object as follows:

    nativeBlob := js.Global.Get("Blob").New([]string{"some blobby data"})
    blob := jsblob.Blob{*nativeBlob}
    fmt.Println( blob.Size() ) // 16
*/
package jsblob

import (
	"github.com/gopherjs/gopherjs/js"
	"sync"
)

// Blob wraps a js.Object
type Blob struct {
	js.Object
}

type Options struct {
	Type    string `js:"type"`
	Endings string `js:"endings"`
}

// New returns a newly created Blob object whose content consists of the
// concatenation of the array of values given in parameter.
func New(parts []interface{}, opts Options) *Blob {
	blob := js.Global.Get("Blob").New(parts, opts)
	return &Blob{*blob}
}

// IsClosed returns true if the Close() method (or the underlying JavaScript
// Blobl.close() method) has been called on the blob. Closed blobs can not be
// read.
func (b *Blob) IsClosed() bool {
	return b.Get("isClosed").Bool()
}

// Size returns the size, in bytes, of the data contained in the Blob object.
func (b *Blob) Size() int {
	return b.Get("size").Int()
}

// Type returns a string indicating the MIME type of the data contained in
// the Blob. If the type is unknown, this string is empty.
func (b *Blob) Type() string {
	return b.Get("type").String()
}

// Close closes the blob object, possibly freeing underlying resources.
func (b *Blob) Close() {
	b.Call("close")
}

// Slice returns a new Blob object containing the specified range of bytes of the source Blob.
func (b *Blob) Slice(start, end int, contenttype string) *Blob {
	newBlob := b.Call("slice", start, end, contenttype)
	return &Blob{*newBlob}
}

// Bytes returns a slice of the contents of the Blob.
func (b *Blob) Bytes() []byte {
	fileReader := js.Global.Get("FileReader").New()
	var wg sync.WaitGroup
	var buf []byte
	wg.Add(1)
	fileReader.Set("onload", js.MakeFunc(func(this *js.Object, _ []*js.Object) interface{} {
		defer wg.Done()
		buf = js.Global.Get("Uint8Array").New(this.Get("result")).Interface().([]uint8)
		return nil
	}))
	fileReader.Call("readAsArrayBuffer", b)
	wg.Wait()
	return buf
}
