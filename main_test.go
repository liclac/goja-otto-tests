package main

import (
	"github.com/dop251/goja"
	"github.com/robertkrimen/otto"
	"testing"
)

const filename = "test.js"
const src = `
var num = Math.random();
(function() {
	num += Math.random();
	return num;
});
`

func BenchmarkLoading(b *testing.B) {
	b.Run("otto", func(b *testing.B) {
		b.Run("Run", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				vm := otto.New()

				for pb.Next() {
					_, err := vm.Run(src)
					if err != nil {
						b.Error(err)
						return
					}
				}
			})
		})
		b.Run("Eval", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				vm := otto.New()

				for pb.Next() {
					_, err := vm.Eval(src)
					if err != nil {
						b.Error(err)
						return
					}
				}
			})
		})
		b.Run("Compile", func(b *testing.B) {
			baseVM := otto.New()

			script, err := baseVM.Compile(filename, src)
			if err != nil {
				b.Error(err)
				return
			}

			b.ResetTimer()

			b.Run("Run", func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					vm := baseVM.Copy()

					for pb.Next() {
						_, err := vm.Run(script)
						if err != nil {
							b.Error(err)
							return
						}
					}
				})
			})
		})
	})
	b.Run("goja", func(b *testing.B) {
		b.Run("RunScript", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				vm := goja.New()

				for pb.Next() {
					_, err := vm.RunScript(filename, src)
					if err != nil {
						b.Error(err)
						return
					}
				}
			})
		})
		b.Run("RunString", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				vm := goja.New()

				for pb.Next() {
					_, err := vm.RunString(src)
					if err != nil {
						b.Error(err)
						return
					}
				}
			})
		})
		b.Run("Compile", func(b *testing.B) {
			for _, strict := range []bool{true, false} {
				name := "strict"
				if !strict {
					name = "loose"
				}
				b.Run(name, func(b *testing.B) {
					pgm, err := goja.Compile(name, src, strict)
					if err != nil {
						b.Error(err)
						return
					}

					b.Run("RunProgram", func(b *testing.B) {
						b.RunParallel(func(pb *testing.PB) {
							vm := goja.New()

							for pb.Next() {
								_, err := vm.RunProgram(pgm)
								if err != nil {
									b.Error(err)
									return
								}
							}
						})
					})
				})
			}
		})
	})
}

func BenchmarkCalling(b *testing.B) {
	b.Run("otto", func(b *testing.B) {
		vm := otto.New()

		fn, err := vm.Run(src)
		if err != nil {
			b.Error(err)
			return
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			if _, err := fn.Call(otto.UndefinedValue()); err != nil {
				b.Error(err)
				return
			}
		}
	})
	b.Run("goja", func(b *testing.B) {
		vm := goja.New()

		v, err := vm.RunString(src)
		if err != nil {
			b.Error(err)
			return
		}

		var callable goja.Callable
		if err := vm.ExportTo(v, &callable); err != nil {
			b.Error(err)
			return
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			if _, err := callable(goja.Undefined()); err != nil {
				b.Error(err)
				return
			}
		}
	})
}
