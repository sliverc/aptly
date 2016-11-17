package deb

import (
	"bytes"
	"fmt"
	"github.com/smira/aptly/database"
	"github.com/smira/go-uuid/uuid"
	"github.com/ugorji/go/codec"
	"log"
	"sync"
)

// LocalRepo is a collection of packages created locally
type LocalRepo struct {
	// Permanent internal ID
	UUID string `json:"-"`
	// User-assigned name
	Name string
	// Comment
	Comment string
	// DefaultDistribution
	DefaultDistribution string `codec:",omitempty"`
	// DefaultComponent
	DefaultComponent string `codec:",omitempty"`
	// Uploaders configuration
	Uploaders *Uploaders `code:",omitempty" json:"-"`
	// "Snapshot" of current list of packages
	packageRefs *PackageRefList
}

// NewLocalRepo creates new instance of Debian local repository
func NewLocalRepo(name string, comment string) *LocalRepo {
	return &LocalRepo{
		UUID:    uuid.New(),
		Name:    name,
		Comment: comment,
	}
}

// String interface
func (repo *LocalRepo) String() string {
	if repo.Comment != "" {
		return fmt.Sprintf("[%s]: %s", repo.Name, repo.Comment)
	}
	return fmt.Sprintf("[%s]", repo.Name)
}

// NumPackages return number of packages in local repo
func (repo *LocalRepo) NumPackages() int {
	if repo.packageRefs == nil {
		return 0
	}
	return repo.packageRefs.Len()
}

// RefList returns package list for repo
func (repo *LocalRepo) RefList() *PackageRefList {
	return repo.packageRefs
}

// UpdateRefList changes package list for local repo
func (repo *LocalRepo) UpdateRefList(reflist *PackageRefList) {
	repo.packageRefs = reflist
}

// Encode does msgpack encoding of LocalRepo
func (repo *LocalRepo) Encode() []byte {
	var buf bytes.Buffer

	encoder := codec.NewEncoder(&buf, &codec.MsgpackHandle{})
	encoder.Encode(repo)

	return buf.Bytes()
}

// Decode decodes msgpack representation into LocalRepo
func (repo *LocalRepo) Decode(input []byte) error {
	decoder := codec.NewDecoderBytes(input, &codec.MsgpackHandle{})
	return decoder.Decode(repo)
}

// Key is a unique id in DB
func (repo *LocalRepo) Key() []byte {
	return []byte("L" + repo.UUID)
}

// RefKey is a unique id for package reference list
func (repo *LocalRepo) RefKey() []byte {
	return []byte("E" + repo.UUID)
}

// LocalRepoCollection does listing, updating/adding/deleting of LocalRepos
type LocalRepoCollection struct {
	*sync.RWMutex
	db   database.Storage
	list []*LocalRepo
	// mutual exclusion for simple actions such as looping list and
	// reading/updating simple Repo fields
	SimpleMutex *sync.Mutex
}

// NewLocalRepoCollection loads LocalRepos from DB and makes up collection
func NewLocalRepoCollection(db database.Storage) *LocalRepoCollection {
	result := &LocalRepoCollection{
		RWMutex:      &sync.RWMutex{},
		db:           db,
		SimpleMutex: &sync.Mutex{},
	}

	blobs := db.FetchByPrefix([]byte("L"))
	result.list = make([]*LocalRepo, 0, len(blobs))

	for _, blob := range blobs {
		r := &LocalRepo{}
		if err := r.Decode(blob); err != nil {
			log.Printf("Error decoding repo: %s\n", err)
		} else {
			result.list = append(result.list, r)
		}
	}

	return result
}

// Add appends new repo to collection and saves it
func (collection *LocalRepoCollection) Add(repo *LocalRepo) error {
	collection.SimpleMutex.Lock()
	list := collection.list
	collection.SimpleMutex.Unlock()

	for _, r := range list {
		if r.Name == repo.Name {
			return fmt.Errorf("local repo with name %s already exists", repo.Name)
		}
	}

	err := collection.Update(repo)
	if err != nil {
		return err
	}

	collection.SimpleMutex.Lock()
	collection.list = append(collection.list, repo)
	collection.SimpleMutex.Unlock()
	return nil
}

// Update stores updated information about repo in DB
func (collection *LocalRepoCollection) Update(repo *LocalRepo) error {
	collection.SimpleMutex.Lock()
	{
		var list []*LocalRepo
		copy(list, collection.list)
		for i, r := range list {
			if r.UUID == repo.UUID {
				list[i] = repo
				break
			}
		}
		collection.list = list
	}
	collection.SimpleMutex.Unlock()

	err := collection.db.Put(repo.Key(), repo.Encode())
	if err != nil {
		return err
	}
	if repo.packageRefs != nil {
		err = collection.db.Put(repo.RefKey(), repo.packageRefs.Encode())
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadComplete loads additional information for local repo
func (collection *LocalRepoCollection) LoadComplete(repo *LocalRepo) error {
	encoded, err := collection.db.Get(repo.RefKey())
	if err == database.ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}

	repo.packageRefs = &PackageRefList{}
	return repo.packageRefs.Decode(encoded)
}

// ByName looks up repository by name
func (collection *LocalRepoCollection) ByName(name string) (LocalRepo, error) {
	collection.SimpleMutex.Lock()
	list := collection.list
	collection.SimpleMutex.Unlock()

	for _, r := range list {
		if r.Name == name {
			return *r, nil
		}
	}
	return LocalRepo{}, fmt.Errorf("local repo with name %s not found", name)
}

// ByUUID looks up repository by uuid
func (collection *LocalRepoCollection) ByUUID(uuid string) (LocalRepo, error) {
	collection.SimpleMutex.Lock()
	list := collection.list
	collection.SimpleMutex.Unlock()

	for _, r := range list {
		if r.UUID == uuid {
			return *r, nil
		}
	}
	return LocalRepo{}, fmt.Errorf("local repo with uuid %s not found", uuid)
}

// ForEach runs method for each repository
func (collection *LocalRepoCollection) ForEach(handler func(*LocalRepo) error) error {
	var err error

	collection.SimpleMutex.Lock()
	list := collection.list
	collection.SimpleMutex.Unlock()

	for _, r := range list {
		err = handler(r)
		if err != nil {
			return err
		}
	}
	return err
}

// Len returns number of remote repos
func (collection *LocalRepoCollection) Len() int {
	collection.SimpleMutex.Lock()
	defer collection.SimpleMutex.Unlock()
	return len(collection.list)
}

// Drop removes remote repo from collection
func (collection *LocalRepoCollection) Drop(repo *LocalRepo) error {
	collection.SimpleMutex.Lock()
	repoPosition := -1

	for i, r := range collection.list {
		if r == repo {
			repoPosition = i
			break
		}
	}

	if repoPosition == -1 {
		collection.SimpleMutex.Unlock()
		panic("local repo not found!")
	}

	collection.list[len(collection.list)-1], collection.list[repoPosition], collection.list =
		nil, collection.list[len(collection.list)-1], collection.list[:len(collection.list)-1]
	collection.SimpleMutex.Unlock()

	err := collection.db.Delete(repo.Key())
	if err != nil {
		return err
	}

	return collection.db.Delete(repo.RefKey())
}
