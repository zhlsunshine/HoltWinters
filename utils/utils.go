package utils

import (
    "os"
    "fmt"
    "os/exec"
    "io/ioutil"
    "path/filepath"
)


func Powerf(x float64, n int) float64 {
    ans := 1.0 
    for n != 0 { 
        if n%2 == 1 { 
            ans *= x 
        } 
        x *= x 
        n /= 2 
    } 
    return ans 
}

func SaveData(data []byte, filedir string) {
    // fmt.Println("data: ", data)
    // jstoreusr, _ := json.Marshal(data)
    ioutil.WriteFile(filedir, data, 0777)
}

func CheckMDir(fileFolder string, filePath string) bool {
    if _, err := os.Stat(fileFolder); os.IsNotExist(err) {
        if err != nil {
            fmt.Errorf("create user data directory: %s error!+%v\n", fileFolder, err)
        }
        fmt.Printf("create user data directory: %s\n", fileFolder)
        os.Mkdir(fileFolder, 0777)
        return false
    }

    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        fmt.Printf("create user data file: %s\n", filePath)
        os.Create(filePath)
        return false
    }
    return true
}

func GetCurrPath() string {
    file, _ := exec.LookPath(os.Args[0])
    path, _ := filepath.Abs(file)
    dir, _ := filepath.Split(path)
    return dir
}
