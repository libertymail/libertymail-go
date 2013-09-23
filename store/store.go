package store

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"libertymail-go/bits/wif"
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

	_, err = os.Stat(s.privFile)
	if err != nil {
		privdb = false
	}

	_, err = os.Stat(s.pubFile)
	if err != nil {
		pubdb = false
	}

	s.priv, err = sql.Open("sqlite3", s.privFile)
	if err != nil {
		return err
	}

	s.pub, err = sql.Open("sqlite3", s.pubFile)
	if err != nil {
		return err
	}

	if !privdb {
		sql := `
    	create table addresses (version integer, privacy integer, identifier text, wif text); 
    	delete from addresses;
		`
		_, err = s.priv.Exec(sql)
		if err != nil {
			return err
		}
	}

	if !pubdb {
		sql := `
    	create table pubkeys (key blob); 
    	delete from pubkeys;
		`
		_, err = s.pub.Exec(sql)
		if err != nil {
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
