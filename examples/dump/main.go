package main

func main() {
	FromInts(1,2,3,4,5).Last().MapFloat64(func (v int) float64 {
		return float64(v) * 1.65
	}).ToOne();

}