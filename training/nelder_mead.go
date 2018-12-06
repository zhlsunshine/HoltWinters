package training 

import (
    "fmt"
    "math"
    "sort"
    "encoding/json"
    "HoltWinters/model"
    "HoltWinters/utils"
    "HoltWinters/holt-winters"
)

// Note: Alpha , Gamma, Rho and Sigma are respectively the reflection, expansion, contraction and shrink coefficients. Standard values are Alpha = 1, Gamma = 2, Rho = 0.5 and Sigma = 0.5. 
// Constraint: Alpha > 0, 0 < Rho < 1, Gamma > 1, Gamma > Alpha, 0 < Sigma < 1.

/*
 *Desc: Initiate the simplex points for the Nelder-Mead Method
 *Auth: zhanghuailong
 *Time: 2018-11-30
 */
func InitSimplexPoints() *model.Simplex {
    fmt.Println("InitSimplexPoints Begin")
    var rSimplex = new(model.Simplex)
    rSimplex.UnitStep = float64(0.001)
    rSimplex.UnitVector = float64(1)
    bLeftParam := rSimplex.UnitVector / (float64(model.VerticsNUM) * math.Sqrt(2)) 
    bRightParam := math.Sqrt(float64(model.VerticsNUM + 1)) - 1
    bParam := bLeftParam * bRightParam
    aParam := bParam + rSimplex.UnitVector / math.Sqrt(2)

    // Initial the three vertics
    X0 := &model.SimplexVertic{float64(0), float64(0), float64(0), float64(0)}
    X1 := &model.SimplexVertic{aParam, bParam, bParam, float64(0)}
    X2 := &model.SimplexVertic{bParam, aParam, bParam, float64(0)}
    X3 := &model.SimplexVertic{bParam, bParam, aParam, float64(0)}

    rSimplex.Vertics = append(rSimplex.Vertics, *X0)
    rSimplex.Vertics = append(rSimplex.Vertics, *X1)
    rSimplex.Vertics = append(rSimplex.Vertics, *X2)
    rSimplex.Vertics = append(rSimplex.Vertics, *X3)

    fmt.Println("The Initial Simplex is : ", *rSimplex)
    return rSimplex
}

/*
 *Desc: Order vertics values for simplex
 *Auth: zhanghuailong
 *Time: 2018-11-30
 */
func OrderVerticsValue(rSimplex *model.Simplex, series []*model.RawData, r_series []*model.RawData, trainP *model.TrainProp) *model.Simplex {
    fmt.Println("OrderVerticsValue Begin")
    for index, item := range rSimplex.Vertics {
        if rSimplex.Vertics[index].VerticValue == 0 {
            pData, _ := HW.AdditiveHoltWinters(series, trainP.WindowSize, trainP.PredictSize, item.Alpha, item.Beta, item.Gamma)
            fitValue := Fitting4NelderMead(r_series, pData)
            rSimplex.Vertics[index].VerticValue = fitValue
        }
    } 
    sort.Sort(rSimplex.Vertics)
    fmt.Println("Sorted simplex vertics: ", *rSimplex)
    return rSimplex
}

/*
 *Desc: Calculate the centroid for simplex
 *Auth: zhanghuailong
 *Time: 2018-11-30
 */
func GetCentroid(rSimplex *model.Simplex) *model.SimplexVertic {
    fmt.Println("GetCentroid Begin")
    sumA := float64(0)
    sumB := float64(0)
    sumG := float64(0)
    for i:=0; i<rSimplex.Vertics.Len() - 2; i++ {
        sumA += rSimplex.Vertics[i].Alpha
        sumB += rSimplex.Vertics[i].Beta
        sumG += rSimplex.Vertics[i].Gamma
    } 

    return &model.SimplexVertic{sumA/float64(model.VerticsNUM), sumB/float64(model.VerticsNUM), sumG/float64(model.VerticsNUM), 0}
}

/*
 *Desc: Fitting produce for predict data 
 *Auth: zhanghuailong
 *Time: 2018-11-30
 */
func Fitting4NelderMead(r_series []*model.RawData, p_series []*model.PredictData) float64 {
    // fmt.Println("Fitting4NelderMead Begin")
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
    return variance
    // fmt.Println("Fitting Data End")
}

/*
 *Desc: Calculate the Reflection Point for simplex 
 *Auth: zhanghuailong
 *Time: 2018-12-3
 */
