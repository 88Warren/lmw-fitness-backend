// database/seed.go

package database

import "gorm.io/gorm"

func SeedDB(db *gorm.DB) {
	ExerciseSeed()
	ProgramSeed()
	BeginnerWorkoutDaySeed()
	AdvancedWorkoutDaySeed()
	BlogSeed(db)
	UserSeed(db)
}
