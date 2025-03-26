package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CronFieldParser interface {
	Parse(field string) ([]int, error)
}

type BaseCronField struct {
	Min, Max int
}

func (b *BaseCronField) parseField(value string) ([]int, error) {
	var results []int

	switch {
	case value == "*":
		for i := b.Min; i <= b.Max; i++ {
			results = append(results, i)
		}
		return results, nil

	case strings.HasPrefix(value, "*/"):
		step, err := strconv.Atoi(strings.TrimPrefix(value, "*/"))
		if err != nil {
			return nil, fmt.Errorf("invalid step value: %s", value)
		}
		for i := b.Min; i <= b.Max; i += step {
			results = append(results, i)
		}
		return results, nil

	case strings.Contains(value, "-"):
		rangeParts := strings.Split(value, "-")
		if len(rangeParts) != 2 {
			return nil, fmt.Errorf("invalid range format: %s", value)
		}
		start, err1 := strconv.Atoi(rangeParts[0])
		end, err2 := strconv.Atoi(rangeParts[1])
		if err1 != nil || err2 != nil || start > end || start < b.Min || end > b.Max {
			return nil, fmt.Errorf("invalid range: %s", value)
		}
		for i := start; i <= end; i++ {
			results = append(results, i)
		}
		return results, nil

	case strings.Contains(value, ","):
		entries := strings.Split(value, ",")
		for _, entry := range entries {
			num, err := strconv.Atoi(entry)
			if err != nil || num < b.Min || num > b.Max {
				return nil, fmt.Errorf("invalid value: %s", entry)
			}
			results = append(results, num)
		}
		return results, nil

	default:
		num, err := strconv.Atoi(value)
		if err != nil || num < b.Min || num > b.Max {
			return nil, fmt.Errorf("invalid value: %s", value)
		}
		return []int{num}, nil
	}
}

type MinuteField struct{ BaseCronField }

func (m *MinuteField) Parse(field string) ([]int, error) {
	m.Min, m.Max = 0, 59
	return m.parseField(field)
}

type HourField struct{ BaseCronField }

func (h *HourField) Parse(field string) ([]int, error) {
	h.Min, h.Max = 0, 23
	return h.parseField(field)
}

type DayOfMonthField struct{ BaseCronField }

func (d *DayOfMonthField) Parse(field string) ([]int, error) {
	d.Min, d.Max = 1, 31
	return d.parseField(field)
}

type MonthField struct{ BaseCronField }

func (m *MonthField) Parse(field string) ([]int, error) {
	m.Min, m.Max = 1, 12
	return m.parseField(field)
}

type DayOfWeekField struct{ BaseCronField }

func (d *DayOfWeekField) Parse(field string) ([]int, error) {
	d.Min, d.Max = 0, 6
	return d.parseField(field)
}

type CronExpressionParser struct {
	fields map[string]CronFieldParser
}

func NewCronExpressionParser() *CronExpressionParser {
	return &CronExpressionParser{
		fields: map[string]CronFieldParser{
			"minute":       &MinuteField{},
			"hour":         &HourField{},
			"day of month": &DayOfMonthField{},
			"month":        &MonthField{},
			"day of week":  &DayOfWeekField{},
		},
	}
}

func (cp *CronExpressionParser) Parse(expression string) (map[string][]int, string, error) {
	parts := strings.Fields(expression)
	if len(parts) < 6 {
		return nil, "", fmt.Errorf("invalid cron expression: expected at least 6 fields")
	}

	cronFields := []string{"minute", "hour", "day of month", "month", "day of week"}
	parsedValues := make(map[string][]int)

	for i, field := range cronFields {
		parser, exists := cp.fields[field]
		if !exists {
			return nil, "", fmt.Errorf("no parser for field: %s", field)
		}
		values, err := parser.Parse(parts[i])
		if err != nil {
			return nil, "", err
		}
		parsedValues[field] = values
	}

	command := strings.Join(parts[5:], " ")
	return parsedValues, command, nil
}

func PrintFormattedOutput(parsed map[string][]int, command string) {
	for field, values := range parsed {
		fmt.Printf("%-14s %s\n", field, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(values)), " "), "[]"))
	}
	fmt.Printf("%-14s %s\n", "command", command)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./deliveroo-cron-parser \"*/15 0 1,15 * 1-5 /usr/bin/find\"")
		os.Exit(1)
	}

	expression := os.Args[1]
	parser := NewCronExpressionParser()
	parsed, command, err := parser.Parse(expression)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	PrintFormattedOutput(parsed, command)
}
