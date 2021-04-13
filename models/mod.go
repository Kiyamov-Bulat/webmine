package models

import (
	"io/ioutil"
	"log"
	"path"
	"time"
)

const MOD_DIR_PATH = "./mods"

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
	file, err := ioutil.ReadFile(path.Join(MOD_DIR_PATH, mod.Name))
	if err != nil {
		return mod
	}
	mod.Content = file
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
	ioutil.WriteFile(path.Join(MOD_DIR_PATH, mod.Name), mod.Content, 0644)
}

func (mod *Mod) Create() {
	GetDB().Create(mod)
}

func GetAllModNames() []Mod {
	mods := []Mod{}
	getAllItems("mods", &mods)
	return mods
}

func (mod *Mod) SetContent() {
	files, err := ioutil.ReadDir(MOD_DIR_PATH)
	if err != nil {
		log.Println("The setting of mod content error!")
	}
	for _, file := range files {
		if mod.Name == file.Name() && mod.ModTime == file.ModTime() && mod.Size == file.Size() {
			tmp, err := ioutil.ReadFile(mod.Name)
			if err != nil {
				log.Println("File read error!")
			}
			mod.Content = tmp
			return
		}
	}
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
