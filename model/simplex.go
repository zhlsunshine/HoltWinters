package model

import (
    "math"
)

type SimplexVertic struct {
    Alpha        float64
    Beta         float64
    Gamma        float64
    VerticValue  float64
}

type SimplexVertics []SimplexVertic

//Len()
func (s SimplexVertics) Len() int {
    return len(s)
}

//Less():
func (s SimplexVertics) Less(i, j int) bool {
    return s[i].VerticValue < s[j].VerticValue
}

//Swap()
func (s SimplexVertics) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

type Simplex struct {
    Vertics      SimplexVertics 
    UnitStep     float64
    UnitVector   float64
}

const VerticsNUM = 4 
const Alpha = float64(1)
const Gamma = float64(2)
const Rho   = float64(0.5)
const Sigma = float64(0.5)


// Abs returns the vector with nonnegative components.
func (v *SimplexVertic) Abs() *SimplexVertic { return &SimplexVertic{math.Abs(v.Alpha), math.Abs(v.Beta), math.Abs(v.Gamma), float64(0)} }

// Add returns the standard vector sum of v and ov.
func (v *SimplexVertic) Add(ov *SimplexVertic) *SimplexVertic { return &SimplexVertic{v.Alpha + ov.Alpha, v.Beta + ov.Beta, v.Gamma + ov.Gamma, float64(0)} }

// Sub returns the standard vector difference of v and ov.
func (v *SimplexVertic) Sub(ov *SimplexVertic) *SimplexVertic { return &SimplexVertic{v.Alpha - ov.Alpha, v.Beta - ov.Beta, v.Gamma - ov.Gamma, float64(0)} }

// Mul returns the standard scalar product of v and m.
func (v *SimplexVertic) Mul(m float64) *SimplexVertic { return &SimplexVertic{m * v.Alpha, m * v.Beta, m * v.Gamma, float64(0)} }

// Dot returns the standard dot product of v and ov.
func (v *SimplexVertic) Dot(ov *SimplexVertic) float64 { return v.Alpha*ov.Alpha + v.Beta*ov.Beta + v.Gamma*ov.Gamma }

// Cross returns the standard cross product of v and ov.
func (v *SimplexVertic) Cross(ov *SimplexVertic) *SimplexVertic {
    return &SimplexVertic{
                v.Beta*ov.Gamma - v.Gamma*ov.Beta,
                v.Gamma*ov.Alpha - v.Alpha*ov.Gamma,
                v.Alpha*ov.Beta - v.Beta*ov.Alpha,
                float64(0),
    }
}

// Norm returns the vector's norm.
func (v *SimplexVertic) Norm() float64 { return math.Sqrt(v.Dot(v)) }

// Distance returns the Euclidean distance between v and ov.
func (v *SimplexVertic) Distance(ov *SimplexVertic) float64 { return v.Sub(ov).Norm() }

//   -1 if v <  ov
//    0 if v == ov
//   +1 if v >  ov
// values they are lexicographically equal.
func (v *SimplexVertic) Cmp(ov *SimplexVertic) int {
    if v.Alpha < ov.Alpha {
        return -1
    }
    if v.Alpha > ov.Alpha {
        return 1
    }
    // First elements were the same, try the next.
    if v.Beta < ov.Beta {
        return -1
    }
    if v.Beta > ov.Beta {
        return 1
    }
    // Second elements were the same return the final compare.
    if v.Gamma < ov.Gamma {
        return -1
    }
    if v.Gamma > ov.Gamma {
        return 1
    }
    // Both are equal
    return 0
}
