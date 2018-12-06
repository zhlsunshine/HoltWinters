package main

import (
    "fmt"
    "HoltWinters/model"
    "HoltWinters/utils"
    "HoltWinters/training"
    "HoltWinters/holt-winters"
)


func main() {
    defer func() {
        if err := recover(); err != nil {
            fmt.Errorf("Internal error: %v", err)
        }
    }()
    fmt.Println("Execute ConfigureLogger")
    model.HWPInstance.SSEP = float64(9223372036854775807)
    utils.CheckMDir(model.MFDir, model.MFName) 
    var trainP = new(model.TrainProp)
    trainP.Precision = 3
    trainP.TrainStep = float64(0.001) 
    trainP.WindowSize = 12
    trainP.PredictSize = 12
    var rawdata [109]float64 = [109]float64{30,21,29,31,40,48,53,47,37,39,31,29,17,9,20,24,27,35,41,38,27,31,27,26,21,13,21,18,33,35,40,36,22,24,21,20,17,14,17,19,26,29,40,31,20,24,18,26,17,9,17,21,28,32,46,33,23,28,22,27,18,8,17,21,31,34,44,38,31,30,26,32,48,53,47,37,39,31,29,17,9,20,24,27,35,31,20,24,18,17,21,28,32,46,33,23,28,22,27,18,8,17,21,33,23,28,22,27,18}

    var series []*model.RawData
    for i:=0; i<109; i++ {
        series = append(series, &model.RawData{float64(rawdata[i]), "", 0})
    } 
    fmt.Println("Series: ", series)
    alpha, beta, gamma := float64(0.716), float64(0.029), float64(0.993)
    pData, _ := HW.AdditiveHoltWinters(series, trainP.WindowSize, trainP.PredictSize, alpha, beta, gamma)
    for index, item := range pData {
        fmt.Println(index, "th predicted data is ", *item)
    }

    var o_data [97]float64 = [97]float64{30,21,29,31,40,48,53,47,37,39,31,29,17,9,20,24,27,35,41,38,27,31,27,26,21,13,21,18,33,35,40,36,22,24,21,20,17,14,17,19,26,29,40,31,20,24,18,26,17,9,17,21,28,32,46,33,23,28,22,27,18,8,17,21,31,34,44,38,31,30,26,32,48,53,47,37,39,31,29,17,9,20,24,27,35,31,20,24,18,17,21,28,32,46,33,23,28}
    var s_o_data []*model.RawData
    for  i:=0; i<97; i++ {
        s_o_data = append(s_o_data, &model.RawData{float64(o_data[i]), "", 0})
    }

    var r_data [12]float64 = [12]float64{22,27,18,8,17,21,33,23,28,22,27,18}
    var r_o_data []*model.RawData
    for  i:=0; i<12; i++ {
        r_o_data = append(r_o_data, &model.RawData{float64(r_data[i]), "", 0})
    }
 
    //training.TrainingController(s_o_data, r_o_data, trainP)
    training.NelderMeadTraining(s_o_data, r_o_data, trainP)
}
