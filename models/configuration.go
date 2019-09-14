package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"projectionist/consts"
	"strconv"
)

const (
	FileIDNAMEPtrn = "%d|%s.json"
	FilePath       = "%s/%s.json"
	FilePathPtrn   = "%s/" + FileIDNAMEPtrn
)

type Configuration struct {
	ID      int                    `json:"id"`
	Name    string                 `json:"username"`
	Config  map[string]interface{} `json:"config"`
	Deleted int                    `json:"deleted"`
}

func (c *Configuration) SetDBCtx(interface{}) error {

	return nil
}

func (c *Configuration) Validate() error {
	if c.ID == 0 {
		return fmt.Errorf("id must be not 0")
	}

	if c.Name == "" {
		return fmt.Errorf("name must be not empty")
	}

	if c.Config == nil {
		return fmt.Errorf("config field must be not empty")
	}

	return nil
}

func (c *Configuration) IsExistByName() (error, bool) {
	var matched, err = filepath.Glob(c.Name)
	if err != nil {
		return err, false
	}
	if len(matched) == 0 {
		return consts.ErrNotFound, false
	}

	return nil, true
}

func (c *Configuration) Count() (int, error) {
	var count int
	var err = filepath.Walk(consts.PathSaveCfgs, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		count++

		return nil
	})

	return count, err
}

func (c *Configuration) Save() error {
	var savePath = fmt.Sprintf(FilePathPtrn, consts.PathSaveCfgs, c.ID, c.Name)
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

func (c *Configuration) GetByName(name string) error {
	bfile, err := ioutil.ReadFile(fmt.Sprintf(FilePath, consts.PathSaveCfgs, name))
	if err != nil {
		return err
	}

	return json.Unmarshal(bfile, &c)
}

func (c *Configuration) GetByID(id int64) error {
	var matched, err = filepath.Glob(fmt.Sprintf(FilePath, consts.PathSaveCfgs, strconv.Itoa(int(id))))
	if err != nil {
		return err
	}

	if len(matched) == 0 {
		return consts.ErrNotFound
	}
	var filePath = matched[0]

	if filepath.Ext(filePath) != "json" {
		return fmt.Errorf("file %s extension not .json", filePath)
	}

	bfile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var cfg = Configuration{}

	err = json.Unmarshal(bfile, &cfg)
	if err != nil {
		return err
	}

	if int64(cfg.ID) != id {
		return consts.ErrNotFound
	}

	if cfg.Deleted > 0 {
		return fmt.Errorf("this file id:%d deleted", id)
	}

	c = &cfg

	return nil
}

func (c *Configuration) Pagination(start, stop int) ([]Model, error) {
	var result []Model
	var err = filepath.Walk(consts.PathSaveCfgs, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != "json" {
			return nil
		}

		bfile, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}

		var cfg = Configuration{}

		err = json.Unmarshal(bfile, &cfg)
		if err != nil {
			return err
		}

		if cfg.Deleted > 0 {
			return nil
		}

		if cfg.ID >= start && cfg.ID < stop {
			result = append(result, &cfg)
		}

		return nil
	})
	return result, err
}

func (c *Configuration) Update(id int) error {
	var matched, err = filepath.Glob(fmt.Sprintf(FilePath, consts.PathSaveCfgs, strconv.Itoa(id)))
	if err != nil {
		return err
	}

	if len(matched) == 0 {
		return consts.ErrNotFound
	}
	var filePath = matched[0]

	if filepath.Ext(filePath) != "json" {
		return fmt.Errorf("file %s extension not .json", filePath)
	}

	var updatePath = fmt.Sprintf(FilePathPtrn, consts.PathSaveCfgs, id, c.Name)
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f, err := os.Create(updatePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return os.Remove(filePath)
}

func (c *Configuration) Delete(id int) error {
	var matched, err = filepath.Glob(fmt.Sprintf(FilePath, consts.PathSaveCfgs, strconv.Itoa(id)))
	if err != nil {
		return err
	}

	if len(matched) == 0 {
		return consts.ErrNotFound
	}
	var filePath = matched[0]

	if filepath.Ext(filePath) != "json" {
		return fmt.Errorf("file %s extension not .json", filePath)
	}

	return os.Remove(filePath)
}

func (c *Configuration) GetID() int {
	return c.ID
}

func (c *Configuration) SetID(id int) {
	c.ID = id
}

func (c *Configuration) GetName() string {
	return c.Name
}

func (c *Configuration) SetDeleted() {
	c.Deleted = 1
}
