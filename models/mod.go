package models

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

const MOD_DIR_PATH = "./data/mods"
const MOD_ZIP_NAME = "mods.zip"
const MOD_ZIP_PATH = "./tmp/" + MOD_ZIP_NAME

type Mod struct {
	MineModel   `gorm:"embedded"`
	Name        string    `json:"name"`
	Content     []byte    `json:"file_content" sql:"-"`
	Size        int64     `json:"size"`
	ModTime     time.Time `json:"mod_time"`
	Description string    `json:"description"`
}

func GetMod(id uint) *Mod {
	mod := &Mod{}
	//setItemFromDB("mods", mod, mod.ID)
	GetDB().Where("id = ?", id).First(mod)
	mod.SetContent()
	return mod
}

func GetAllMods() []Mod {
	mods := []Mod{}
	getAllItems("mods", &mods)
	for _, mod := range mods {
		mod.SetContent()
	}
	return mods
}

func (mod *Mod) SaveFile() {
	ioutil.WriteFile(mod.GetFullFilePath(), mod.Content, 0644)
}

func (mod *Mod) Create() {
	GetDB().Create(mod)
}

func (mod *Mod) Delete() {
	os.Remove(mod.GetFullFilePath())
	GetDB().Delete(mod)
}

func (mod *Mod) SetContent() {
	file, err := ioutil.ReadFile(mod.GetFullFilePath())
	if err != nil {
		log.Println("File read error!")
	}
	mod.Content = file
}

func (mod *Mod) GetFullFilePath() string {
	return path.Join(MOD_DIR_PATH, mod.Name)
}

func GetAllModNames() []Mod {
	mods := []Mod{}
	getAllItems("mods", &mods)
	return mods
}

func initModes() {
	files, err := ioutil.ReadDir(MOD_DIR_PATH)
	if err != nil {
		log.Println("It's error", err)
	}
	mods := make([]Mod, len(files))
	for i, file := range files {
		if !file.IsDir() {
			mods[i].Name = file.Name()
			mods[i].ModTime = file.ModTime()
			mods[i].Size = file.Size()
			GetDB().Where(&mods[i]).First(&mods[i])
			if !mods[i].IsValid() {
				mods[i].Create()
			}
		}
	}
	db := GetDB()
	for _, mod := range mods {
		db = db.Not(mod)
	}
	db.Delete(&Mod{})
}
