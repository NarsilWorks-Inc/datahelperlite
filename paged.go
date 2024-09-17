package datahelperlite

// Pager is an interface for data
//
// A pager must be implemented as a singleton or a global object
// where storing, fetching and trashing should provide a ReadWrite mutex
// A possible implementation of this is via Redis
type Pager interface {
	Store(pageId string, pageNumber int, data []byte) error // Store data via page id. If page id is not set, this function should fail.
	Fetch(pageId string, pageNumber int) ([]byte, error)    // Retrieves a page data
	Trash(pageId string) error                              // Trash a page cache
	SizeSet(pageId string, size int)                        // Set size of a page
	SizeGet(pageId string) int                              // Get size of a page
}
