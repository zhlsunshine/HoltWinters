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
    Value     float64    `json:"factor"`
}

type TrainProp struct {
    Precision   int
    TrainStep   float64
    WindowSize  int
    PredictSize int
}

type HWParam struct {
    Alpha    float64
    Beta     float64
    Gamma    float64
    SSEP     float64
}

// Concurrency for training model
const MAX_ROUTINE_NUM = 100000
// Training times for the Nelder Mead Method
const MAX_TRAINING_TIMES = 1000000
var RoutinePool = make(chan int64, MAX_ROUTINE_NUM)

// The best model parameter for the current training
var HWPInstance HWParam

/*lock for sync codes*/
var Lock = new(sync.Mutex)

//model file directory
var MFDir = utils.GetCurrPath()
var MFName = utils.GetCurrPath() + "HWModel"
