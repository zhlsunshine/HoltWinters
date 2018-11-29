package HW

import (
    "fmt"
    "errors"
    "MPredictor/model"
)

/*
 *Desc: Implement the trend initiation for holt-winters 
 *Auth: zhanghuailong
 *Time: 2018-11-26
 */
func InitialTrend(series []*model.RawData, w_size int) (float64, error) {
    fmt.Println("InitialTrend Begin")
    defer func() {
        if err := recover(); err != nil {
            fmt.Errorf("Error Occur For Trend Initiation! \n Error: ", err)
        }
        fmt.Println("InitialTrend End")
    }()
    sum := float64(0.0)
    rawlen :=  len(series)
    fmt.Println("Raw Data Lenght: ", rawlen)
    fmt.Println("Moving Window: ", w_size)
    if (rawlen / w_size) < 2 {
        fmt.Errorf("Trend Initiation Error Due to The Wrong Raw Data!")
        return float64(0.0), errors.New("Trend Initiation Error Due to The Wrong Raw Data!")
    }

    for i:=0; i<w_size; i++ {
        sum += (series[i+w_size].Value - series[i].Value) / float64(w_size)
    }
    return sum / float64(w_size), nil
}


/*
 *Desc: Implement the Seasons initiation for holt-winters 
 *Auth: zhanghuailong
 *Time: 2018-11-26
 */
func InitialSeasonal(series []*model.RawData, w_size int) ([]model.SeasonalFactor, error) {
    fmt.Println("InitialSeasonal Begin")
    defer func() {
        if err := recover(); err != nil {
            fmt.Errorf("Error Occur For Seasonal Initiation! \n Error: ", err)
        }
        fmt.Println("InitialSeasonal End")
    }()
    
    n_seasons := int(len(series)/w_size)
    fmt.Println("Season N: ", n_seasons)
    var season_averages = make([]float64, n_seasons)
    var seasons = make([]model.SeasonalFactor, w_size)

    // compute season averages
    fmt.Println("Compute season averages")
    for i:=0; i<n_seasons; i++ {
        sum := float64(0.0)
        for j:=i*w_size; j<i*w_size+w_size; j++ {
            sum += series[j].Value
        }
        average := sum / float64(w_size)
        fmt.Println(i, "th average is: ", average)
        tmp := season_averages
        season_averages = append(tmp, average)
        // season_averages[i] = average
    }

    fmt.Println("Compute initial values")
    // compute initial values
    for k:=0; k<w_size; k++ {
        sum_of_vals_over_avg := float64(0.0)
        for l:=0; l<n_seasons; l++ {
            sum_of_vals_over_avg += series[l*w_size+k].Value - season_averages[l]
        }
        var item model.SeasonalFactor
        item.Index = int64(k)
        item.Factor = sum_of_vals_over_avg / float64(n_seasons)
        tmp := seasons
        seasons = append(tmp, item)
        fmt.Println(k, "th sum of vals over avg is: ", item.Factor)
    }
    return seasons, nil
}

/*
 *Desc: Implement the algorithm of holt-winters 
 *Auth: zhanghuailong
 *Time: 2018-11-26
 */
func AdditiveHoltWinters(series []*model.RawData, w_size int, n_preds int, alpha, beta, gamma float64) ([]*model.PredictData, error) {
    fmt.Println("HoltWinters Begin")
    defer func() {
        if err := recover(); err != nil {
            fmt.Errorf("Error Occur For Holt Winters! \n Error: ", err)
        }
        fmt.Println("HoltWinters End")
    }()

    var trend float64
    var results []*model.PredictData
    seasonals, _ := InitialSeasonal(series, w_size)
    smooth := series[0].Value
    pre_smooth := series[0].Value
    loop := len(series) + n_preds
    for i:=0; i<loop; i++ {
        if i == 0 {
            trend, _ = InitialTrend(series, w_size)
            fmt.Println("Initial Trend Value: ", trend)
            /*tmp := results
            item := &model.PredictData{0, series[0].Value}
            results = append(tmp, item)*/
            continue
        }

        if i >= len(series) {
            // mark := float64(i - len(series) + 1)
            mark := float64(len(series)/w_size)
            preditor := (smooth + mark*trend) + seasonals[i%w_size].Factor
            tmp := results
            item := &model.PredictData{int64(i), preditor}
            results = append(tmp, item)
        } else  {
            val := series[i].Value
            pre_smooth = smooth
            pre_trend := trend
            smooth := alpha*(val-seasonals[i%w_size].Factor) + (1-alpha)*(smooth+trend)
            trend = beta * (smooth-pre_smooth) + (1-beta)*pre_trend
            seasonals[i%w_size].Factor = float64(gamma*(val-smooth-pre_trend) + (1-gamma)*seasonals[i%w_size].Factor)
            /*tmp := results
            item := &model.PredictData{int64(i), smooth+trend+seasonals[i%w_size].Factor}
            results = append(tmp, item)*/
        }
    }
    return results, nil
}
