package utils

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