func CalReflectionPoint(rSimplex *model.Simplex, centPoint *model.SimplexVertic, Alpha float64) *model.SimplexVertic {
    fmt.Println("CalReflectionPoint Start")
    endI   := rSimplex.Vertics.Len() - 1
    
    rPoint := centPoint.Sub(&rSimplex.Vertics[endI])
    rPoint = rPoint.Mul(Alpha)
    rPoint = centPoint.Add(rPoint)

    return rPoint.Abs() 
}

/*
 *Desc: Calculate the Expansion Point for simplex 
 *Auth: zhanghuailong
 *Time: 2018-12-3
 */
func CalExpansionPoint(rSimplex *model.Simplex, rPoint, centPoint *model.SimplexVertic, Alpha, Gamma float64) *model.SimplexVertic {
    fmt.Println("CalExpansionPoint Start")

    ePoint := rPoint.Sub(centPoint)
    ePoint = ePoint.Mul(Gamma)
    ePoint = centPoint.Add(ePoint)

    return ePoint.Abs()
}

/*
 *Desc: Calculate the Contraction Point for simplex 
 *Auth: zhanghuailong
 *Time: 2018-12-3
 */
func CalContractionPoint(rSimplex *model.Simplex, centPoint *model.SimplexVertic, Rho float64, isOut bool) *model.SimplexVertic {
    fmt.Println("CalContractionPoint Start")
    endI   := rSimplex.Vertics.Len() - 1

    cPoint := (&rSimplex.Vertics[endI]).Sub(centPoint)
    cPoint = cPoint.Mul(Rho) 
    if isOut {
        cPoint = centPoint.Add(cPoint)
    } else {
        cPoint = centPoint.Sub(cPoint)
    }
    
    return cPoint.Abs()
}

/*
 *Desc: Shrink All Points for simplex 
 *Auth: zhanghuailong
 *Time: 2018-12-3
 */
func ShrinkAllPoints(rSimplex *model.Simplex, Sigma float64) *model.Simplex {
    fmt.Println("ShrinkAllPoints Start")
    startI := 0    

    for i:=1; i<rSimplex.Vertics.Len(); i++ {
        sPoint := (&rSimplex.Vertics[i]).Sub(&rSimplex.Vertics[startI])
        sPoint = sPoint.Mul(Sigma)
        sPoint = (&rSimplex.Vertics[startI]).Add(sPoint)
        sPoint = sPoint.Abs()

        rSimplex.Vertics[i].Alpha = sPoint.Alpha
        rSimplex.Vertics[i].Beta = sPoint.Beta
        rSimplex.Vertics[i].Gamma = sPoint.Gamma
    }
    
    return rSimplex
}

/*
 *Desc: Training By Nelder-Mead Method 
 *Auth: zhanghuailong
 *Time: 2018-12-3
 */
