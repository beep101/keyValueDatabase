package storage

import (
	"errors"
	"os"
)

const (
	TreeExt = ".dbs"
	DataExt = ".dbd"
)

type DbFiles struct {
	//basic characteristic
	dbAddress string
	dbName    string
	//data files
	tree *os.File
	data *os.File
}

func (dbf *DbFiles) checkFiles() bool {
	ext := true

	_, err := os.Stat(dbf.dbAddress + dbf.dbName + TreeExt)
	if os.IsNotExist(err) {
		ext = false
	}
	_, err = os.Stat(dbf.dbAddress + dbf.dbName + DataExt)
	if os.IsNotExist(err) {
		ext = false
	}

	return ext
}

func (dbf *DbFiles) createFiles() error {
	dbs, err := os.Create(dbf.dbAddress + dbf.dbName + TreeExt)
	if err != nil {
		return err
	}

	dbd, err := os.Create(dbf.dbAddress + dbf.dbName + DataExt)
	if err != nil {
		return err
	}

	dbf.tree = dbs
	dbf.data = dbd

	err = dbf.WriteRootAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return err
	}
	err = dbf.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (dbf *DbFiles) openFiles() error {
	dbs, err := os.OpenFile(dbf.dbAddress+dbf.dbName+TreeExt, os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	dbd, err := os.OpenFile(dbf.dbAddress+dbf.dbName+DataExt, os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	dbf.tree = dbs
	dbf.data = dbd

	return nil
}

func Check(dbadr, dbname string) bool {
	dbf := &DbFiles{dbAddress: dbadr, dbName: dbname}
	return dbf.checkFiles()
}

// open and close funcs
func Open(dbadr, dbname string) (*DbFiles, error) {
	dbf := &DbFiles{dbAddress: dbadr, dbName: dbname}
	ext := dbf.checkFiles()
	if !ext {
		return nil, errors.New("Error: Missing database files")
	}
	err := dbf.openFiles()
	if err != nil {
		return nil, err
	}
	return dbf, nil
}

func CreateAndOpen(dbadr, dbname string) (*DbFiles, error) {
	dbf := &DbFiles{dbAddress: dbadr, dbName: dbname}
	ext := dbf.checkFiles()
	if ext {
		return nil, errors.New("Error: Some files already exist")
	}
	err := dbf.createFiles()
	if err != nil {
		return nil, err
	}
	return dbf, nil
}

func (dbf *DbFiles) Close() error {
	err := dbf.tree.Close()
	if err != nil {
		return errors.New("Error: Cannot close structure")
	}
	dbf.tree = nil

	err = dbf.data.Close()
	if err != nil {
		return errors.New("Error: Cannot close data")
	}
	dbf.data = nil

	return nil
}

//commit func
func (dbf *DbFiles) Commit() error {
	err := dbf.tree.Sync()
	if err != nil {
		return errors.New("Error: Cannot commit structure")
	}
	err = dbf.data.Sync()
	if err != nil {
		return errors.New("Error: Cannot commit data")
	}
	return nil
}

//write funcs
func (dbf *DbFiles) WriteTreeElem(address int64, elem []byte) (int64, error) {
	if dbf.tree == nil {
		return -1, errors.New("Error: File not opened")
	}
	if address == -1 {
		fs, err := dbf.tree.Stat()
		if err != nil {
			panic("Cannot read stats for file")
		}
		address = fs.Size()
	}

	_, err := dbf.tree.WriteAt(elem, address)
	if err != nil {
		return -1, err
	}

	return address, nil
}

func (dbf *DbFiles) WriteData(data []byte) (int64, error) {
	if dbf.data == nil {
		return -1, errors.New("Error: File not opened")
	}

	fs, err := dbf.data.Stat()
	if err != nil {
		panic("Cannot read stats for opened file")
	}
	address := fs.Size()

	_, err = dbf.data.WriteAt(data, address)
	if err != nil {
		return -1, err
	}
	return address, nil
}

func (dbf *DbFiles) WriteRootAddress(adr []byte) error {
	_, err := dbf.WriteTreeElem(0, adr)
	return err
}

//read funcs
func (dbf *DbFiles) ReadFullTree() ([]byte, error) {
	return readFile(dbf.tree)
}

func (dbf *DbFiles) ReadData(address int64, length int) ([]byte, error) {
	if dbf.data == nil {
		return nil, errors.New("Error: File not opened")
	}

	b := make([]byte, length)
	n, err := dbf.data.ReadAt(b, address)
	if err != nil {
		return nil, err
	}
	if n != length {
		return nil, errors.New("Error: Bad read")
	}

	return b, nil
}

func readFile(file *os.File) ([]byte, error) {
	if file == nil {
		return nil, errors.New("Error: File not opened")
	}

	fs, err := file.Stat()
	if err != nil {
		panic("Cannot read stats for opened file")
	}
	length := fs.Size()

	data := make([]byte, length)

	n, err := file.Read(data)
	if err != nil {
		return nil, err
	}
	if n != len(data) {
		return nil, errors.New("Error: Cannot read data")
	}
	return data, nil
}
