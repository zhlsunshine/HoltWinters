package model

import (
    "sync"
    "HoltWinters/utils"
)

type RawData struct {
    Value     float64    `json:"value"`
    Label     string     `json:"label"`
    Timestamp int64      `json:"timestamp"`
}

type SeasonalFactor struct {
    Index     int64      `json:"index"` 
    Factor    float64    `json:"factor"`
}

type PredictData struct {
    Index     int64      `json:"index"` 
    Value     float64    `json:"value"`
}

type TrainProp struct {
    Precision   int      `json:"precision"`
    TrainStep   float64  `json:"trainstep"`
    WindowSize  int      `json:"windowsize"`
    PredictSize int      `json:"predictsize"`
    TrainMode   string   `json:"trainmode"`
}

type HWParam struct {
    Alpha    float64
    Beta     float64
    Gamma    float64
    SSEP     float64
}

// Configuration struct define
type ConStruct struct {
    Dbconnection struct {
        DBName    string `json:"dbname"`
        DBAddress string `json:"dbaddress"`
        UserName  string `json:"username"`
        PassWord  string `json:"password"`
        Queries   []struct {
            Name string `json:"name"`
            Iql  string `json:"iql"`
        } `json:"queries"`
    } `json:"dbconnection"`
    Endpoint struct {
        APIAddr  string `json:"apiaddress"`
        UserName string `json:"username"`
        PassWord string `json:"password"`
        Queries  []struct {
            Name string `json:"name"`
            Iql  string `json:"iql"`
        } `json:"queries"`

    } `json:"endpoint"`
    FileData struct {
        FileDir  string `json:"filedir"`
        DataName string `json:"dataname"`
    } `json:"filedata"`
    Training TrainProp `json:"training"`
}

// Concurrency for training model
const MAX_ROUTINE_NUM = 100000
// Training times for the Nelder Mead Method
const MAX_TRAINING_TIMES = 100000
// Repeat training times for the Nelder Mead Method
const REPEAT_TRAINING_TIMES = 10000
// Handle the infinite value 
const PRECISE_FOR_INF = 10000

var RoutinePool = make(chan int64, MAX_ROUTINE_NUM)

// The best model parameter for the current training
var HWPInstance HWParam

// History data for predict
var HisData []*RawData

/*lock for sync codes*/
var Lock = new(sync.Mutex)

//model file directory
var MFDir = utils.GetCurrPath()
var MFName = utils.GetCurrPath() + "HWModel"

// Configuration instance
var CInstance ConStruct 

// Configuration file directory
var ConfDir string
// Log Configuration file directory
var LogConf string
