package model

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

// Concurrency for training model
const MAX_ROUTINE_NUM = 100000
var RoutinePool = make(chan int64, MAX_ROUTINE_NUM)
