package main

import (
	"reflect"
	"testing"
)

func TestMinuteField(t *testing.T) {
	minuteParser := &MinuteField{BaseCronField{Min: 0, Max: 59}}

	tests := []struct {
		input    string
		expected []int
		wantErr  bool
	}{
		{"*", generateRange(0, 59), false},
		{"*/15", []int{0, 15, 30, 45}, false},
		{"1,5,10", []int{1, 5, 10}, false},
		{"0-10", generateRange(0, 10), false},
		{"70", nil, true},
	}

	runTests(t, minuteParser, tests)
}

func TestHourField(t *testing.T) {
	hourParser := &HourField{BaseCronField{Min: 0, Max: 23}}

	tests := []struct {
		input    string
		expected []int
		wantErr  bool
	}{
		{"*", generateRange(0, 23), false},
		{"*/6", []int{0, 6, 12, 18}, false},
		{"0-12", generateRange(0, 12), false},
		{"24", nil, true},
	}

	runTests(t, hourParser, tests)
}

func TestDayOfMonthField(t *testing.T) {
	dayParser := &DayOfMonthField{BaseCronField{Min: 1, Max: 31}}

	tests := []struct {
		input    string
		expected []int
		wantErr  bool
	}{
		{"*", generateRange(1, 31), false},
		{"5,10,20", []int{5, 10, 20}, false},
		{"1-7", generateRange(1, 7), false},
		{"32", nil, true},
		{"31", generateRange(31, 31), false},
	}

	runTests(t, dayParser, tests)
}

func TestMonthField(t *testing.T) {
	monthParser := &MonthField{BaseCronField{Min: 1, Max: 12}}

	tests := []struct {
		input    string
		expected []int
		wantErr  bool
	}{
		{"*", generateRange(1, 12), false},
		{"1,6,12", []int{1, 6, 12}, false},
		{"1-6", generateRange(1, 6), false},
		{"13", nil, true},
	}

	runTests(t, monthParser, tests)
}

func TestDayOfWeekField(t *testing.T) {
	weekParser := &DayOfWeekField{BaseCronField{Min: 0, Max: 6}}

	tests := []struct {
		input    string
		expected []int
		wantErr  bool
	}{
		{"*", generateRange(0, 6), false},
		{"0,3,6", []int{0, 3, 6}, false},
		{"1-5", generateRange(1, 5), false},
		{"7", nil, true},
	}

	runTests(t, weekParser, tests)
}

func TestCronExpressionParser(t *testing.T) {
	parser := NewCronExpressionParser()

	tests := []struct {
		expression string
		wantErr    bool
	}{
		{"*/15 0 1,15 * 1-5 /usr/bin/find", false},
		{"60 0 1 1 1 echo 'invalid'", true},
		{"*/5 * * * * /some/command", false},
		{"*/5 * * *", true},
	}

	for _, tt := range tests {
		_, _, err := parser.Parse(tt.expression)
		if (err != nil) != tt.wantErr {
			t.Errorf("Parse(%q) unexpected error: %v", tt.expression, err)
		}
	}
}

func generateRange(start, end int) []int {
	var result []int
	for i := start; i <= end; i++ {
		result = append(result, i)
	}
	return result
}

func runTests(t *testing.T, parser CronFieldParser, tests []struct {
	input    string
	expected []int
	wantErr  bool
}) {
	for _, tt := range tests {
		result, err := parser.Parse(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
		}
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("Parse(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}
