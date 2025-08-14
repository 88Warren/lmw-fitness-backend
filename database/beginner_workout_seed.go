package database

import (
	"log"

	"github.com/88warren/lmw-fitness-backend/models"
	"gorm.io/gorm"
)

// Helper function to get an exercise ID by name
// In a real application, this would be a more robust lookup.
func getExerciseIDByName(db *gorm.DB, name string) (uint, error) {
	var exercise models.Exercise
	if err := db.Where("name = ?", name).First(&exercise).Error; err != nil {
		return 0, err
	}
	return exercise.ID, nil
}

// Helper function to get the program ID by name
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

	// Fetch all necessary exercise IDs upfront
	exIDs := make(map[string]uint)
	exercises := []string{
		"Press Ups", "Squats", "Plank Hold", "Burpees", "Starjumps", "Diamond Sit Ups",
		"Press Ups (on Knees)", "Crunches", "Glute Bridges", "Lunges", "Calf Raises",
		"Leg Raises", "Donkey Kicks", "Wide Arm Press Ups (on Knees)", "Tricep Dips",
		"Plank Shoulder Taps", "Superman", "Walkaways", "Cross Jabs", "Dorsal Raises",
		"Mountain Climbers", "Squat Jumps", "Ab Twists", "Bicycles", "Flutter Kicks",
		"Half Sit Ups", "Heel Taps", "Bearcrawls", "T-Runs", "Y Shaped Lunges", "H.O.G. Press Ups",
		"Moving Press Ups", "Scissors", "V Press Ups", "Straddle Sit Ups", "Sit Ups",
		"Broad Jumps", "Squat Kicks", "Squat Twists", "Switch Kicks", "Standing Mountain Climbers",
		"Sprints", "Sprawls", "Ski Jumps", "Plank Jabs", "Overhead Jabs", "Reverse Lunge",
		"Oblique Hops", "Oblique Plank", "Lateral Lunges", "Plank Leg Raises", "Knees to Chest",
		"High Knees", "Jack Knife", "Elbows to Knee", "Calf Jumps", "Box Jumps", "Belt Kicks",
		"Toe Taps",
	}

	for _, name := range exercises {
		id, err := getExerciseIDByName(DB, name)
		if err != nil {
			log.Fatalf("Failed to find exercise '%s': %v", name, err)
		}
		exIDs[name] = id
	}

	// Helper function to create a workout day and handle errors for BEGINNER program
	createWorkoutDay := func(day models.WorkoutDay) {
		var existingDay models.WorkoutDay
		if err := DB.Where("program_id = ? AND day_number = ?", day.ProgramID, day.DayNumber).First(&existingDay).Error; err == nil {
			// Log for beginner program
			log.Printf("Beginner Program - Day %d already exists, skipping creation.", day.DayNumber)
			return
		}
		if err := DB.Create(&day).Error; err != nil {
			// Log for beginner program
			log.Printf("Failed to create Beginner Program - Day %d: %v", day.DayNumber, err)
		} else {
			// Log for beginner program
			log.Printf("Successfully created Beginner Program - Day %d: %s", day.DayNumber, day.Title)
		}
	}

	// --- DAY 1 ---
	day1 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   1,
		Title:       "FITNESS TEST + Full Body Introduction",
		Description: "After the fitness test, complete a full body circuit.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "Fitness Test",
				BlockNotes: "1 min work – max effort, 2 mins rest for each exercise. Make sure you record your results you will need them for day 30",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Ups"], Reps: "Max Effort", Duration: "1 min", RestDuration: "2 mins"},
					{Order: 2, ExerciseID: exIDs["Squats"], Reps: "Max Effort", Duration: "1 min", RestDuration: "2 mins"},
					{Order: 3, ExerciseID: exIDs["Plank Hold"], Reps: "Max Effort", Duration: "1 min", RestDuration: "2 mins"},
					{Order: 4, ExerciseID: exIDs["Burpees"], Reps: "Max Effort", Duration: "1 min", RestDuration: "2 mins"},
					{Order: 5, ExerciseID: exIDs["Starjumps"], Reps: "Max Effort", Duration: "1 min", RestDuration: "2 mins"},
					{Order: 6, ExerciseID: exIDs["Diamond Sit Ups"], Reps: "Max Effort", Duration: "1 min", RestDuration: "2 mins"},
				},
			},
			{
				BlockType:   "Circuit",
				BlockRounds: 2,
				BlockNotes:  "90 seconds rest between rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squats"], Reps: "10 reps"},
					{Order: 2, ExerciseID: exIDs["Press Ups (on Knees)"], Reps: "10 reps"},
					{Order: 3, ExerciseID: exIDs["Crunches"], Reps: "15 reps"},
					{Order: 4, ExerciseID: exIDs["Glute Bridges"], Reps: "15 reps"},
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
		Description: "A circuit to build lower body strength and endurance.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 3,
				BlockNotes:  "45 seconds work, 15 seconds rest.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squats"], Duration: "45s", RestDuration: "15s"},
					{Order: 2, ExerciseID: exIDs["Lunges"], Duration: "45s", RestDuration: "15s"},
					{Order: 3, ExerciseID: exIDs["Calf Raises"], Duration: "45s", RestDuration: "15s"},
					{Order: 4, ExerciseID: exIDs["Glute Bridges"], Duration: "45s", RestDuration: "15s"},
					{Order: 5, ExerciseID: exIDs["Leg Raises"], Duration: "45s", RestDuration: "15s"},
					{Order: 6, ExerciseID: exIDs["Donkey Kicks"], Duration: "45s", RestDuration: "15s"},
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
		Description: "A circuit to build upper body strength.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 3,
				BlockNotes:  "30 seconds work, 15 seconds rest, with 60 seconds rest between rounds.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Wide Arm Press Ups (on Knees)"], Duration: "30s", RestDuration: "15s"},
					{Order: 2, ExerciseID: exIDs["Tricep Dips"], Duration: "30s", RestDuration: "15s", Tips: "Can do on chair, sofa or step if you have one."},
					{Order: 3, ExerciseID: exIDs["Plank Shoulder Taps"], Duration: "30s", RestDuration: "15s", Tips: "Can do on knees"},
					{Order: 4, ExerciseID: exIDs["Superman"], Duration: "30s", RestDuration: "15s"},
					{Order: 6, ExerciseID: exIDs["Plank Hold"], Duration: "30s", RestDuration: "15s"},
					{Order: 6, ExerciseID: exIDs["Walkaways"], Duration: "30s", RestDuration: "15s"},
					{Order: 7, ExerciseID: exIDs["Cross Jabs"], Duration: "30s", RestDuration: "15s"},
					{Order: 8, ExerciseID: exIDs["Dorsal Raises"], Duration: "30s", RestDuration: "15s"},
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
		Description: "As many rounds as possible (AMRAP) in 12 minutes.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "AMRAP",
				BlockNotes: "12 minutes total.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpees"], Reps: "5", Tips: "Modified if needed"},
					{Order: 2, ExerciseID: exIDs["Mountain Climbers"], Reps: "10"},
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
		Description: "Every Minute on the Minute (EMOM) for 4 rounds.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "EMOM",
				BlockRounds: 4,
				BlockNotes:  "Complete the designated reps at the top of each minute. If you finish early, rest until the next minute starts.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squats"], Reps: "10", Tips: "Minute 1"},
					{Order: 2, ExerciseID: exIDs["Press Ups (on Knees)"], Reps: "8", Tips: "Minute 2"},
					{Order: 3, ExerciseID: exIDs["Lunges"], Reps: "6 per leg", Tips: "Minute 3"},
					{Order: 4, ExerciseID: exIDs["Crunches"], Reps: "15", Tips: "Minute 4"},
					{Order: 5, ExerciseID: exIDs["Starjumps"], Reps: "10", Tips: "Minute 5 (low impact)"},
				},
			},
		},
	}
	createWorkoutDay(day5)

	// --- DAY 6 ---
	day6 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   6,
		Title:       "Active Recovery & Core",
		Description: "A timed core workout for recovery and stability.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 4,
				BlockNotes:  "30 seconds each exercise. 60 Second rest between rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Leg Raises"], Duration: "30s"},
					{Order: 2, ExerciseID: exIDs["Bicycles"], Duration: "30s"},
					{Order: 3, ExerciseID: exIDs["Flutter Kicks"], Duration: "30s"},
					{Order: 4, ExerciseID: exIDs["Half Sit Ups"], Duration: "30s"},
					{Order: 5, ExerciseID: exIDs["Heel Taps"], Duration: "30s"},
					{Order: 6, ExerciseID: exIDs["Glute Bridges"], Duration: "30s"},
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
		Description: "An EMOM workout alternating between two exercise sets for 16 minutes.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "EMOM",
				BlockNotes: "16 minutes total, alternating between odd and even minutes.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, Tips: "Odd minutes", Reps: "10 Squats + 5 Press Ups (on Knees)", ExerciseID: exIDs["Squats"]},
					{Order: 2, Tips: "Even minutes", Reps: "30-second plank + 10 Calf Raises", ExerciseID: exIDs["Plank Hold"]},
				},
			},
		},
	}
	createWorkoutDay(day7)

	// --- DAY 8 ---
	day8 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   8,
		Title:       "Upper Body Strength & Mobility",
		Description: "Optional mobility day AND/OR a regular workout focusing on upper body strength.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 3,
				BlockNotes:  "30 seconds rest between sets.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Ups"], Reps: "10", Tips: "Knee or full"},
					{Order: 2, ExerciseID: exIDs["Tricep Dips"], Reps: "15"},
					{Order: 3, ExerciseID: exIDs["Walkaways"], Reps: "8"},
					{Order: 4, ExerciseID: exIDs["Jack Knife"], Reps: "12"},
					{Order: 5, ExerciseID: exIDs["Plank Leg Raises"], Reps: "10 per leg"},
					{Order: 6, ExerciseID: exIDs["Diamond Sit Ups"], Reps: "15"},
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
				BlockNotes:  "40 seconds work, 20 seconds rest.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Kicks"], Duration: "40s", RestDuration: "20s"},
					{Order: 2, ExerciseID: exIDs["Y Shaped Lunges"], Duration: "40s", RestDuration: "20s"},
					{Order: 3, ExerciseID: exIDs["Calf Jumps"], Duration: "40s", RestDuration: "20s"},
					{Order: 4, ExerciseID: exIDs["Squat Jumps"], Duration: "40s", RestDuration: "20s"},
					{Order: 5, ExerciseID: exIDs["Lateral Lunges"], Duration: "40s", RestDuration: "20s"},
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
				BlockRounds: 5,
				BlockNotes:  "8 sets, 20s on, 10s off. Alternate between exercises in each round.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Mountain Climbers"], Reps: "Tabata Round 1"},
					{Order: 2, ExerciseID: exIDs["Plank Hold"], Reps: "Tabata Round 1"},
					{Order: 3, ExerciseID: exIDs["Bicycles"], Reps: "Tabata Round 2"},
					{Order: 4, ExerciseID: exIDs["Leg Raises"], Reps: "Tabata Round 2"},
					{Order: 5, ExerciseID: exIDs["Ab Twists"], Reps: "Tabata Round 3"},
					{Order: 6, ExerciseID: exIDs["Flutter Kicks"], Reps: "Tabata Round 3"},
					{Order: 7, ExerciseID: exIDs["Crunches"], Reps: "Tabata Round 4"},
					{Order: 8, ExerciseID: exIDs["Scissors"], Reps: "Tabata Round 4"},
					{Order: 9, ExerciseID: exIDs["Superman"], Reps: "Tabata Round 5"},
					{Order: 10, ExerciseID: exIDs["Half Sit Ups"], Reps: "Tabata Round 5"},
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
		Description: "As many rounds as possible (AMRAP) in 15 minutes.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "AMRAP",
				BlockNotes: "15 minutes total.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpees"], Reps: "3"},
					{Order: 2, ExerciseID: exIDs["Squat Twists"], Reps: "6"},
					{Order: 3, ExerciseID: exIDs["Press Ups"], Reps: "9", Tips: "Knee or full"},
					{Order: 4, ExerciseID: exIDs["High Knees"], Reps: "12"},
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
				BlockNotes:  "45 seconds work, 15 seconds rest.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Overhead Jabs"], Duration: "45s", RestDuration: "15s"},
					{Order: 2, ExerciseID: exIDs["Broad Jumps"], Duration: "45s", RestDuration: "15s"},
					{Order: 3, ExerciseID: exIDs["Squat Kicks"], Duration: "45s", RestDuration: "15s"},
					{Order: 4, ExerciseID: exIDs["Standing Mountain Climbers"], Duration: "45s", RestDuration: "15s"},
					{Order: 5, ExerciseID: exIDs["Calf Jumps"], Duration: "45s", RestDuration: "15s"},
					{Order: 6, ExerciseID: exIDs["Sprints"], Duration: "45s", RestDuration: "15s"},
					{Order: 7, ExerciseID: exIDs["Oblique Hops"], Duration: "45s", RestDuration: "15s"},
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
				BlockNotes:  "Core Circuit: 30s work, 15s rest.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Half Sit Ups"], Duration: "30s", RestDuration: "15s"},
					{Order: 2, ExerciseID: exIDs["Scissors"], Duration: "30s", RestDuration: "15s"},
					{Order: 3, ExerciseID: exIDs["Flutter Kicks"], Duration: "30s", RestDuration: "15s"},
					{Order: 4, ExerciseID: exIDs["Elbows to Knee"], Duration: "30s", RestDuration: "15s"},
				},
			},
			{
				BlockType:   "Circuit",
				BlockRounds: 3,
				BlockNotes:  "Cardio Circuit: 45s work, 15s rest.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpees"], Duration: "45s", RestDuration: "15s"},
					{Order: 2, ExerciseID: exIDs["Starjumps"], Duration: "45s", RestDuration: "15s"},
					{Order: 3, ExerciseID: exIDs["T-Runs"], Duration: "45s", RestDuration: "15s"},
					{Order: 4, ExerciseID: exIDs["Switch Kicks"], Duration: "45s", RestDuration: "15s"},
				},
			},
		},
	}
	createWorkoutDay(day13)

	// --- DAY 14 ---
	day14 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   14,
		Title:       "EMOM 16 minutes (alternate)",
		Description: "An EMOM workout alternating between two exercises.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "EMOM",
				BlockNotes: "16 minutes total. Odd minutes: HOG Press Ups (20 reps). Even minutes: Squats (20 reps). Rest remainder of minute.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["H.O.G. Press Ups"], Reps: "20", Tips: "Odd minutes"},
					{Order: 2, ExerciseID: exIDs["Squats"], Reps: "20", Tips: "Even minutes"},
				},
			},
		},
	}
	createWorkoutDay(day14)

	// --- DAY 15 ---
	day15 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   15,
		Title:       "Power Circuit & Mobility",
		Description: "Optional mobility AND/OR a power circuit focusing on endurance.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 4,
				BlockNotes:  "35s work, 25s rest.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Jumps"], Duration: "35s", RestDuration: "25s"},
					{Order: 2, ExerciseID: exIDs["Press Ups"], Duration: "35s", RestDuration: "25s", Tips: "Modified"},
					{Order: 3, ExerciseID: exIDs["Mountain Climbers"], Duration: "35s", RestDuration: "25s"},
					{Order: 4, ExerciseID: exIDs["Lunges"], Duration: "35s", RestDuration: "25s"},
					{Order: 5, ExerciseID: exIDs["Burpees"], Duration: "35s", RestDuration: "25s", Tips: "Modified"},
					{Order: 6, ExerciseID: exIDs["Starjumps"], Duration: "35s", RestDuration: "25s"},
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
		Description: "An EMOM workout targeting the upper body.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "EMOM",
				BlockRounds: 3,
				BlockNotes:  "12 minutes total. 45s work, then rest until the next minute starts.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Ups"], Reps: "Max", Duration: "45s", Tips: "Minute 1"},
					{Order: 2, ExerciseID: exIDs["Tricep Dips"], Reps: "Max", Duration: "45s", Tips: "Minute 2"},
					{Order: 3, ExerciseID: exIDs["V Press Ups"], Reps: "Max", Duration: "45s", Tips: "Minute 3"},
					{Order: 4, ExerciseID: exIDs["Plank Hold"], Duration: "45s", Tips: "Minute 4"},
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
		Description: "A ladder workout for lower body endurance.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "Ladder",
				BlockNotes: "Pyramid up and down: 3-4-5-6-7-6-5-4-3 reps.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Jumps"]},
					{Order: 2, ExerciseID: exIDs["Lunges"], Tips: "Each leg / Jumps"},
					{Order: 3, ExerciseID: exIDs["Calf Jumps"], Tips: "Raises"},
					{Order: 4, ExerciseID: exIDs["Box Jumps"]},
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
				BlockRounds: 2,
				BlockNotes:  "6 Tabata rounds of each exercise (20s work, 10s rest). Rest 90 seconds between rounds.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpees"], Reps: "Round 1"},
					{Order: 2, ExerciseID: exIDs["Mountain Climbers"], Reps: "Round 2"},
					{Order: 3, ExerciseID: exIDs["Starjumps"], Reps: "Round 3"},
					{Order: 4, ExerciseID: exIDs["High Knees"], Reps: "Round 4"},
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
				BlockNotes:  "10 reps each. 60 seconds rest between rounds.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squats"], Reps: "10"},
					{Order: 2, ExerciseID: exIDs["Press Ups"], Reps: "10"},
					{Order: 3, ExerciseID: exIDs["Lunges"], Reps: "10 per leg"},
					{Order: 4, ExerciseID: exIDs["Tricep Dips"], Reps: "10"},
					{Order: 5, ExerciseID: exIDs["Plank Shoulder Taps"], Reps: "10 per side"},
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
		Description: "An AMRAP focused on core strength.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "AMRAP",
				BlockNotes: "18 minutes total. Complete as many rounds as possible.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Sit Ups"], Reps: "10"},
					{Order: 2, ExerciseID: exIDs["Bicycles"], Reps: "15 each side"},
					{Order: 3, ExerciseID: exIDs["Flutter Kicks"], Reps: "20"},
					{Order: 4, ExerciseID: exIDs["Ab Twists"], Reps: "25"},
					{Order: 5, ExerciseID: exIDs["Plank Hold"], Duration: "30s"},
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
		Description: "A descending ladder workout.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "Ladder",
				BlockNotes: "Descending ladder (10,9,8...1). Complete all exercises at each number.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpees"]},
					{Order: 2, ExerciseID: exIDs["Squat Jumps"]},
					{Order: 3, ExerciseID: exIDs["Press Ups"]},
					{Order: 4, ExerciseID: exIDs["Starjumps"]},
				},
			},
		},
	}
	createWorkoutDay(day21)

	// --- DAY 22 ---
	day22 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   22,
		Title:       "Strength Test Prep & Mobility",
		Description: "Optional mobility AND/OR a workout to prepare for a strength test.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "Strength Test",
				BlockNotes: "Light movement and stretching.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Ups"], Reps: "5 sets of max reps"},
					{Order: 2, ExerciseID: exIDs["Squats"], Reps: "5 sets of 15-20 reps"},
					{Order: 3, ExerciseID: exIDs["Plank Hold"], Reps: "3 max effort holds"},
					{Order: 4, ExerciseID: exIDs["Walkaways"], Reps: "5 sets of 10 reps"},
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
		Description: "An EMOM circuit followed by a timed core workout.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "EMOM",
				BlockRounds: 4,
				BlockNotes:  "4 rounds total.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Moving Press Ups"], Reps: "8", Tips: "Minute 1"},
					{Order: 2, ExerciseID: exIDs["Oblique Plank"], Reps: "10", Tips: "Minute 2"},
					{Order: 3, ExerciseID: exIDs["Bearcrawls"], Reps: "Forward + back x 4", Tips: "Minute 3"},
				},
			},
			{
				BlockType:   "Circuit",
				BlockRounds: 3,
				BlockNotes:  "Timed Core: 30s work, 15s rest.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Toe Taps"], Duration: "30s", RestDuration: "15s"},
					{Order: 2, ExerciseID: exIDs["Straddle Sit Ups"], Duration: "30s", RestDuration: "15s"},
					{Order: 3, ExerciseID: exIDs["Scissors"], Duration: "30s", RestDuration: "15s"},
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
		Description: "A full body workout with multiple Tabata pairs.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "Tabata",
				BlockNotes: "8 rounds, 20s work, 10s rest for each pair. Rest 30s between pairs.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squats"], Tips: "Tabata 1"},
					{Order: 2, ExerciseID: exIDs["Starjumps"], Tips: "Tabata 1"},
					{Order: 3, ExerciseID: exIDs["Press Ups"], Tips: "Tabata 2"},
					{Order: 4, ExerciseID: exIDs["Mountain Climbers"], Tips: "Tabata 2"},
					{Order: 5, ExerciseID: exIDs["Lunges"], Tips: "Tabata 3"},
					{Order: 6, ExerciseID: exIDs["High Knees"], Tips: "Tabata 3"},
					{Order: 7, ExerciseID: exIDs["Ab Twists"], Tips: "Tabata 4"},
					{Order: 8, ExerciseID: exIDs["Diamond Sit Ups"], Tips: "Tabata 4"},
					{Order: 9, ExerciseID: exIDs["Elbows to Knee"], Tips: "Tabata 5"},
					{Order: 10, ExerciseID: exIDs["Ski Jumps"], Tips: "Tabata 5"},
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
		Description: "As many rounds as possible (AMRAP) in 20 minutes.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "AMRAP",
				BlockNotes: "20 minutes total.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpees"], Reps: "5"},
					{Order: 2, ExerciseID: exIDs["Squat Jumps"], Reps: "10"},
					{Order: 3, ExerciseID: exIDs["Press Ups"], Reps: "15"},
					{Order: 4, ExerciseID: exIDs["Mountain Climbers"], Reps: "20"},
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
				BlockNotes:  "50s work, 10s rest. 60-90s rest between rounds.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Sprints"], Duration: "50s", RestDuration: "10s", Tips: "On the spot"},
					{Order: 2, ExerciseID: exIDs["Sprawls"], Duration: "50s", RestDuration: "10s"},
					{Order: 3, ExerciseID: exIDs["T-Runs"], Duration: "50s", RestDuration: "10s"},
					{Order: 4, ExerciseID: exIDs["Ski Jumps"], Duration: "50s", RestDuration: "10s"},
					{Order: 5, ExerciseID: exIDs["Box Jumps"], Duration: "50s", RestDuration: "10s"},
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
		Description: "An EMOM circuit for full body conditioning.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "EMOM",
				BlockRounds: 5,
				BlockNotes:  "5 rounds total. Perform reps at the top of each minute.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Jumps"], Reps: "12", Tips: "Minute 1"},
					{Order: 2, ExerciseID: exIDs["Press Ups"], Reps: "12", Tips: "Minute 2 (full or on knees)"},
					{Order: 3, ExerciseID: exIDs["Lunges"], Reps: "10 per leg", Tips: "Minute 3"},
					{Order: 4, ExerciseID: exIDs["Mountain Climbers"], Reps: "30 total", Tips: "Minute 4"},
					{Order: 5, ExerciseID: exIDs["Jack Knife"], Reps: "10", Tips: "Minute 5"},
				},
			},
		},
	}
	createWorkoutDay(day27)

	// --- DAY 28 ---
	day28 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   28,
		Title:       "Full Body Pyramid Workout",
		Description: "A pyramid workout with a mix of exercises.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "Pyramid",
				BlockNotes: "Pyramid up and down: 1-10 reps then back down to 1.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Jumps"], Tips: "Or Squats"},
					{Order: 2, ExerciseID: exIDs["Press Ups"], Tips: "Or on Knees"},
					{Order: 3, ExerciseID: exIDs["Plank Jabs"], Reps: "2 Jabs = 1 rep"},
					{Order: 4, ExerciseID: exIDs["Reverse Lunge"]},
					{Order: 5, ExerciseID: exIDs["Diamond Sit Ups"]},
				},
			},
		},
	}
	createWorkoutDay(day28)

	// --- DAY 29 ---
	day29 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   29,
		Title:       "AMRAP & Mobility",
		Description: "An optional mobility day or an AMRAP workout.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "AMRAP",
				BlockRounds: 2,
				BlockNotes:  "Complete as many reps as possible. 2 x 8 min work, 2 min rest in between.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["High Knees"], Reps: "20"},
					{Order: 2, ExerciseID: exIDs["Cross Jabs"], Reps: "20"},
					{Order: 3, ExerciseID: exIDs["Diamond Sit Ups"], Reps: "20"},
					{Order: 4, ExerciseID: exIDs["Belt Kicks"], Reps: "20"},
				},
			},
		},
	}
	createWorkoutDay(day29)

	// --- DAY 30 ---
	day30 := models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   30,
		Title:       "FINAL FITNESS TEST",
		Description: "Repeat the original assessment to measure progress.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "Fitness Assessment",
				BlockNotes: "Complete original assessment and compare results.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Ups"], Reps: "Max Reps", Duration: "1 min", RestDuration: "2 mins"},
					{Order: 2, ExerciseID: exIDs["Squats"], Reps: "Max Reps", Duration: "1 min", RestDuration: "2 mins"},
					{Order: 3, ExerciseID: exIDs["Plank Hold"], Reps: "Max Time", RestDuration: "2 mins"},
					{Order: 4, ExerciseID: exIDs["Burpees"], Reps: "Max Reps", Duration: "1 min", RestDuration: "2 mins"},
					{Order: 5, ExerciseID: exIDs["Starjumps"], Reps: "Max Reps", Duration: "1 min", RestDuration: "2 mins"},
					{Order: 6, ExerciseID: exIDs["Diamond Sit Ups"], Reps: "Max Reps", Duration: "1 min", RestDuration: "2 mins"},
				},
			},
		},
	}
	createWorkoutDay(day30)

	log.Println("Successfully seeded Beginner Workout data.")
}
