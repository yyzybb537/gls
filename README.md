# gls
Goroutine local storage library

# Benchmark

| Operate | Performance |
| ------| ------ |
| Benchmark_Goid | 5.30 ns/op |
| Benchmark_Set | 271 ns/op |
| Benchmark_Get | 155 ns/op |
| Benchmark_Set_4Threads | 129 ns/op |
| Benchmark_Get_4Threads | 58.6 ns/op |
| Benchmark_GetNil_4Threads | 48.1 ns/op |
