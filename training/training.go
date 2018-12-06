package training

import (
    "fmt"
    "sync"
    "encoding/json"
    "HoltWinters/model"
    "HoltWinters/utils"
    "HoltWinters/holt-winters"
)

/*
 *Desc: Training Controller for the whole training process
 *Auth: zhanghuailong
 *Time: 2018-11-28
 */
func TrainingController(series []*model.RawData, r_series []*model.RawData, trainP *model.TrainProp) {
    defer func() {
        if err := recover(); err != nil {
            fmt.Errorf("TrainingController error: %v", err)
        }
    }()

    fmt.Println("Begin Training Process")
    var waitgroup = new(sync.WaitGroup)
    PLoop := utils.Powerf(10, trainP.Precision)
    loop := int(PLoop) 
    alpha, beta, gamma := float64(0), float64(0), float64(0)
    for ai:=1; ai<=loop; ai++ {
        alpha = float64(ai) / float64(PLoop)
        for bi:=1; bi<=loop; bi++ {
            beta = float64(bi) / float64(PLoop)
            for gi:=1; gi<=loop; gi++ {
                gamma = float64(gi) / float64(PLoop)
                model.RoutinePool <- int64(float64(gi)*PLoop*PLoop*PLoop + float64(bi)*PLoop*PLoop + float64(ai)*PLoop)
                waitgroup.Add(1)
                go TrainJob(series, r_series, trainP.WindowSize, trainP.PredictSize, waitgroup, alpha, beta, gamma)
            }
        } 
    }
    waitgroup.Wait()
}

/*
 *Desc: Training Job for each time, it's a go routine
 *Auth: zhanghuailong
 *Time: 2018-11-28
 */
func TrainJob(series []*model.RawData, r_series []*model.RawData, w_size int, n_preds int, wg *sync.WaitGroup, alpha, beta, gamma float64) {
    defer func() {
        if err := recover(); err != nil {
            fmt.Errorf("TrainJob error: %v", err)
        }
        wg.Done()
        times := <-model.RoutinePool
        fmt.Println(times, "th training model")
    }()
    pData, _ := HW.AdditiveHoltWinters(series, w_size, n_preds, alpha, beta, gamma)
    FittingPredictData(r_series, pData, alpha, beta, gamma)
}

/*
 *Desc: Fitting produce for predict data 
 *Auth: zhanghuailong
 *Time: 2018-11-28
 */
func FittingPredictData(r_series []*model.RawData, p_series []*model.PredictData, alpha, beta, gamma float64) float64 {
    // fmt.Println("Fitting Data Begin")
    if len(r_series) != len(p_series) {
        fmt.Errorf("Data Fitting Error! Real series length:[", len(r_series), "]  Predicted series length:[", len(p_series), "]" )
        return float64(0)
    }
    s_minu := float64(0)
    variance := float64(0)
     
    for i:=0; i<len(r_series); i++ {
        s_minu = p_series[i].Value - r_series[i].Value
        variance += utils.Powerf(s_minu, 2) 
    }
    if variance < model.HWPInstance.SSEP {
        model.Lock.Lock()
        model.HWPInstance.SSEP = variance
        model.HWPInstance.Alpha = alpha
        model.HWPInstance.Beta = beta
        model.HWPInstance.Gamma = gamma
        jPackage, _ := json.Marshal(&model.HWPInstance)
        utils.SaveData(jPackage, model.MFName)
        model.Lock.Unlock()
        fmt.Println("Current Model is: A:[", model.HWPInstance.Alpha, "]   B:[", model.HWPInstance.Beta, "]   G:[", model.HWPInstance.Gamma, "],  SSEP value is: [", model.HWPInstance.SSEP, "]")
    }
    fmt.Println("###########################################################################################")
    fmt.Println("A:[", alpha, "]   B:[", beta, "]   G:[", gamma, "]")
    fmt.Println("Predicted Data Series: ")
    for index, item := range p_series {
        fmt.Println("the", index, "th predicted data is ", *item)
    }
    fmt.Println("The SSEP (Sum of Squared Errors for Prediction) For the Prediction is: [", variance, "]")
    fmt.Println("\n")
    return variance
    // fmt.Println("Fitting Data End")
}
