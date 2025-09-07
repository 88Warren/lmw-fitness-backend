package database

import (
	"log"

	"github.com/88warren/lmw-fitness-backend/models"
	"gorm.io/gorm"
)

func getExerciseIDByName(db *gorm.DB, name string) (uint, error) {
	var exercise models.Exercise
	if err := db.Where("name = ?", name).First(&exercise).Error; err != nil {
		return 0, err
	}
	return exercise.ID, nil
}

func getProgramIDByName(db *gorm.DB, name string) (uint, error) {
	var program models.WorkoutProgram
	if err := db.Where("name = ?", name).First(&program).Error; err != nil {
		return 0, err
	}
	return program.ID, nil
}

func BeginnerWorkoutDaySeed() {
	log.Println("Seeding workout day, block, and exercise data for the 30-Day Beginner Program...")

	programID, err := getProgramIDByName(DB, "beginner-program")
	if err != nil {
		log.Fatalf("Failed to find '30-Day Beginner Program': %v", err)
	}

	exIDs := make(map[string]uint)
	exercises := []string{
		"Ab Twists", "Bearcrawls", "Belt Kicks", "Bicycles", "Box Jumps", "Broad Jumps",
		"Burpees", "Burpees (modified)", "Calf Jumps", "Calf Raises", "Cross Jabs",
		"Crunches", "Diamond Press Ups (on Knees)", "Diamond Sit Ups", "Donkey Kicks",
		"Dorsal Raises", "Elbows to Knee", "Flutter Kicks", "Glute Bridges", "H.O.G. Press Ups", "H.O.G. Press Ups (on Knees)",
		"Half Sit Ups", "Heel Taps", "High Knees", "Jack Knife", "Knees to Chest", "Lateral Lunges",
		"Leg Raises", "Lunges", "Mobility", "Mountain Climbers", "Moving Press Ups", "Oblique Hops",
		"Oblique Plank", "Overhead Jabs", "Plank Hold", "Plank Jabs", "Plank Leg Raises",
		"Plank Shoulder Taps", "Press Ups", "Press Ups (on Knees)", "Reverse Lunge", "Scissors",
		"Sit Ups", "Ski Jumps", "Sprawls", "Sprints", "Squat Jumps", "Squat Kicks", "Squat Twists",
		"Squats", "Squat Hold", "Starjumps", "Standing Mountain Climbers", "Straddle Sit Ups", "Superman",
		"Switch Kicks", "T-Runs", "Toe Taps", "Tricep Dips (Floor)", "Tricep Dips (with Chair)",
		"V Press Ups", "Walkaways", "Wide Arm Press Ups (on Knees)", "Y Shaped Lunges",
	}

	for _, name := range exercises {
		id, err := getExerciseIDByName(DB, name)
		if err != nil {
			log.Fatalf("Failed to find exercise '%s': %v", name, err)
		}
		exIDs[name] = id
	}

	createWorkoutDay := func(day models.WorkoutDay) {
		var existingDay models.WorkoutDay
		if err := DB.Where("program_id = ? AND day_number = ?", day.ProgramID, day.DayNumber).First(&existingDay).Error; err == nil {
			log.Printf("Beginner Program - Day %d already exists, skipping creation.", day.DayNumber)
			return
		}
		if err := DB.Create(&day).Error; err != nil {
			log.Printf("Failed to create Beginner Program - Day %d: %v", day.DayNumber, err)
		} else {
			log.Printf("Successfully created Beginner Program - Day %d: %s", day.DayNumber, day.Title)
		}
	}

	// --- DAY 1 ---
	day1 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   1,
		Title:       "Fitness Assessment",
		Description: "Complete these 8 exercises for 1 minute each. Make sure you record your results, you will need them for day 30 - there's a table attached in your email to help",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "Fitness Assessment",
				BlockNotes: "Try to do as many reps as possible. Use the whole 2 mins rest after each exercise, to be able to give 100% effort for the next exercise.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Ups (on Knees)"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 2, ExerciseID: exIDs["Squats"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 3, ExerciseID: exIDs["Plank Hold"], Reps: "Max Effort", Duration: "Max Time", Rest: "2 mins"},
					{Order: 4, ExerciseID: exIDs["Burpees (modified)"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 5, ExerciseID: exIDs["Starjumps"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 6, ExerciseID: exIDs["Sit Ups"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 7, ExerciseID: exIDs["Lunges"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins", Tips: "2 Lunges = 1 rep"},
					{Order: 8, ExerciseID: exIDs["Tricep Dips (with Chair)"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
				},
			},
		},
	}
	createWorkoutDay(day1)

	// --- DAY 2 ---
	day2 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   2,
		Title:       "Lower Body Focus",
		Description: "A circuit focusing on your lower body to build strength and endurance.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 3,
				BlockNotes:  "Exercise for 40 seconds and then rest for 20 seconds. Repeat the circuit 3 times. Full duration 18 minutes.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squats"], Duration: "40s", Rest: "20s"},
					{Order: 2, ExerciseID: exIDs["Lunges"], Duration: "40s", Rest: "20s"},
					{Order: 3, ExerciseID: exIDs["Calf Raises"], Duration: "40s", Rest: "20s"},
					{Order: 4, ExerciseID: exIDs["Glute Bridges"], Duration: "40s", Rest: "20s"},
					{Order: 5, ExerciseID: exIDs["Squat Kicks"], Duration: "40s", Rest: "20s"},
					{Order: 6, ExerciseID: exIDs["Donkey Kicks"], Duration: "40s", Rest: "20s"},
				},
			},
		},
	}
	createWorkoutDay(day2)

	// --- DAY 3 ---
	day3 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   3,
		Title:       "Upper Body Focus",
		Description: "A circuit focusing on your upper body to build strength and endurance.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 3,
				RoundRest:   "60s",
				BlockNotes:  "Exercise for 30 seconds and then rest for 15 seconds. Repeat the circuit 3 times. Rest 60 seconds between rounds. Full duration 20 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Wide Arm Press Ups (on Knees)"], Duration: "30s", Rest: "15s"},
					{Order: 2, ExerciseID: exIDs["Tricep Dips (with Chair)"], Duration: "30s", Rest: "15s"},
					{Order: 3, ExerciseID: exIDs["Plank Shoulder Taps"], Duration: "30s", Rest: "15s"},
					{Order: 4, ExerciseID: exIDs["Superman"], Duration: "30s", Rest: "15s"},
					{Order: 6, ExerciseID: exIDs["Plank Hold"], Duration: "30s", Rest: "15s"},
					{Order: 6, ExerciseID: exIDs["Walkaways"], Duration: "30s", Rest: "15s"},
					{Order: 7, ExerciseID: exIDs["Cross Jabs"], Duration: "30s", Rest: "15s"},
					{Order: 8, ExerciseID: exIDs["Dorsal Raises"], Duration: "30s", Rest: "15s"},
				},
			},
		},
	}
	createWorkoutDay(day3)

	// --- DAY 4 ---
	day4 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   4,
		Title:       "Cardio & Core",
		Description: "As many rounds as possible (AMRAP) in 12 minutes. Track your rounds with the counter!",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "AMRAP",
				BlockNotes: "12 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpees (modified)"], Reps: "5"},
					{Order: 2, ExerciseID: exIDs["Mountain Climbers"], Reps: "10", Tips: "2 Climbers = 1 rep"},
					{Order: 3, ExerciseID: exIDs["Squat Jumps"], Reps: "15"},
					{Order: 4, ExerciseID: exIDs["Ab Twists"], Reps: "20"},
				},
			},
		},
	}
	createWorkoutDay(day4)

	// --- DAY 5 ---
	day5 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   5,
		Title:       "Full Body Circuit",
		Description: "Every Minute on the Minute (EMOM), complete the number of reps within the minute. Faster you complete the more rest",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "EMOM",
				BlockRounds: 3,
				BlockNotes:  "15 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, Instructions: "Minute 1", ExerciseID: exIDs["Squats"], Reps: "10"},
					{Order: 2, Instructions: "Minute 2", ExerciseID: exIDs["Press Ups (on Knees)"], Reps: "8"},
					{Order: 3, Instructions: "Minute 3", ExerciseID: exIDs["Lunges"], Reps: "6", Tips: "2 Lunges = 1 rep"},
					{Order: 4, Instructions: "Minute 4", ExerciseID: exIDs["Crunches"], Reps: "15"},
					{Order: 5, Instructions: "Minute 5", ExerciseID: exIDs["Starjumps"], Reps: "10"},
				},
			},
		},
	}
	createWorkoutDay(day5)

	// --- DAY 6 ---
	day6 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   6,
		Title:       "Core Blast",
		Description: "A timed core workout for stability.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 4,
				RoundRest:   "60s",
				BlockNotes:  "Exercise for 30 seconds, no rest between exercises. 60 Second rest between rounds. Full duration 16 mintutes.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Leg Raises"], Duration: "30s", Rest: "0s"},
					{Order: 2, ExerciseID: exIDs["Bicycles"], Duration: "30s", Rest: "0s"},
					{Order: 3, ExerciseID: exIDs["Flutter Kicks"], Duration: "30s", Rest: "0s"},
					{Order: 4, ExerciseID: exIDs["Half Sit Ups"], Duration: "30s", Rest: "0s"},
					{Order: 5, ExerciseID: exIDs["Heel Taps"], Duration: "30s", Rest: "0s"},
					{Order: 6, ExerciseID: exIDs["Glute Bridges"], Duration: "30s", Rest: "0s"},
				},
			},
		},
	}
	createWorkoutDay(day6)

	// --- DAY 7 ---
	day7 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   7,
		Title:       "Full Body Flow",
		Description: "Every Minute on the Minute (EMOM), complete the number of reps within the minute. Faster you complete the more rest",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "EMOM",
				BlockRounds: 4,
				BlockNotes:  "16 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, Instructions: "Minute 1", ExerciseID: exIDs["Squats"], Reps: "15"},
					{Order: 2, Instructions: "Minute 2", ExerciseID: exIDs["Plank Hold"], Reps: "30s"},
					{Order: 3, Instructions: "Minute 3", ExerciseID: exIDs["Press Ups (on Knees)"], Reps: "15"},
					{Order: 4, Instructions: "Minute 4", ExerciseID: exIDs["Starjumps"], Reps: "30"},
				},
			},
		},
	}
	createWorkoutDay(day7)

	// --- DAY 8 ---
	day8 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   8,
		Title:       "Recovery day & Optional workout",
		Description: "Mobility AND/OR Workout: Upper body strength.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Mobility",
				BlockRounds: 1,
				BlockNotes:  "A mobility session to stretch your tight muscle. Prevent injury and aid recovery",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Mobility"]},
				},
			},
			{
				BlockType:   "Circuit",
				BlockRounds: 3,
				BlockNotes:  "Exercise for 40 seconds and then rest for 20 seconds. Repeat the circuit 3 times. Full duration 18 minutes.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Ups (on Knees)"], Duration: "40s", Rest: "20s"},
					{Order: 2, ExerciseID: exIDs["Tricep Dips (with Chair)"], Duration: "40s", Rest: "20s"},
					{Order: 3, ExerciseID: exIDs["Walkaways"], Duration: "40s", Rest: "20s"},
					{Order: 4, ExerciseID: exIDs["Jack Knife"], Duration: "40s", Rest: "20s"},
					{Order: 5, ExerciseID: exIDs["Plank Leg Raises"], Duration: "40s", Rest: "20s"},
					{Order: 6, ExerciseID: exIDs["Diamond Sit Ups"], Duration: "40s", Rest: "20s"},
				},
			},
		},
	}
	createWorkoutDay(day8)

	// --- DAY 9 ---
	day9 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   9,
		Title:       "Lower Body Power",
		Description: "A circuit to build explosive lower body strength.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 4,
				BlockNotes:  "Exercise for 40 seconds and then rest for 20 seconds. Repeat the circuit 4 times. Full duration 20 minutes.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Kicks"], Duration: "40s", Rest: "20s"},
					{Order: 2, ExerciseID: exIDs["Y Shaped Lunges"], Duration: "40s", Rest: "20s"},
					{Order: 3, ExerciseID: exIDs["Calf Jumps"], Duration: "40s", Rest: "20s"},
					{Order: 4, ExerciseID: exIDs["Squat Jumps"], Duration: "40s", Rest: "20s"},
					{Order: 5, ExerciseID: exIDs["Lateral Lunges"], Duration: "40s", Rest: "20s"},
				},
			},
		},
	}
	createWorkoutDay(day9)

	// --- DAY 10 ---
	day10 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   10,
		Title:       "Core & Stability",
		Description: "A Tabata workout focused on core stability.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Tabata",
				BlockRounds: 8,
				BlockNotes:  "20s work / 10s rest x 8 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Mountain Climbers"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Plank Hold"], Duration: "20s", Rest: "10s"},
				},
			},
			{
				BlockType:   "Tabata",
				BlockRounds: 8,
				BlockNotes:  "20s work / 10s rest x 8 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Bicycles"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Leg Raises"], Duration: "20s", Rest: "10s"},
				},
			},
			{
				BlockType:   "Tabata",
				BlockRounds: 8,
				BlockNotes:  "20s work / 10s rest x 8 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Ab Twists"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Flutter Kicks"], Duration: "20s", Rest: "10s"},
				},
			},
		},
	}
	createWorkoutDay(day10)

	// --- DAY 11 ---
	day11 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   11,
		Title:       "Full Body Cardio",
		Description: "As many rounds as possible (AMRAP) in 15 minutes. Track your rounds with the counter!",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "AMRAP",
				BlockNotes: "15 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpees (modified)"], Reps: "3"},
					{Order: 2, ExerciseID: exIDs["Squat Twists"], Reps: "6"},
					{Order: 3, ExerciseID: exIDs["Press Ups (on Knees)"], Reps: "9"},
					{Order: 4, ExerciseID: exIDs["High Knees"], Reps: "12", Tips: "2 High Knees = 1 rep"},
				},
			},
		},
	}
	createWorkoutDay(day11)

	// --- DAY 12 ---
	day12 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   12,
		Title:       "Plyometrics",
		Description: "A plyometric circuit to build explosive power.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 3,
				BlockNotes:  "Exercise for 45 seconds and then rest for 15 seconds. Repeat the circuit 3 times. Full duration 21 minutes.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Overhead Jabs"], Duration: "45s", Rest: "15s"},
					{Order: 2, ExerciseID: exIDs["Broad Jumps"], Duration: "45s", Rest: "15s"},
					{Order: 3, ExerciseID: exIDs["Squat Kicks"], Duration: "45s", Rest: "15s"},
					{Order: 4, ExerciseID: exIDs["Standing Mountain Climbers"], Duration: "45s", Rest: "15s"},
					{Order: 5, ExerciseID: exIDs["Calf Jumps"], Duration: "45s", Rest: "15s"},
					{Order: 6, ExerciseID: exIDs["Sprints"], Duration: "45s", Rest: "15s"},
					{Order: 7, ExerciseID: exIDs["Oblique Hops"], Duration: "45s", Rest: "15s"},
				},
			},
		},
	}
	createWorkoutDay(day12)

	// --- DAY 13 ---
	day13 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   13,
		Title:       "Core Focus & Full Body Cardio",
		Description: "A core circuit followed by a cardio circuit.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 3,
				BlockNotes:  "Exercise for 30 seconds and then rest for 15 seconds. Repeat the circuit 3 times. Full duration 9 minutes.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Half Sit Ups"], Duration: "30s", Rest: "15s"},
					{Order: 2, ExerciseID: exIDs["Scissors"], Duration: "30s", Rest: "15s"},
					{Order: 3, ExerciseID: exIDs["Flutter Kicks"], Duration: "30s", Rest: "15s"},
					{Order: 4, ExerciseID: exIDs["Elbows to Knee"], Duration: "30s", Rest: "15s"},
				},
			},
			{
				BlockType:   "Circuit",
				BlockRounds: 3,
				BlockNotes:  "Exercise for 45 seconds and then rest for 15 seconds. Repeat the circuit 3 times. Full duration 12 minutes.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpees"], Duration: "45s", Rest: "15s"},
					{Order: 2, ExerciseID: exIDs["Starjumps"], Duration: "45s", Rest: "15s"},
					{Order: 3, ExerciseID: exIDs["T-Runs"], Duration: "45s", Rest: "15s"},
					{Order: 4, ExerciseID: exIDs["Switch Kicks"], Duration: "45s", Rest: "15s"},
				},
			},
		},
	}
	createWorkoutDay(day13)

	// --- DAY 14 ---
	day14 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   14,
		Title:       "Full Body Switch Up",
		Description: "Every Minute on the Minute (EMOM), complete the number of reps within the minute. Faster you complete the more rest",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "EMOM",
				BlockRounds: 4,
				BlockNotes:  "16 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, Instructions: "Minute 1", ExerciseID: exIDs["H.O.G. Press Ups (on Knees)"], Reps: "20"},
					{Order: 2, Instructions: "Minute 2", ExerciseID: exIDs["Squats"], Reps: "20"},
					{Order: 3, Instructions: "Minute 3", ExerciseID: exIDs["Diamond Press Ups (on Knees)"], Reps: "20"},
					{Order: 4, Instructions: "Minute 4", ExerciseID: exIDs["Squat Hold"], Duration: "30s"},
				},
			},
		},
	}
	createWorkoutDay(day14)

	// --- DAY 15 ---
	day15 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   15,
		Title:       "Recovery day & Optional workout",
		Description: "Mobility AND/OR Workout: Power circuit focusing on endurance.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Mobility",
				BlockRounds: 1,
				BlockNotes:  "A mobility session to stretch your tight muscle. Prevent injury and aid recovery",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Mobility"]},
				},
			},
			{
				BlockType:   "Circuit",
				BlockRounds: 4,
				BlockNotes:  "Exercise for 35 seconds and then rest for 25 seconds. Repeat the circuit 4 times. Full duration 20 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Jumps"], Duration: "35s", Rest: "25s"},
					{Order: 2, ExerciseID: exIDs["Press Ups"], Duration: "35s", Rest: "25s"},
					{Order: 3, ExerciseID: exIDs["Mountain Climbers"], Duration: "35s", Rest: "25s"},
					{Order: 4, ExerciseID: exIDs["Lunges"], Duration: "35s", Rest: "25s"},
					{Order: 5, ExerciseID: exIDs["Burpees"], Duration: "35s", Rest: "25s"},
					{Order: 6, ExerciseID: exIDs["Starjumps"], Duration: "35s", Rest: "25s"},
				},
			},
		},
	}
	createWorkoutDay(day15)

	// --- DAY 16 ---
	day16 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   16,
		Title:       "Upper Body Challenge",
		Description: "Every Minute on the Minute (EMOM), complete the number of reps within the minute. Faster you complete the more rest",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "EMOM",
				BlockRounds: 4,
				BlockNotes:  "16 minutes total",
				Exercises: []models.WorkoutExercise{
					{Order: 1, Instructions: "Minute 1", ExerciseID: exIDs["Press Ups (on Knees)"], Reps: "15"},
					{Order: 2, Instructions: "Minute 2", ExerciseID: exIDs["Tricep Dips (with Chair)"], Reps: "15"},
					{Order: 3, Instructions: "Minute 3", ExerciseID: exIDs["V Press Ups"], Reps: "15"},
					{Order: 4, Instructions: "Minute 4", ExerciseID: exIDs["Plank Hold"], Reps: "30s"},
				},
			},
		},
	}
	createWorkoutDay(day16)

	// --- DAY 17 ---
	day17 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   17,
		Title:       "Lower Body Endurance",
		Description: "A pyramid workout for lower body endurance. Complete all exercises from 3 reps to 7 reps and back to 3 reps",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "For Time",
				BlockNotes: "Pyramid: 3-4-5-6-7-6-5-4-3 reps.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Jumps"], Reps: "3, 4, 5, 6, 7, 6, 5, 4, 3"},
					{Order: 2, ExerciseID: exIDs["Lunges"], Reps: "3, 4, 5, 6, 7, 6, 5, 4, 3", Tips: "2 Lunges = 1 rep"},
					{Order: 3, ExerciseID: exIDs["Calf Jumps"], Reps: "3, 4, 5, 6, 7, 6, 5, 4, 3"},
					{Order: 4, ExerciseID: exIDs["Box Jumps"], Reps: "3, 4, 5, 6, 7, 6, 5, 4, 3", Tips: "1 Full box = 1 rep"},
				},
			},
		},
	}
	createWorkoutDay(day17)

	// --- DAY 18 ---
	day18 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   18,
		Title:       "Metabolic Conditioning",
		Description: "A metabolic conditioning workout with multiple Tabata rounds.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Tabata",
				BlockRounds: 8,
				BlockNotes:  "20s work / 10s rest x 8 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpees"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["High Knees"], Duration: "20s", Rest: "10s"},
				},
			},
			{
				BlockType:   "Tabata",
				BlockRounds: 8,
				BlockNotes:  "20s work / 10s rest x 8 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Mountain Climbers"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Squat Jumps"], Duration: "20s", Rest: "10s"},
				},
			},
			{
				BlockType:   "Tabata",
				BlockRounds: 8,
				BlockNotes:  "20s work / 10s rest x 8 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Starjumps"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Sprints"], Duration: "20s", Rest: "10s"},
				},
			},
		},
	}
	createWorkoutDay(day18)

	// --- DAY 19 ---
	day19 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   19,
		Title:       "Full Body Strength",
		Description: "A full body strength circuit with sets and reps.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 5,
				BlockNotes:  "Exercise for 40 seconds and then rest for 20 seconds. Repeat the circuit 5 times. Full duration 25 minutes.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squats"], Duration: "40s", Rest: "20s"},
					{Order: 2, ExerciseID: exIDs["Press Ups"], Duration: "40s", Rest: "20s"},
					{Order: 3, ExerciseID: exIDs["Lunges"], Duration: "40s", Rest: "20s"},
					{Order: 4, ExerciseID: exIDs["Tricep Dips (with Chair)"], Duration: "40s", Rest: "20s"},
					{Order: 5, ExerciseID: exIDs["Plank Shoulder Taps"], Duration: "40s", Rest: "20s"},
				},
			},
		},
	}
	createWorkoutDay(day19)

	// --- DAY 20 ---
	day20 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   20,
		Title:       "Core Intensive",
		Description: "As many rounds as possible (AMRAP) in 18 minutes. Track your rounds with the counter!",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "AMRAP",
				BlockNotes: "18 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Sit Ups"], Reps: "10"},
					{Order: 2, ExerciseID: exIDs["Bicycles"], Reps: "15", Tips: "2 Bicycles = 1 rep"},
					{Order: 3, ExerciseID: exIDs["Flutter Kicks"], Reps: "20", Tips: "2 Kicks = 1 rep"},
					{Order: 4, ExerciseID: exIDs["Ab Twists"], Reps: "25"},
					{Order: 5, ExerciseID: exIDs["Plank Shoulder Taps"], Reps: "15", Tips: "2 Taps = 1 rep"},
				},
			},
		},
	}
	createWorkoutDay(day20)

	// --- DAY 21 ---
	day21 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   21,
		Title:       "Cardio Finisher",
		Description: "A descending ladder workout. Complete all exercises from 10 reps down to 1 rep.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "For Time",
				BlockNotes: "Descending ladder: 10-9-8-7-6-5-4-3-2-1 reps.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpees"], Reps: "10, 9, 8, 7, 6, 5, 4, 3, 2, 1"},
					{Order: 2, ExerciseID: exIDs["Squat Jumps"], Reps: "10, 9, 8, 7, 6, 5, 4, 3, 2, 1"},
					{Order: 3, ExerciseID: exIDs["Press Ups"], Reps: "10, 9, 8, 7, 6, 5, 4, 3, 2, 1"},
					{Order: 4, ExerciseID: exIDs["Starjumps"], Reps: "10, 9, 8, 7, 6, 5, 4, 3, 2, 1"},
				},
			},
		},
	}
	createWorkoutDay(day21)

	// --- DAY 22 ---
	day22 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   22,
		Title:       "Recovery day & Optional workout",
		Description: "Mobility AND/OR Workout: Full body circuit.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Mobility",
				BlockRounds: 1,
				BlockNotes:  "A mobility session to stretch your tight muscle. Prevent injury and aid recovery",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Mobility"]},
				},
			},
			{
				BlockType:   "Circuit",
				BlockRounds: 4,
				BlockNotes:  "Exercise for 40 seconds and then rest for 20 seconds. Repeat the circuit 3 times. Full duration 16 minutes.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Ups (on Knees)"], Duration: "40s", Rest: "20s"},
					{Order: 2, ExerciseID: exIDs["Squats"], Duration: "40s", Rest: "20s"},
					{Order: 3, ExerciseID: exIDs["Plank Hold"], Duration: "40s", Rest: "20s"},
					{Order: 4, ExerciseID: exIDs["Walkaways"], Duration: "40s", Rest: "20s"},
				},
			},
		},
	}
	createWorkoutDay(day22)

	// --- DAY 23 ---
	day23 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   23,
		Title:       "Upper Body & Core Challenge",
		Description: "An EMOM circuit followed by a core circuit.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "EMOM",
				BlockRounds: 4,
				BlockNotes:  "12 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, Instructions: "Minute 1", ExerciseID: exIDs["Moving Press Ups"], Reps: "8"},
					{Order: 2, Instructions: "Minute 2", ExerciseID: exIDs["Oblique Plank"], Reps: "10"},
					{Order: 3, Instructions: "Minute 3", ExerciseID: exIDs["Bearcrawls"], Reps: "4", Tips: "1 x forward and 1 x backward = 1 rep"},
				},
			},
			{
				BlockType:   "Circuit",
				BlockRounds: 3,
				BlockNotes:  "Exercise for 30 seconds and then rest for 15 seconds. Repeat the circuit 3 times. Full duration 9 minutes.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Toe Taps"], Duration: "30s", Rest: "15s"},
					{Order: 2, ExerciseID: exIDs["Straddle Sit Ups"], Duration: "30s", Rest: "15s"},
					{Order: 3, ExerciseID: exIDs["Scissors"], Duration: "30s", Rest: "15s"},
					{Order: 4, ExerciseID: exIDs["Ab Twists"], Duration: "30s", Rest: "15s"},
				},
			},
		},
	}
	createWorkoutDay(day23)

	// --- DAY 24 ---
	day24 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   24,
		Title:       "Full Body Tabata",
		Description: "A full body workout with multiple Tabata blocks.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Tabata",
				BlockRounds: 8,
				BlockNotes:  "20s work / 10s rest x 8 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squats"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Starjumps"], Duration: "20s", Rest: "10s"},
				},
			},
			{
				BlockType:   "Tabata",
				BlockRounds: 8,
				BlockNotes:  "20s work / 10s rest x 8 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Ups"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Mountain Climbers"], Duration: "20s", Rest: "10s"},
				},
			},
			{
				BlockType:   "Tabata",
				BlockRounds: 8,
				BlockNotes:  "20s work / 10s rest x 8 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Lunges"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["High Knees"], Duration: "20s", Rest: "10s"},
				},
			},
			{
				BlockType:   "Tabata",
				BlockRounds: 8,
				BlockNotes:  "20s work / 10s rest x 8 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Elbows to Knee"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Ski Jumps"], Duration: "20s", Rest: "10s"},
				},
			},
		},
	}
	createWorkoutDay(day24)

	// --- DAY 25 ---
	day25 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   25,
		Title:       "Endurance Challenge",
		Description: "As many rounds as possible (AMRAP) in 20 minutes. Track your rounds with the counter!",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "AMRAP",
				BlockNotes: "20 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpees"], Reps: "5"},
					{Order: 2, ExerciseID: exIDs["Squat Jumps"], Reps: "10"},
					{Order: 3, ExerciseID: exIDs["Press Ups"], Reps: "15"},
					{Order: 4, ExerciseID: exIDs["Mountain Climbers"], Reps: "20", Tips: "2 Climberss = 1 rep"},
					{Order: 5, ExerciseID: exIDs["Starjumps"], Reps: "25"},
				},
			},
		},
	}
	createWorkoutDay(day25)

	// --- DAY 26 ---
	day26 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   26,
		Title:       "Agility",
		Description: "A timed agility circuit.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 4,
				RoundRest:   "60s",
				BlockNotes:  "Exercise for 50 seconds and then rest for 10 seconds. Repeat the circuit 4 times. Rest 60 seconds between rounds. Full duration 24 minutes.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Sprints"], Duration: "50s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Sprawls"], Duration: "50s", Rest: "10s"},
					{Order: 3, ExerciseID: exIDs["T-Runs"], Duration: "50s", Rest: "10s"},
					{Order: 4, ExerciseID: exIDs["Ski Jumps"], Duration: "50s", Rest: "10s"},
					{Order: 5, ExerciseID: exIDs["Box Jumps"], Duration: "50s", Rest: "10s"},
				},
			},
		},
	}
	createWorkoutDay(day26)

	// --- DAY 27 ---
	day27 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   27,
		Title:       "Full Body Conditioning",
		Description: "Every Minute on the Minute (EMOM), complete the number of reps within the minute. Faster you complete the more rest",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "EMOM",
				BlockRounds: 4,
				BlockNotes:  "20 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, Instructions: "Minute 1", ExerciseID: exIDs["Squat Jumps"], Reps: "12"},
					{Order: 2, Instructions: "Minute 2", ExerciseID: exIDs["Press Ups (on Knees)"], Reps: "12"},
					{Order: 3, Instructions: "Minute 3", ExerciseID: exIDs["Lunges"], Reps: "10", Tips: "2 Lunges = 1 rep"},
					{Order: 4, Instructions: "Minute 4", ExerciseID: exIDs["Mountain Climbers"], Reps: "20", Tips: "2 Climbers = 1 rep"},
					{Order: 5, Instructions: "Minute 5", ExerciseID: exIDs["Jack Knife"], Reps: "10"},
				},
			},
		},
	}
	createWorkoutDay(day27)

	// --- DAY 28 ---
	day28 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   28,
		Title:       "Full Body Workout",
		Description: "A pyramid workout with a mix of exercises. Complete all exercises from 1 to 10reps and back down to 1 rep",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "For Time",
				BlockNotes: "Pyramid: 1-2-3-4-5-6-7-8-9-10-9-8-7-6-5-4-3-2-1 reps.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Jumps"], Reps: "1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1"},
					{Order: 2, ExerciseID: exIDs["Press Ups"], Reps: "1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1"},
					{Order: 3, ExerciseID: exIDs["Plank Jabs"], Reps: "1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1", Tips: "2 Jabs = 1 rep"},
					{Order: 4, ExerciseID: exIDs["Reverse Lunge"], Reps: "1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1", Tips: "2 Lunges = 1 rep"},
					{Order: 5, ExerciseID: exIDs["Diamond Sit Ups"], Reps: "1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1"},
				},
			},
		},
	}
	createWorkoutDay(day28)

	// --- DAY 29 ---
	day29 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   29,
		Title:       "Recovery day & Optional workout",
		Description: "Mobility AND/OR AMRAP (As many rounds as possible ) in 20 minutes. Track your rounds with the counter!",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Mobility",
				BlockRounds: 1,
				BlockNotes:  "A mobility session to stretch your tight muscle. Prevent injury and aid recovery",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Mobility"]},
				},
			},
			{
				BlockType:   "AMRAP",
				BlockRounds: 2,
				BlockNotes:  "20 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["High Knees"], Reps: "20", Tips: "2 High Knees = 1 rep"},
					{Order: 2, ExerciseID: exIDs["Cross Jabs"], Reps: "20", Tips: "2 Jabs = 1 rep"},
					{Order: 3, ExerciseID: exIDs["Diamond Sit Ups"], Reps: "20"},
					{Order: 4, ExerciseID: exIDs["Belt Kicks"], Reps: "10", Tips: "2 Belt Kicks = 1 rep"},
				},
			},
		},
	}
	createWorkoutDay(day29)

	// --- DAY 30 ---
	day30 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   30,
		Title:       "FINAL FITNESS Assessment",
		Description: "Complete this fitness assessment one more time and compare the results from Day 1.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "Fitness Assessment",
				BlockNotes: "Push yourself as hard as you did on day 1 and note your improvements.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Ups"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 2, ExerciseID: exIDs["Squats"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 3, ExerciseID: exIDs["Plank Hold"], Reps: "Max Effort", Duration: "Max Time", Rest: "2 mins"},
					{Order: 4, ExerciseID: exIDs["Burpees"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 5, ExerciseID: exIDs["Starjumps"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 6, ExerciseID: exIDs["Sit Ups"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 7, ExerciseID: exIDs["Lunges"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins", Tips: "2 Lunges = 1 rep"},
					{Order: 8, ExerciseID: exIDs["Tricep Dips (with Chair)"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
				},
			},
		},
	}
	createWorkoutDay(day30)

	log.Println("Successfully seeded Beginner Workout data.")
}
