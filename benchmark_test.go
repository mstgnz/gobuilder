package gobuilder

import (
	"testing"
)

var benchmarkGb = NewGoBuilder(Postgres)

func BenchmarkSelect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkGb.Table("users").Select().Where("id", "=", i).Sql()
	}
}

func BenchmarkInsert(b *testing.B) {
	args := map[string]any{"firstname": "John", "lastname": "Doe"}
	for i := 0; i < b.N; i++ {
		benchmarkGb.Table("users").Create(args).Sql()
	}
}

func BenchmarkUpdate(b *testing.B) {
	args := map[string]any{"firstname": "Jane"}
	for i := 0; i < b.N; i++ {
		benchmarkGb.Table("users").Update(args).Where("id", "=", i).Sql()
	}
}

func BenchmarkDelete(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkGb.Table("users").Delete().Where("id", "=", i).Sql()
	}
}
