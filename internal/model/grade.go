package model

// IsValidGrade 检查成绩是否有效
func IsValidGrade(grade string) bool {
	validGrades := map[string]bool{
		"A":  true,
		"A-": true,
		"B+": true,
		"B":  true,
		"B-": true,
		"C+": true,
		"C":  true,
		"C-": true,
		"D+": true,
		"D":  true,
		"F":  true,
	}
	return validGrades[grade]
}

// IsPassingGrade 检查是否及格
func IsPassingGrade(grade string) bool {
	failingGrades := map[string]bool{
		"F": true,
	}
	return !failingGrades[grade]
}
