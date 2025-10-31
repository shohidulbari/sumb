package db

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/blevesearch/bleve/v2"
	bolt "go.etcd.io/bbolt"
	berrosrs "go.etcd.io/bbolt/errors"
)

type Note struct {
	Body string `json:"body"`
}

type NoteResponse struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

func GetDb() (*bolt.DB, bleve.Index) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dbDir := filepath.Join(homeDir, ".sumb")
	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		log.Fatal(err)
	}
	dbPath := filepath.Join(dbDir, "sumb.db")
	blevePath := filepath.Join(dbDir, "sumb.bleve")

	db, err := bolt.Open(dbPath, 0o600, nil)
	if err != nil {
		log.Fatal(err)
	}

	indexMapping := bleve.NewIndexMapping()
	index, err := bleve.New(blevePath, indexMapping)
	if err != nil {
		index, err = bleve.Open(blevePath)
		if err != nil {
			log.Fatal(err)
		}
	}

	return db, index
}

func encodeInt64(n string) []byte {
	id, err := strconv.Atoi(n)
	if err != nil {
		log.Fatal(err)
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(id))
	return b
}

func decodeInt64(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}

func generateID(db *bolt.DB) string {
	currentSeq, err := GetByID([]byte("__id_seq__"), db)
	if err != nil {
		currentSeq = &Note{Body: "0"}
		err := db.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte("notes"))
			if err != nil {
				return err
			}
			noteData, err := json.Marshal(currentSeq)
			if err != nil {
				return err
			}
			return bucket.Put([]byte("__id_seq__"), noteData)
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	currentId, err := strconv.ParseInt(currentSeq.Body, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	newId := currentId + 1
	newSeq := &Note{Body: strconv.FormatInt(newId, 10)}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("notes"))
		if bucket == nil {
			return berrosrs.ErrBucketNotFound
		}
		noteData, err := json.Marshal(newSeq)
		if err != nil {
			return err
		}
		return bucket.Put([]byte("__id_seq__"), noteData)
	})
	if err != nil {
		log.Fatal(err)
	}

	return strconv.FormatInt(newId, 10)
}

func Create(body string) (*NoteResponse, error) {
	db, index := GetDb()
	defer db.Close()
	defer index.Close()

	note := &Note{
		Body: body,
	}
	id := generateID(db)
	idBytes := encodeInt64(id)

	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("notes"))
		if err != nil {
			return err
		}
		noteData, err := json.Marshal(note)
		if err != nil {
			return err
		}
		return bucket.Put(idBytes, noteData)
	})
	if err != nil {
		return nil, err
	}

	err = index.Index(id, note)
	if err != nil {
		return nil, err
	}

	return &NoteResponse{ID: id, Body: body}, nil
}

func Update(id string, body string) error {
	db, index := GetDb()
	defer db.Close()
	defer index.Close()

	idBytes := encodeInt64(id)

	note := &Note{
		Body: body,
	}
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("notes"))
		if bucket == nil {
			return berrosrs.ErrBucketNotFound
		}
		noteData, err := json.Marshal(note)
		if err != nil {
			return err
		}
		return bucket.Put(idBytes, noteData)
	})
	if err != nil {
		return err
	}

	err = index.Index(id, note)
	if err != nil {
		return err
	}

	return nil
}

func Delete(id string) error {
	db, index := GetDb()
	defer db.Close()
	defer index.Close()

	idBytes := encodeInt64(id)

	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("notes"))
		if bucket == nil {
			return berrosrs.ErrBucketNotFound
		}
		return bucket.Delete(idBytes)
	})
	if err != nil {
		return err
	}

	err = index.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func GetByID(id []byte, db *bolt.DB) (*Note, error) {
	var note Note
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("notes"))
		if bucket == nil {
			return berrosrs.ErrBucketNotFound
		}
		noteData := bucket.Get(id)
		if noteData == nil {
			return berrosrs.ErrInvalid
		}
		return json.Unmarshal(noteData, &note)
	})
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func Search(queryString string) ([]NoteResponse, error) {
	db, index := GetDb()
	defer db.Close()
	defer index.Close()

	query := bleve.NewQueryStringQuery(queryString)
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	var results []NoteResponse
	for _, hit := range searchResult.Hits {
		noteId := hit.ID
		idBytes := encodeInt64(noteId)
		note, err := GetByID(idBytes, db)
		if err != nil {
			return nil, err
		}
		noteResponse := &NoteResponse{
			ID:   noteId,
			Body: note.Body,
		}
		results = append(results, *noteResponse)
	}

	return results, nil
}

func ListLatestNotes(n int) ([]NoteResponse, error) {
	db, index := GetDb()
	defer db.Close()
	defer index.Close()

	var results []NoteResponse
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("notes"))
		if bucket == nil {
			return berrosrs.ErrBucketNotFound
		}

		c := bucket.Cursor()
		k, v := c.Last()
		count := 0
		for k != nil && count < n {
			if string(k) != "__id_seq__" {
				var note Note
				if err := json.Unmarshal(v, &note); err != nil {
					return err
				}
				noteId := decodeInt64(k)
				noteResponse := NoteResponse{
					ID:   strconv.FormatInt(noteId, 10),
					Body: note.Body,
				}
				results = append(results, noteResponse)
				count++
			}
			k, v = c.Prev()
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func GetNoteById(id string) (*NoteResponse, error) {
	db, index := GetDb()
	defer db.Close()
	defer index.Close()

	idBytes := encodeInt64(id)
	note, err := GetByID(idBytes, db)
	if err != nil {
		return nil, err
	}

	return &NoteResponse{
		ID:   id,
		Body: note.Body,
	}, nil
}
