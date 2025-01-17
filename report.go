package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

func (p *processor) reportFlaky() {
	var flaky []flakyTest

	for test, count := range p.failed {
		if p.passed[test] != 0 {
			flaky = append(flaky, flakyTest{
				test:   test,
				passed: p.passed[test],
				failed: count,
			})
		}
	}

	sort.Slice(flaky, func(i, j int) bool {
		return flaky[i].test > flaky[j].test
	})

	if len(flaky) > 0 {
		p.counts["flaky"] = len(flaky)

		if p.fl.Markdown {
			fmt.Println("## Flaky tests")
			fmt.Println("<details>")
			fmt.Printf("<summary>Tests: %d</summary>\n\n", len(flaky))

			fmt.Println("| Pass | Fail | Test |")
			fmt.Println("| - | - | - |")

			for _, ft := range flaky {
				fmt.Printf("| %d | %d | %s |\n", ft.passed, ft.failed, ft.test)
			}

			fmt.Println("</details>")
		} else {
			fmt.Println("Flaky tests:")

			for _, ft := range flaky {
				fmt.Printf("%s: %d passed, %d failed\n", ft.test, ft.passed, ft.failed)
			}
		}

		fmt.Println()
	}
}

func (p *processor) reportSlowest() {
	sort.Slice(p.slowest, func(i, j int) bool {
		return p.slowest[i].Elapsed > p.slowest[j].Elapsed
	})

	if len(p.slowest) > 0 {
		if p.fl.Markdown {
			fmt.Println("## Slow tests")
			fmt.Println("<details>")
			fmt.Printf("<summary>Total slow runs: %d</summary>\n\n", len(p.slowest))

			fmt.Println("| Result | Duration | Package | Test |")
			fmt.Println("| - | - | - | - |")

			for i, l := range p.slowest {
				if i >= p.fl.Slowest {
					break
				}

				dur := time.Duration(l.Elapsed * float64(time.Second))
				fmt.Printf("| %s | %s | %s | %s |\n", l.Action, dur.String(), l.Package, l.Test)
			}

			fmt.Println("</details>")
		} else {
			fmt.Println("Slowest tests:")

			for i, l := range p.slowest {
				if i >= p.fl.Slowest {
					break
				}

				dur := time.Duration(l.Elapsed * float64(time.Second))
				fmt.Printf("%s %s %s %s\n", l.Action, l.Package, l.Test, dur.String())
			}
		}

		fmt.Println()
	}
}

func (p *processor) reportRaces() {
	if len(p.strippedDataRaces) > 0 {
		var keys []string

		for k := range p.strippedDataRaces {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		if p.fl.Markdown {
			fmt.Println("## Data races")
			fmt.Println("<details>")
			fmt.Printf("<summary>Total data races: %d, unique: %d</summary>\n\n",
				len(p.dataRaces), len(p.strippedDataRaces))

			for _, k := range keys {
				r := p.strippedDataRaces[k]
				t := p.strippedTests[k]

				fmt.Println("<details>")
				fmt.Printf("<summary><code>%s</code></summary>\n\n", t[0])

				if len(t) > 1 {
					fmt.Println("Other affected tests:")
					fmt.Println("```")

					for _, tt := range t[1:] {
						fmt.Println(tt)
					}

					fmt.Println("```")
				}

				fmt.Println("\n```")
				fmt.Println(r)
				fmt.Println("```")
				fmt.Println("</details>")
				fmt.Println()
			}

			fmt.Println("</details>")
			fmt.Println()
		} else {
			fmt.Println("Data races:")

			for _, k := range keys {
				t := p.strippedTests[k]

				if len(t) > 3 {
					t = append(t[0:3], "...")
				}

				fmt.Println(strings.Join(t, ", "))
				fmt.Println(p.strippedDataRaces[k])
			}

			fmt.Println()
		}
	}
}

func (p *processor) reportPackages() {
	if len(p.packageElapsed) > 0 {
		sort.Slice(p.packageElapsed, func(i, j int) bool {
			return p.packageElapsed[i].Elapsed > p.packageElapsed[j].Elapsed
		})

		if p.fl.Markdown {
			fmt.Println("## Slowest test packages")
			fmt.Println("<details>")
			fmt.Printf("<summary>Total packages with tests: %d</summary>\n\n", len(p.packageElapsed))

			fmt.Println("| Duration | Package |")
			fmt.Println("| - | - |")

			for i, ps := range p.packageElapsed {
				dur := time.Duration(ps.Elapsed * float64(time.Second))

				fmt.Printf("| %s | %s |\n", dur, ps.Package)

				if i > p.fl.Slowest {
					break
				}
			}

			fmt.Println("</details>")
			fmt.Println()
		}
	}
}

func (p *processor) reportFailed() {
	if len(p.failures) > 0 {
		if p.fl.Markdown {
			fmt.Println("## Failed tests")
			fmt.Println("<details>")
			fmt.Printf("<summary>Failed: %d</summary>\n\n", len(p.failures))

			for test, output := range p.failures {
				fmt.Println("<details>")
				fmt.Printf("<summary><code>%s</code></summary>\n\n", test)

				fmt.Println("```")
				fmt.Println(strings.Join(output, ""))
				fmt.Println("```")

				fmt.Println("</details>")
			}

			fmt.Println("</details>")
			fmt.Println()
		}
	}
}

func (p *processor) report() {
	p.reportFlaky()
	p.reportSlowest()
	p.reportRaces()
	p.reportPackages()
	p.reportFailed()

	if p.fl.Markdown {
		fmt.Println("## Metrics")
		fmt.Println()

		fmt.Printf("```\n%v\n```\n\n", p.counts)
		fmt.Println("Elapsed:", p.elapsed.String())
		fmt.Println("Slow:", p.elapsedSlow.String())

		fmt.Println()

		fmt.Println("## Elapsed distribution (seconds)")
		fmt.Println("```")
		fmt.Println(p.hist.String())
		fmt.Println("```")
	} else {
		fmt.Printf("Metrics: %v\n", p.counts)
		fmt.Println("Elapsed:", p.elapsed.String())
		fmt.Println("Slow:", p.elapsedSlow.String())

		fmt.Println()

		fmt.Println("Elapsed distribution:")
		fmt.Println(p.hist.String())
	}
}