func NelderMeadTraining(series []*model.RawData, r_series []*model.RawData, trainP *model.TrainProp) {
    defer func() {
        if err := recover(); err != nil {
            fmt.Errorf("NelderMeadTraining error: %v", err)
        }
    }()
    fmt.Println("NelderMeadTraining Begin")
    
    // initate the whole simplex
    rSimplex := InitSimplexPoints()
    fmt.Println("Initiate simplex points: ", *rSimplex)
    preCentPoint := &model.SimplexVertic{float64(0), float64(0), float64(0), float64(0)} 

    for i:=0; i<model.MAX_TRAINING_TIMES; i++ {
        isShrink := false 
        startI := 0
        rSimplex = OrderVerticsValue(rSimplex, series, r_series, trainP)
        centPoint := GetCentroid(rSimplex)
        fmt.Println("Centroid Point: ", *centPoint)
        // Search the optimal model as soon as possible
        if preCentPoint.Cmp(centPoint) == 0 {
            centPoint.Alpha = centPoint.Alpha + rSimplex.UnitStep
            centPoint.Beta = centPoint.Beta + rSimplex.UnitStep
            centPoint.Gamma = centPoint.Gamma + rSimplex.UnitStep
        }
        endI   := rSimplex.Vertics.Len() - 1
        rPoint := CalReflectionPoint(rSimplex, centPoint, model.Alpha)
        pData, _ := HW.AdditiveHoltWinters(series, trainP.WindowSize, trainP.PredictSize, rPoint.Alpha, rPoint.Beta, rPoint.Gamma)
        fitValue := Fitting4NelderMead(r_series, pData)
        rPoint.VerticValue = fitValue
        fmt.Println("Reflection Point: ", *rPoint, "   rSimplex best Vertics: ", rSimplex.Vertics[startI], "   rSimplex worst Vertics: ", rSimplex.Vertics[endI])

        if rPoint.VerticValue >= rSimplex.Vertics[startI].VerticValue {
            if rPoint.VerticValue < rSimplex.Vertics[endI-1].VerticValue {
                rSimplex.Vertics[endI].Alpha = rPoint.Alpha
                rSimplex.Vertics[endI].Beta = rPoint.Beta
                rSimplex.Vertics[endI].Gamma = rPoint.Gamma
                rSimplex.Vertics[endI].VerticValue = rPoint.VerticValue
                isShrink = true
            }
        } 

        if rPoint.VerticValue < rSimplex.Vertics[startI].VerticValue {
            ePoint := CalExpansionPoint(rSimplex, rPoint, centPoint, model.Alpha, model.Gamma)
            pData, _ := HW.AdditiveHoltWinters(series, trainP.WindowSize, trainP.PredictSize, ePoint.Alpha, ePoint.Beta, ePoint.Gamma)
            fitValue := Fitting4NelderMead(r_series, pData)
            ePoint.VerticValue = fitValue
            fmt.Println("Expansion Point: ", *ePoint)
            if ePoint.VerticValue < rPoint.VerticValue {
                rSimplex.Vertics[endI].Alpha = ePoint.Alpha
                rSimplex.Vertics[endI].Beta = ePoint.Beta
                rSimplex.Vertics[endI].Gamma = ePoint.Gamma
                rSimplex.Vertics[endI].VerticValue = ePoint.VerticValue
            } else {
                rSimplex.Vertics[endI].Alpha = rPoint.Alpha
                rSimplex.Vertics[endI].Beta = rPoint.Beta
                rSimplex.Vertics[endI].Gamma = rPoint.Gamma
                rSimplex.Vertics[endI].VerticValue = rPoint.VerticValue
            }
            isShrink = true
        }

        if rPoint.VerticValue >= rSimplex.Vertics[endI-1].VerticValue {
            if rPoint.VerticValue < rSimplex.Vertics[endI].VerticValue {
                cPoint := CalContractionPoint(rSimplex, centPoint, model.Rho, true)
                pData, _ := HW.AdditiveHoltWinters(series, trainP.WindowSize, trainP.PredictSize, cPoint.Alpha, cPoint.Beta, cPoint.Gamma)
                fitValue := Fitting4NelderMead(r_series, pData)
                cPoint.VerticValue = fitValue
                fmt.Println("Out Contraction Point: ", *cPoint)
                if cPoint.VerticValue < rSimplex.Vertics[endI].VerticValue {
                    rSimplex.Vertics[endI].Alpha = cPoint.Alpha
                    rSimplex.Vertics[endI].Beta = cPoint.Beta
                    rSimplex.Vertics[endI].Gamma = cPoint.Gamma
                    rSimplex.Vertics[endI].VerticValue = cPoint.VerticValue
                    isShrink = true
                }
            }
        }

        if rPoint.VerticValue >= rSimplex.Vertics[endI].VerticValue {
                cPoint := CalContractionPoint(rSimplex, centPoint, model.Rho, false)
                pData, _ := HW.AdditiveHoltWinters(series, trainP.WindowSize, trainP.PredictSize, cPoint.Alpha, cPoint.Beta, cPoint.Gamma)
                fitValue := Fitting4NelderMead(r_series, pData)
                cPoint.VerticValue = fitValue
                fmt.Println("In Contraction Point: ", *cPoint)
                if cPoint.VerticValue < rSimplex.Vertics[endI].VerticValue {
                    rSimplex.Vertics[endI].Alpha = cPoint.Alpha
                    rSimplex.Vertics[endI].Beta = cPoint.Beta
                    rSimplex.Vertics[endI].Gamma = cPoint.Gamma
                    rSimplex.Vertics[endI].VerticValue = rPoint.VerticValue
                    isShrink = true
                }
        }


        if !isShrink {
            rSimplex = ShrinkAllPoints(rSimplex, model.Sigma)
        }
        preCentPoint = centPoint

        /*preCentPoint.Alpha = centPoint.Alpha
        preCentPoint.Beta = centPoint.Beta
        preCentPoint.Gamma = centPoint.Gamma
        preCentPoint.VerticValue = centPoint.VerticValue*/

        jPackage, _ := json.Marshal(rSimplex)
        utils.SaveData(jPackage, model.MFName)
        fmt.Println("Current Simplex Model is: ", *rSimplex) 
    }
}
