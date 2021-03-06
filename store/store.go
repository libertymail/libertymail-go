package store

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"libertymail-go/bits/encoding/wif"
	"libertymail-go/proto"
)

type store struct {
	priv, pub         *sql.DB
	privFile, pubFile string
}

func NewStore(privFile, pubFile string) *store {

	s := new(store)
	s.privFile = privFile
	s.pubFile = pubFile
	return s
}

func (s *store) Open() error {

	var err error
	privdb, pubdb := true, true

	if _, err = os.Stat(s.privFile); err != nil {
		privdb = false
	}

	if _, err = os.Stat(s.pubFile); err != nil {
		pubdb = false
	}

	if s.priv, err = sql.Open("sqlite3", s.privFile); err != nil {
		return err
	}

	if s.pub, err = sql.Open("sqlite3", s.pubFile); err != nil {
		return err
	}

	if !privdb {
		sql := `
    	create table addresses (version integer, privacy integer, identifier text, wif text); 
    	delete from addresses;
		`
		if _, err = s.priv.Exec(sql); err != nil {
			return err
		}
	}

	if !pubdb {
		sql := `
    	create table pubkeys (key blob); 
    	delete from pubkeys;
		`
		if _, err = s.pub.Exec(sql); err != nil {
			return err
		}
	}

	return nil
}

func (s *store) Close() {

	s.priv.Close()
	s.pub.Close()
}

func (s *store) SaveAddress(addr *proto.Address) error {

	tx, err := s.priv.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert into addresses(version, privacy, identifier, wif) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	w, err := wif.Encode(addr.Key)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(addr.Version, addr.Privacy, addr.Identifier, w)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}
