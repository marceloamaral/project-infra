package flakefinder

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("report.go", func() {

	When("sorting test data", func() {
		tests := []string{"t1", "t2", "t3"}

		It("returns all tests", func() {
			data := map[string]map[string]*Details{
				"t3": {"a": &Details{Failed: 4, Succeeded: 1, Skipped: 2, Severity: HeavilyFlaky, Jobs: []*Job{}}},
			}

			Expect(SortTestsByRelevance(data, tests)).To(BeEquivalentTo([]string{"t3", "t1", "t2"}))
		})

		It("returns no duplicated tests", func() {
			data := map[string]map[string]*Details{
				"t1": {
					"a": &Details{Failed: 3, Succeeded: 1, Skipped: 2, Severity: MostlyFlaky, Jobs: []*Job{}},
					"b": &Details{Failed: 3, Succeeded: 1, Skipped: 2, Severity: Unimportant, Jobs: []*Job{}},
				},
				"t2": {"a": &Details{Failed: 2, Succeeded: 1, Skipped: 2, Severity: MildlyFlaky, Jobs: []*Job{}}},
				"t3": {"a": &Details{Failed: 4, Succeeded: 1, Skipped: 2, Severity: HeavilyFlaky, Jobs: []*Job{}}},
			}

			Expect(SortTestsByRelevance(data, tests)).To(BeEquivalentTo([]string{"t3", "t1", "t2"}))
		})

		It("returns no duplicated tests for the end", func() {
			data := map[string]map[string]*Details{
				"t1": {
					"a": &Details{Failed: 3, Succeeded: 1, Skipped: 2, Severity: MostlyFlaky, Jobs: []*Job{}},
					"b": &Details{Failed: 3, Succeeded: 1, Skipped: 2, Severity: Unimportant, Jobs: []*Job{}},
				},
				"t2": {
					"a": &Details{Failed: 2, Succeeded: 1, Skipped: 2, Severity: MildlyFlaky, Jobs: []*Job{}},
					"b": &Details{Failed: 2, Succeeded: 1, Skipped: 2, Severity: Unimportant, Jobs: []*Job{}},
				},
				"t3": {"a": &Details{Failed: 4, Succeeded: 1, Skipped: 2, Severity: HeavilyFlaky, Jobs: []*Job{}}},
			}

			Expect(SortTestsByRelevance(data, tests)).To(BeEquivalentTo([]string{"t3", "t1", "t2"}))
		})

		It("returns tests sorted descending by severity", func() {
			data := map[string]map[string]*Details{
				"t1": {"a": &Details{Failed: 3, Succeeded: 1, Skipped: 2, Severity: MostlyFlaky, Jobs: []*Job{}}},
				"t2": {"a": &Details{Failed: 2, Succeeded: 1, Skipped: 2, Severity: MildlyFlaky, Jobs: []*Job{}}},
				"t3": {"a": &Details{Failed: 4, Succeeded: 1, Skipped: 2, Severity: HeavilyFlaky, Jobs: []*Job{}}},
			}

			Expect(SortTestsByRelevance(data, tests)).To(BeEquivalentTo([]string{"t3", "t1", "t2"}))
		})

		It("returns tests of same severity sorted descending by number of severity points", func() {
			data := map[string]map[string]*Details{
				"t1": {
					"a": &Details{Failed: 3, Succeeded: 1, Skipped: 2, Severity: HeavilyFlaky, Jobs: []*Job{}},
					"b": &Details{Failed: 3, Succeeded: 1, Skipped: 2, Severity: MostlyFlaky, Jobs: []*Job{}},
				},
				"t2": {
					"a": &Details{Failed: 2, Succeeded: 1, Skipped: 2, Severity: HeavilyFlaky, Jobs: []*Job{}},
					"b": &Details{Failed: 2, Succeeded: 1, Skipped: 2, Severity: MildlyFlaky, Jobs: []*Job{}},
				},
				"t3": {
					"a": &Details{Failed: 4, Succeeded: 1, Skipped: 2, Severity: HeavilyFlaky, Jobs: []*Job{}},
					"b": &Details{Failed: 4, Succeeded: 1, Skipped: 2, Severity: HeavilyFlaky, Jobs: []*Job{}},
				},
			}

			Expect(SortTestsByRelevance(data, tests)).To(BeEquivalentTo([]string{"t3", "t1", "t2"}))
		})

		DescribeTable("returns tests of same severity weighted by total number of tests", func(t1Failed, t2Failed, t3Failed, t1Succeeded, t2Succeeded, t3Succeeded []int, expected []string) {
			data := map[string]map[string]*Details{}
			data["t1"] = map[string]*Details{}
			for index, failed := range t1Failed {
				data["t1"][fmt.Sprint(index)] = &Details{Failed: failed, Succeeded: t1Succeeded[index], Skipped: 2, Severity: HeavilyFlaky, Jobs: []*Job{}}
			}
			data["t2"] = map[string]*Details{}
			for index, failed := range t2Failed {
				data["t2"][fmt.Sprint(index)] = &Details{Failed: failed, Succeeded: t2Succeeded[index], Skipped: 2, Severity: HeavilyFlaky, Jobs: []*Job{}}
			}
			data["t3"] = map[string]*Details{}
			for index, failed := range t3Failed {
				data["t3"][fmt.Sprint(index)] = &Details{Failed: failed, Succeeded: t3Succeeded[index], Skipped: 2, Severity: HeavilyFlaky, Jobs: []*Job{}}
			}

			Expect(SortTestsByRelevance(data, tests)).To(BeEquivalentTo(expected))
		},
			Entry("zeros shouldn't be a problem",
				[]int{2}, []int{1}, []int{3}, []int{0}, []int{0}, []int{0}, []string{"t3", "t1", "t2"},
			),
			Entry("the more failures the higher",
				[]int{2}, []int{1}, []int{3}, []int{1}, []int{1}, []int{1}, []string{"t3", "t1", "t2"},
			),
			Entry("multiple values with zeros",
				[]int{2, 0}, []int{1, 0}, []int{3, 0}, []int{1, 0}, []int{1, 0}, []int{1, 0}, []string{"t3", "t1", "t2"},
			),
			Entry("multiple values",
				[]int{4, 5}, []int{3, 4}, []int{6, 7}, []int{1, 1}, []int{1, 1}, []int{1, 1}, []string{"t3", "t1", "t2"},
			),
			Entry("errors high, ratios small",
				[]int{6, 7}, []int{4, 5}, []int{11, 12}, []int{5, 6}, []int{3, 4}, []int{10, 11}, []string{"t3", "t1", "t2"},
			),
			Entry("higher ratio, the higher",
				[]int{6, 7}, []int{4, 5}, []int{11, 12}, []int{3, 4}, []int{2, 3}, []int{2, 1}, []string{"t3", "t1", "t2"},
			),
			Entry("mixed lengths",
				[]int{8, 10}, []int{3, 4, 5}, []int{22}, []int{2, 1}, []int{1, 2, 3}, []int{2}, []string{"t3", "t1", "t2"},
			),
		)

	})

	When("sorting test via severity", func() {

		It("returns tests of same severity sorted descending by number of severity points", func() {
			Expect(BuildUpSortedTestsBySeverity(map[string]map[string]int{
				"t1": {HeavilyFlaky: 2},
				"t2": {HeavilyFlaky: 1},
				"t3": {HeavilyFlaky: 3},
			})).To(BeEquivalentTo([]string{"t3", "t1", "t2"}))
		})

		It("returns tests of same severity and same number of severity points sorted lexically", func() {
			Expect(BuildUpSortedTestsBySeverity(map[string]map[string]int{
				"tb": {HeavilyFlaky: 2},
				"tc": {HeavilyFlaky: 2},
				"ta": {HeavilyFlaky: 2},
			})).To(BeEquivalentTo([]string{"ta", "tb", "tc"}))
		})

		It("returns tests of same severity sorted by lower severity", func() {
			Expect(BuildUpSortedTestsBySeverity(map[string]map[string]int{
				"tb": {HeavilyFlaky: 2, MostlyFlaky: 2},
				"tc": {HeavilyFlaky: 2, MostlyFlaky: 1},
				"ta": {HeavilyFlaky: 2, MostlyFlaky: 3},
			})).To(BeEquivalentTo([]string{"ta", "tb", "tc"}))
		})

		It("returns tests of same severity sorted by lower severity if even lower values present but zero", func() {
			Expect(BuildUpSortedTestsBySeverity(map[string]map[string]int{
				"tb": {HeavilyFlaky: 2, MostlyFlaky: 2, ModeratelyFlaky: 0, MildlyFlaky: 0},
				"tc": {HeavilyFlaky: 2, MostlyFlaky: 1, ModeratelyFlaky: 0, MildlyFlaky: 0},
				"ta": {HeavilyFlaky: 2, MostlyFlaky: 3, ModeratelyFlaky: 0, MildlyFlaky: 0},
			})).To(BeEquivalentTo([]string{"ta", "tb", "tc"}))
		})

		It("returns tests of same severity sorted by lower severity if inbetween values present but zero", func() {
			Expect(BuildUpSortedTestsBySeverity(map[string]map[string]int{
				"tb": {HeavilyFlaky: 2, MostlyFlaky: 0, ModeratelyFlaky: 2, MildlyFlaky: 0},
				"tc": {HeavilyFlaky: 2, MostlyFlaky: 0, ModeratelyFlaky: 1, MildlyFlaky: 0},
				"ta": {HeavilyFlaky: 2, MostlyFlaky: 0, ModeratelyFlaky: 3, MildlyFlaky: 0},
			})).To(BeEquivalentTo([]string{"ta", "tb", "tc"}))
		})

		It("returns tests of same severity sorted by lower severity if some inbetween values zero", func() {
			Expect(BuildUpSortedTestsBySeverity(map[string]map[string]int{
				"tb": {HeavilyFlaky: 2, MostlyFlaky: 0, ModeratelyFlaky: 2, MildlyFlaky: 0},
				"tc": {HeavilyFlaky: 2, MostlyFlaky: 0, ModeratelyFlaky: 0, MildlyFlaky: 0},
				"ta": {HeavilyFlaky: 2, MostlyFlaky: 3, ModeratelyFlaky: 0, MildlyFlaky: 0},
			})).To(BeEquivalentTo([]string{"ta", "tb", "tc"}))
		})

		It("returns tests of same severity sorted by lower severity if some inbetween values zero with more values", func() {
			Expect(BuildUpSortedTestsBySeverity(map[string]map[string]int{
				"tb": {HeavilyFlaky: 1, MostlyFlaky: 0, ModeratelyFlaky: 0, MildlyFlaky: 1, Fine: 0, Unimportant: 0},
				"tc": {HeavilyFlaky: 1, MostlyFlaky: 0, ModeratelyFlaky: 0, MildlyFlaky: 0, Fine: 0, Unimportant: 0},
				"ta": {HeavilyFlaky: 1, MostlyFlaky: 0, ModeratelyFlaky: 0, MildlyFlaky: 2, Fine: 0, Unimportant: 0},
			})).To(BeEquivalentTo([]string{"ta", "tb", "tc"}))
		})

	})

	DescribeTable("When comparing severity",
		func(a, b *TestToSeverityOccurrences, expected bool) {
			bySeverity := []*TestToSeverityOccurrences{a, b}
			Expect(BySeverity.Less(bySeverity, 0, 1)).To(BeEquivalentTo(expected))
		},

		Entry("ta -> Sev(2) less than tb -> Sev(2) is false",
			&TestToSeverityOccurrences{Name: "ta", SeverityOccurrences: []int{2}},
			&TestToSeverityOccurrences{Name: "tb", SeverityOccurrences: []int{2}},
			false,
		),
		Entry("tb -> Sev(2) less than ta -> Sev(2) is true",
			&TestToSeverityOccurrences{Name: "tb", SeverityOccurrences: []int{2}},
			&TestToSeverityOccurrences{Name: "ta", SeverityOccurrences: []int{2}},
			true,
		),
		Entry("ta -> Sev(2) less than ta -> Sev(2) is false",
			&TestToSeverityOccurrences{Name: "ta", SeverityOccurrences: []int{2}},
			&TestToSeverityOccurrences{Name: "ta", SeverityOccurrences: []int{2}},
			false,
		),
		Entry("tb -> Sev(3) is less than ta -> Sev(2) is false",
			&TestToSeverityOccurrences{Name: "tb", SeverityOccurrences: []int{3}},
			&TestToSeverityOccurrences{Name: "ta", SeverityOccurrences: []int{2}},
			false,
		),
		Entry("ta -> Sev(2) is less than tb -> Sev(3) is true",
			&TestToSeverityOccurrences{Name: "ta", SeverityOccurrences: []int{2}},
			&TestToSeverityOccurrences{Name: "tb", SeverityOccurrences: []int{3}},
			true,
		),
		Entry("tb -> Sev(3, 2) is less ta -> Sev(3, 3) is true",
			&TestToSeverityOccurrences{Name: "tb", SeverityOccurrences: []int{3, 2}},
			&TestToSeverityOccurrences{Name: "ta", SeverityOccurrences: []int{3, 3}},
			true,
		),
		Entry("tb -> Sev(3, 3) is less ta -> Sev(3, 2) is false",
			&TestToSeverityOccurrences{Name: "tb", SeverityOccurrences: []int{3, 3}},
			&TestToSeverityOccurrences{Name: "ta", SeverityOccurrences: []int{3, 2}},
			false,
		),
		Entry("tb -> Sev(3, 0, 3) is less ta -> Sev(3, 0, 2) is false",
			&TestToSeverityOccurrences{Name: "tb", SeverityOccurrences: []int{3, 0, 3}},
			&TestToSeverityOccurrences{Name: "ta", SeverityOccurrences: []int{3, 0, 2}},
			false,
		),
		Entry("tb -> Sev(3, 0, 2) is not less ta -> Sev(3, 0, 3) is true",
			&TestToSeverityOccurrences{Name: "tb", SeverityOccurrences: []int{3, 0, 2}},
			&TestToSeverityOccurrences{Name: "ta", SeverityOccurrences: []int{3, 0, 3}},
			true,
		),
		Entry("tb -> Sev(1,0,0,2,0,0) is less ta -> Sev(1,0,0,1,0,0) is false",
			&TestToSeverityOccurrences{Name: "tb", SeverityOccurrences: []int{1, 0, 0, 2, 0, 0}},
			&TestToSeverityOccurrences{Name: "ta", SeverityOccurrences: []int{1, 0, 0, 1, 0, 0}},
			false,
		),
		Entry("ta -> Sev(1,0,0,1,0,0) is less tb -> Sev(1,0,0,2,0,0) is true",
			&TestToSeverityOccurrences{Name: "ta", SeverityOccurrences: []int{1, 0, 0, 1, 0, 0}},
			&TestToSeverityOccurrences{Name: "tb", SeverityOccurrences: []int{1, 0, 0, 2, 0, 0}},
			true,
		),
		Entry("tb -> Sev(1,0,1,0,0,0) is less ta -> Sev(1,0,0,1,0,0) is false",
			&TestToSeverityOccurrences{Name: "tb", SeverityOccurrences: []int{1, 0, 1, 0, 0, 0}},
			&TestToSeverityOccurrences{Name: "ta", SeverityOccurrences: []int{1, 0, 0, 1, 0, 0}},
			false,
		),
		Entry("ta -> Sev(1,0,0,1,0,0) is less tb -> Sev(1,0,1,0,0,0) is true",
			&TestToSeverityOccurrences{Name: "ta", SeverityOccurrences: []int{1, 0, 0, 1, 0, 0}},
			&TestToSeverityOccurrences{Name: "tb", SeverityOccurrences: []int{1, 0, 1, 0, 0, 0}},
			true,
		),
	)

	When("swapping elements", func() {

		It("Works", func() {
			bySeverity := []*TestToSeverityOccurrences{
				{Name: "tb", SeverityOccurrences: []int{3, 0, 2}},
				{Name: "ta", SeverityOccurrences: []int{3, 0, 3}},
			}
			BySeverity.Swap(bySeverity, 0, 1)
			Expect(bySeverity[0].Name).To(BeEquivalentTo("ta"))
			Expect(bySeverity[1].Name).To(BeEquivalentTo("tb"))
		})

	})

	DescribeTable("When calculating severity",
		func(details *Details, expected string) {
			SetSeverity(details)
			Expect(details.Severity).To(BeEquivalentTo(expected))
		},
		Entry("results having no failed tests but successful tests is fine", &Details{Failed: 0, Succeeded: 1, Skipped: 2, Jobs: []*Job{}}, Fine),
		Entry("results having no successful tests is heavily flaky", &Details{Failed: 1, Succeeded: 0, Skipped: 2, Jobs: []*Job{}}, HeavilyFlaky),
		Entry("results being HeavilyFlaky", &Details{Failed: 1, Succeeded: 1, Skipped: 2, Jobs: []*Job{}}, HeavilyFlaky),
		Entry("results being MostlyFlaky", &Details{Failed: 1, Succeeded: 2, Skipped: 2, Jobs: []*Job{}}, MostlyFlaky),
		Entry("results being ModeratelyFlaky", &Details{Failed: 1, Succeeded: 5, Skipped: 2, Jobs: []*Job{}}, ModeratelyFlaky),
		Entry("results being MildlyFlaky", &Details{Failed: 1, Succeeded: 10, Skipped: 2, Jobs: []*Job{}}, MildlyFlaky),
	)
})
