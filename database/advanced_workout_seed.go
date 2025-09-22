package database

import (
	"log"

	"github.com/88warren/lmw-fitness-backend/models"
)

func AdvancedWorkoutDaySeed() {
	log.Println("Seeding workout day, block, and exercise data for the 30-Day Advanced Programme...")

	programID, err := getProgramIDByName(DB, "advanced-program")
	if err != nil {
		log.Fatalf("Failed to find '30-Day Advanced Program': %v", err)
	}

	exIDs := make(map[string]uint)
	exercises := []string{
		"Ab Twists", "Bearcrawls", "Belt Kicks", "Bicycle Legs", "Bicycles", "Broad Jumps",
		"Burpee Sprints", "Burpee Tucks", "Burpees", "Burpees (modified)", "Calf Jumps",
		"Cross Jacks", "Crunches", "Diamond Press Ups", "Diamond Sit Ups", "Explosive Starjumps",
		"Flutter Kicks", "Glute Bridges", "H.O.G. Press Ups", "Half Sit Ups", "Heel Flicks",
		"Heel Taps", "High Knees", "High Low Plank", "Hollow Hold", "Inch Worm", "Jack Knife",
		"Jump Lunge", "Knees to Chest", "Leg Circles", "Leg Raises", "Lunges", "Mobility", "Mountain Climbers",
		"Moving Press Ups", "Oblique Hops", "Oblique Press Ups", "Overhead Jabs (Fast)", "Pike Jumps",
		"Plank Hold", "Plank Jabs", "Plank Leg Raises", "Plank Shoulder Taps", "Plyo Press Ups",
		"Press Up Twists", "Press Ups", "Scissors", "Sit Ups", "Ski Jumps", "Sprawls", "Sprints",
		"Squat Hold", "Squat Jumps", "Squat Kicks", "Squat Twists", "Squats", "Starjumps", "Straddle Sit Ups",
		"Switch Kicks", "T-Runs", "Thrusters", "Tricep Dips (Floor)", "Tricep Dips (with Chair)", "Tuck Jumps", "V Press Ups",
		"Wall Sits", "Walkaways", "Wide Arm Press Ups", "Y Shaped Lunges",
	}

	for _, name := range exercises {
		id, err := getExerciseIDByName(DB, name)
		if err != nil {
			log.Printf("Failed to find exercise '%s'. Please seed this exercise first.", name)
		}
		exIDs[name] = id
	}

	createWorkoutDay := func(day models.WorkoutDay) {
		var existingDay models.WorkoutDay
		if err := DB.Where("program_id = ? AND day_number = ?", day.ProgramID, day.DayNumber).First(&existingDay).Error; err == nil {
			log.Printf("Advanced Program - Day %d already exists, skipping creation.", day.DayNumber)
			return
		}
		if err := DB.Create(&day).Error; err != nil {
			log.Printf("Failed to create Advanced Program - Day %d: %v", day.DayNumber, err)
		} else {
			log.Printf("Successfully created Advanced Program - Day %d: %s", day.DayNumber, day.Title)
		}
	}

	// --- DAY 1 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   1,
		Title:       "Fitness Assessment",
		Description: "Complete these 8 exercises for 1 minute each. Make sure you record your results, you will need them for day 30 - there's a table attached in your email to help",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "Fitness Assessment",
				BlockNotes: "Try to do as many reps as possible. Use the whole 2 mins rest after each exercise, to be able to give 100% effort for the next exercise.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Ups"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 2, ExerciseID: exIDs["Straddle Sit Ups"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 3, ExerciseID: exIDs["Plank Hold"], Reps: "Max Effort", Duration: "Max Time", Rest: "2 mins"},
					{Order: 4, ExerciseID: exIDs["Squat Jumps"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 5, ExerciseID: exIDs["Burpees"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 6, ExerciseID: exIDs["Jump Lunge"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins", Tips: "2 Lunges = 1 rep"},
					{Order: 7, ExerciseID: exIDs["Explosive Starjumps"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 8, ExerciseID: exIDs["Thrusters"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
				},
			},
		},
	})

	// --- DAY 2 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   2,
		Title:       "Upper Body Power",
		Description: "Every Minute on the Minute (EMOM) Complete the following reps within a minute. The quicker yu do them, the more rest you get.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "EMOM",
				BlockRounds: 5,
				BlockNotes:  "20 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Wide Arm Press Ups"], Reps: "15"},
					{Order: 2, ExerciseID: exIDs["Tricep Dips (with Chair)"], Reps: "10"},
					{Order: 3, ExerciseID: exIDs["Plank Shoulder Taps"], Reps: "30", Tips: "2 Taps = 1 rep"},
					{Order: 4, ExerciseID: exIDs["Plyo Press Ups"], Reps: "10"},
				},
			},
		},
	})

	// --- DAY 3 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   3,
		Title:       "Lower Body Strength",
		Description: "Complete for time. Reps can be broken up or done in any order.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "For Time",
				BlockRounds: 4,
				BlockNotes:  "Complete 4 rounds of all 5 exercises for the given number of reps. Complete the workout as fast as you can.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Jumps"], Reps: "30"},
					{Order: 2, ExerciseID: exIDs["Broad Jumps"], Reps: "10"},
					{Order: 3, ExerciseID: exIDs["Jump Lunge"], Reps: "30", Tips: "1 Lunge = 1 rep"},
					{Order: 4, ExerciseID: exIDs["Burpees"], Reps: "10"},
					{Order: 5, ExerciseID: exIDs["Glute Bridges"], Reps: "30"},
				},
			},
		},
	})

	// --- DAY 4 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   4,
		Title:       "Core Circuit",
		Description: "A circuit focused on core strength.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 3,
				RoundRest:   "60s",
				BlockNotes:  "Exercise for 45 seconds and then rest for 15 seconds. Repeat the circuit 3 times. Rest 60 seconds between rounds. Full duration 24 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Plank Shoulder Taps"], Duration: "45s", Rest: "15s"},
					{Order: 2, ExerciseID: exIDs["Bicycle Legs"], Duration: "45s", Rest: "15s"},
					{Order: 3, ExerciseID: exIDs["V Press Ups"], Duration: "45s", Rest: "15s"},
					{Order: 4, ExerciseID: exIDs["Ab Twists"], Duration: "45s", Rest: "15s"},
					{Order: 5, ExerciseID: exIDs["Flutter Kicks"], Duration: "45s", Rest: "15s"},
					{Order: 6, ExerciseID: exIDs["Mountain Climbers"], Duration: "45s", Rest: "15s"},
					{Order: 7, ExerciseID: exIDs["Diamond Sit Ups"], Duration: "45s", Rest: "15s"},
				},
			},
		},
	})

	// --- DAY 5 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   5,
		Title:       "Cardio Intervals",
		Description: "A cardio-focused interval circuit.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 3,
				BlockNotes:  "Exercise for 40 seconds and then rest for 20 seconds. Repeat the circuit 3 times. Full duration 24 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["High Knees"], Duration: "40s", Rest: "20s"},
					{Order: 2, ExerciseID: exIDs["Burpees"], Duration: "40s", Rest: "20s"},
					{Order: 3, ExerciseID: exIDs["Explosive Starjumps"], Duration: "40s", Rest: "20s"},
					{Order: 4, ExerciseID: exIDs["Mountain Climbers"], Duration: "40s", Rest: "20s"},
					{Order: 5, ExerciseID: exIDs["Sprints"], Duration: "40s", Rest: "20s"},
					{Order: 6, ExerciseID: exIDs["Thrusters"], Duration: "40s", Rest: "20s"},
					{Order: 7, ExerciseID: exIDs["Belt Kicks"], Duration: "40s", Rest: "20s"},
					{Order: 8, ExerciseID: exIDs["Heel Flicks"], Duration: "40s", Rest: "20s"},
				},
			},
		},
	})

	// --- DAY 6 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   6,
		Title:       "Full Body & Core Tabata ",
		Description: "4 x Tabata blocks.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Tabata",
				BlockRounds: 8,
				BlockNotes:  "20s work / 10s rest x 8 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpees"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Plank Jabs"], Duration: "20s", Rest: "10s"},
				},
			},
			{
				BlockType:   "Tabata",
				BlockRounds: 8,
				BlockNotes:  "20s work / 10s rest x 8 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["High Knees"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Scissors"], Duration: "20s", Rest: "10s"},
				},
			},
			{
				BlockType:   "Tabata",
				BlockRounds: 8,
				BlockNotes:  "20s work / 10s rest x 8 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Jump Lunge"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Ab Twists"], Duration: "20s", Rest: "10s"},
				},
			},
			{
				BlockType:   "Tabata",
				BlockRounds: 8,
				BlockNotes:  "20s work / 10s rest x 8 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Jumps"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Straddle Sit Ups"], Duration: "20s", Rest: "10s"},
				},
			},
			{
				BlockType:   "Tabata",
				BlockRounds: 8,
				BlockNotes:  "20s work / 10s rest x 8 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Sprints"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Sit Ups"], Duration: "20s", Rest: "10s"},
				},
			},
		},
	})

	// --- DAY 7 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   7,
		Title:       "Full Body Flow",
		Description: "As many rounds as possible (AMRAP) in 25 minutes. Track your rounds with the counter!",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "AMRAP",
				BlockNotes: "25 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Thrusters"], Reps: "10"},
					{Order: 2, ExerciseID: exIDs["Press Ups"], Reps: "15"},
					{Order: 3, ExerciseID: exIDs["Squat Twists"], Reps: "20"},
					{Order: 4, ExerciseID: exIDs["Knees to Chest"], Reps: "25"},
					{Order: 5, ExerciseID: exIDs["Heel Taps"], Reps: "30"},
				},
			},
		},
	})

	// --- DAY 8 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   8,
		Title:       "Recovery day & Optional workout",
		Description: "Mobility. Optional Workout: Every Minute on the Minute (EMOM) Complete the following reps within a minute. The quicker yu do them, the more rest you get",
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
				BlockType:   "EMOM",
				BlockRounds: 4,
				BlockNotes:  "20 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Jumps"], Reps: "20"},
					{Order: 2, ExerciseID: exIDs["Jump Lunge"], Reps: "10", Tips: "2 Lunges = 1 rep"},
					{Order: 3, ExerciseID: exIDs["Burpees"], Reps: "10"},
					{Order: 4, ExerciseID: exIDs["Crunches"], Reps: "20"},
					{Order: 5, ExerciseID: exIDs["Tricep Dips (with Chair)"], Reps: "15"},
				},
			},
		},
	})

	// --- DAY 9 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   9,
		Title:       "Plyometric Power",
		Description: "A pyramid-style workout focused on explosive movements.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "For Time",
				BlockNotes: "A pyramid workout for plyometric power. Complete all exercises from 2 reps to 16 reps and back to 2 reps as quickly as possible.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Tuck Jumps"], Reps: "2, 4, 6, 8, 10, 12, 14, 16, 14, 12, 10, 8, 6, 4, 2"},
					{Order: 2, ExerciseID: exIDs["Plyo Press Ups"], Reps: "2, 4, 6, 8, 10, 12, 14, 16, 14, 12, 10, 8, 6, 4, 2"},
					{Order: 3, ExerciseID: exIDs["Squat Jumps"], Reps: "2, 4, 6, 8, 10, 12, 14, 16, 14, 12, 10, 8, 6, 4, 2"},
					{Order: 4, ExerciseID: exIDs["Explosive Starjumps"], Reps: "2, 4, 6, 8, 10, 12, 14, 16, 14, 12, 10, 8, 6, 4, 2"},
				},
			},
		},
	})

	// --- DAY 10 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   10,
		Title:       "Upper Body Endurance",
		Description: "Complete for time. Reps can be broken up or done in any order.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "For Time",
				BlockNotes: "Complete all 4 exercise for the given number of reps. Work through the exercises as fast as possible.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Ups"], Reps: "100"},
					{Order: 2, ExerciseID: exIDs["Tricep Dips (with Chair)"], Reps: "75"},
					{Order: 3, ExerciseID: exIDs["High Low Plank"], Reps: "50"},
					{Order: 4, ExerciseID: exIDs["Walkaways"], Reps: "25"},
				},
			},
		},
	})

	// --- DAY 11 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   11,
		Title:       "Core Focus & Full Body Burst",
		Description: "Workout 1: Every Minute on the Minute (EMOM) Complete the following reps within a minute. The quicker you do them, the more rest you get. Workout 2: For time workout",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "EMOM",
				BlockRounds: 4,
				BlockNotes:  "12 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Flutter Kicks"], Reps: "30", Tips: "2 Kicks = 1 rep"},
					{Order: 2, ExerciseID: exIDs["Leg Raises"], Reps: "20"},
					{Order: 3, ExerciseID: exIDs["Jack Knife"], Reps: "10"},
				},
			},
			{
				BlockType:   "For Time",
				BlockRounds: 4,
				BlockNotes:  "Complete 4 rounds of all 3 exercises for the given number of reps. Complete the workout as fast as you can.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpee Sprints"], Reps: "12"},
					{Order: 2, ExerciseID: exIDs["Cross Jacks"], Reps: "12", Tips: "2 Jacks = 1 rep"},
					{Order: 3, ExerciseID: exIDs["Explosive Starjumps"], Reps: "12"},
				},
			},
		},
	})

	// --- DAY 12 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   12,
		Title:       "Core Domination",
		Description: "Every Minute on the Minute (EMOM) Complete the following reps within a minute. The quicker yu do them, the more rest you get",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "EMOM",
				BlockRounds: 7,
				BlockNotes:  "21 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Crunches"], Reps: "20"},
					{Order: 2, ExerciseID: exIDs["Leg Raises"], Reps: "15"},
					{Order: 3, ExerciseID: exIDs["Plank Hold"], Duration: "40s"},
				},
			},
		},
	})

	// --- DAY 13 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   13,
		Title:       "Metabolic Mayhem",
		Description: "A circuit designed for metabolic conditioning.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 5,
				BlockNotes:  "Exercise for 40 seconds and then rest for 20 seconds. Repeat the circuits 5 times. Full duration 30 minutes.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpee Sprints"], Duration: "40s", Rest: "20s"},
					{Order: 2, ExerciseID: exIDs["Switch Kicks"], Duration: "40s", Rest: "20s"},
					{Order: 3, ExerciseID: exIDs["Bearcrawls"], Duration: "40s", Rest: "20s"},
					{Order: 4, ExerciseID: exIDs["Sprawls"], Duration: "40s", Rest: "20s"},
					{Order: 5, ExerciseID: exIDs["T-Runs"], Duration: "40s", Rest: "20s"},
					{Order: 6, ExerciseID: exIDs["Ski Jumps"], Duration: "40s", Rest: "20s"},
				},
			},
		},
	})

	// --- DAY 14 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   14,
		Title:       "Full Body AMRAP",
		Description: "As many rounds as possible (AMRAP) in 25 minutes. Track your rounds with the counter!",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "AMRAP",
				BlockRounds: 3,
				BlockNotes:  "25 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Jumps"], Reps: "25"},
					{Order: 2, ExerciseID: exIDs["Wide Arm Press Ups"], Reps: "20"},
					{Order: 3, ExerciseID: exIDs["Mountain Climbers"], Reps: "15", Tips: "2 Climbers = 1 rep"},
					{Order: 4, ExerciseID: exIDs["Glute Bridges"], Reps: "10"},
					{Order: 5, ExerciseID: exIDs["Burpees"], Reps: "5"},
				},
			},
		},
	})

	// --- DAY 15 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   15,
		Title:       "Recovery day & Optional workout",
		Description: "Mobility. Optional 2 x mini workouts: AMRAP & for time pyramid.",
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
				BlockType:  "AMRAP",
				BlockNotes: "12 minutes.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Jumps"], Reps: "10"},
					{Order: 2, ExerciseID: exIDs["Plyo Press Ups"], Reps: "8"},
					{Order: 3, ExerciseID: exIDs["Burpee Tucks"], Reps: "6"},
				},
			},
			{
				BlockType:  "For Time",
				BlockNotes: "A pyramid workout for full body endurance. Complete all exercises from 5 reps to 10 reps and back to 5 reps as quickly as possible.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Inch Worm"], Reps: "5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5"},
					{Order: 2, ExerciseID: exIDs["Pike Jumps"], Reps: "5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5"},
				},
			},
		},
	})

	// --- DAY 16 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   16,
		Title:       "Lower Body Power & Endurance",
		Description: "Workout 1: A lower body circuit. Workout 2: A descending ladder for time.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 3,
				BlockNotes:  "Exercise for 50 seconds and then rest for 10 seconds. Repeat the circuit 3 times. Full duration 9 minutes.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Jump Lunge"], Duration: "50s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Squat Kicks"], Duration: "50s", Rest: "10s"},
					{Order: 3, ExerciseID: exIDs["Ski Jumps"], Duration: "50s", Rest: "10s"},
				},
			},
			{
				BlockType:  "For Time",
				BlockNotes: "A descending ladder workout for lower body endurance. Complete all exercises from 10 reps down to 2 as quickly as possible.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Broad Jumps"], Reps: "10, 8, 6, 4, 2"},
					{Order: 2, ExerciseID: exIDs["Thrusters"], Reps: "10, 8, 6, 4, 2"},
				},
			},
		},
	})

	// --- DAY 17 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   17,
		Title:       "Upper Body Strength & Endurance",
		Description: "Every Minute on the Minute (EMOM) Complete the following reps within a minute. The quicker yu do them, the more rest you get",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "EMOM",
				BlockRounds: 4,
				BlockNotes:  "24 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Diamond Sit Ups"], Reps: "20"},
					{Order: 2, ExerciseID: exIDs["Press Ups"], Reps: "20"},
					{Order: 3, ExerciseID: exIDs["Plank Hold"], Duration: "40s"},
					{Order: 4, ExerciseID: exIDs["Tricep Dips (with Chair)"], Reps: "20"},
					{Order: 5, ExerciseID: exIDs["Straddle Sit Ups"], Reps: "20"},
					{Order: 6, ExerciseID: exIDs["Overhead Jabs (Fast)"], Reps: "20", Tips: "2 Jabs = 1 rep"},
				},
			},
		},
	})

	// --- DAY 18 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   18,
		Title:       "Recovery day & Optional workout",
		Description: "Mobility AND/OR Workout: AMRAP in 20 minutes.",
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
				BlockType:  "AMRAP",
				BlockNotes: "20 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpee Tucks"], Reps: "5"},
					{Order: 2, ExerciseID: exIDs["Jump Lunge"], Reps: "10", Tips: "2 Lunges = 1 rep"},
					{Order: 3, ExerciseID: exIDs["Mountain Climbers"], Reps: "15", Tips: "2 Climbers = 1 rep"},
					{Order: 4, ExerciseID: exIDs["High Knees"], Reps: "20", Tips: "2 High Knees = 1 rep"},
					{Order: 5, ExerciseID: exIDs["Heel Flicks"], Reps: "25", Tips: "2 Flicks = 1 rep"},
				},
			},
		},
	})

	// --- DAY 19 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   19,
		Title:       "Core & Cardio Challenge",
		Description: "5 x Tabata rounds.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Tabata",
				BlockRounds: 6,
				BlockNotes:  "20s work / 10s rest x 6 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Sprints"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Knees to Chest"], Duration: "20s", Rest: "10s"},
				},
			},
			{
				BlockType:   "Tabata",
				BlockRounds: 6,
				BlockNotes:  "20s work / 10s rest x 6 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Explosive Starjumps"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Bicycle Legs"], Duration: "20s", Rest: "10s"},
				},
			},
			{
				BlockType:   "Tabata",
				BlockRounds: 6,
				BlockNotes:  "20s work / 10s rest x 6 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["High Knees"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Diamond Sit Ups"], Duration: "20s", Rest: "10s"},
				},
			},
			{
				BlockType:   "Tabata",
				BlockRounds: 6,
				BlockNotes:  "20s work / 10s rest x 6 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Thrusters"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Ab Twists"], Duration: "20s", Rest: "10s"},
				},
			},
			{
				BlockType:   "Tabata",
				BlockRounds: 6,
				BlockNotes:  "20s work / 10s rest x 6 rounds",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Belt Kicks"], Duration: "20s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Leg Raises"], Duration: "20s", Rest: "10s"},
				},
			},
		},
	})

	// --- DAY 20 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   20,
		Title:       "Full Body Fusion",
		Description: "Complex training circuit for time.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "For Time",
				BlockRounds: 5,
				BlockNotes:  "Complete 5 rounds of all 5 exercises for the given number of reps. Complete the workout as fast as you can.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Thrusters"], Reps: "8"},
					{Order: 2, ExerciseID: exIDs["Burpees"], Reps: "10"},
					{Order: 3, ExerciseID: exIDs["Squat Twists"], Reps: "12"},
					{Order: 4, ExerciseID: exIDs["Press Up Twists"], Reps: "14"},
					{Order: 5, ExerciseID: exIDs["Explosive Starjumps"], Reps: "16"},
				},
			},
		},
	})

	// --- DAY 21 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   21,
		Title:       "Endurance Test",
		Description: "Complete for time. Reps can be broken up or done in any order.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "For Time",
				BlockNotes: "Complete all 4 exercises for the given number of reps. Work through the exercises as fast as possible.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Jumps"], Reps: "200"},
					{Order: 2, ExerciseID: exIDs["Press Ups"], Reps: "150"},
					{Order: 3, ExerciseID: exIDs["Burpees"], Reps: "100"},
					{Order: 4, ExerciseID: exIDs["Tuck Jumps"], Reps: "50"},
				},
			},
		},
	})

	// --- DAY 22 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   22,
		Title:       "Recovery day & Optional workout",
		Description: "Mobility. Optional 2 x mini workout: for time pyramid & AMRAP Finisher ",
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
				BlockType:  "For Time",
				BlockNotes: "A pyramid workout for full body endurance. Complete all exercises from 10 reps to 20 reps and back to 10 reps as quickly as possible.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpees"], Reps: "10, 15, 20, 15, 10"},
					{Order: 2, ExerciseID: exIDs["Squat Jumps"], Reps: "10, 15, 20, 15, 10"},
				},
			},
			{
				BlockType:  "AMRAP",
				BlockNotes: "5 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Tuck Jumps"], Reps: "5"},
					{Order: 2, ExerciseID: exIDs["Wide Arm Press Ups"], Reps: "10"},
					{Order: 3, ExerciseID: exIDs["Crunches"], Reps: "15"},
				},
			},
		},
	})

	// --- DAY 23 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   23,
		Title:       "Upper Body & Core Endurance",
		Description: "Workout 1: Timed circuit. Workout 2: For time pyramid. Workout 3: Burpee finisher.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 3,
				BlockNotes:  "Exercise for 50 seconds and then rest for 10 seconds. Repeat the circuit 3 times. Full duration 12 mintues.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Up Twists"], Duration: "50s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Oblique Press Ups"], Duration: "50s", Rest: "10s"},
					{Order: 3, ExerciseID: exIDs["Plank Leg Raises"], Duration: "50s", Rest: "10s"},
					{Order: 4, ExerciseID: exIDs["Plank Hold"], Duration: "50s", Rest: "10s"},
				},
			},
			{
				BlockType:  "For Time",
				BlockNotes: "Complete the following two exercises in the Pyramid workout. Working from 25 reps down to 15, in groups of 5 and back to 15 reps.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["V Press Ups"], Reps: "25, 20, 15, 20, 25"},
					{Order: 2, ExerciseID: exIDs["Sit Ups"], Reps: "25, 20, 15, 20, 25"},
				},
			},
			{
				BlockType:  "For Time",
				BlockNotes: "Finisher: Complete 50 burpees as quickly as possible.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpees"], Reps: "50"},
				},
			},
		},
	})

	// --- DAY 24 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   24,
		Title:       "Lower Body Endurance & Agility",
		Description: "Workout 1: Timed circuit. Workout 2: AMRAP. Workout 3: Static hold finisher.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 4,
				BlockNotes:  "Exercise for 50 seconds and rest for 10 seconds. Repeat the circuit 4 times. Full duration 16 minutes.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["T-Runs"], Duration: "50s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Y Shaped Lunges"], Duration: "50s", Rest: "10s"},
					{Order: 3, ExerciseID: exIDs["Squat Twists"], Duration: "50s", Rest: "10s"},
					{Order: 4, ExerciseID: exIDs["Calf Jumps"], Duration: "50s", Rest: "10s"},
				},
			},
			{
				BlockType:  "AMRAP",
				BlockNotes: "10 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Switch Kicks"], Reps: "10", Tips: "2 Kicks = 1 rep"},
					{Order: 2, ExerciseID: exIDs["Thrusters"], Reps: "10"},
					{Order: 3, ExerciseID: exIDs["Broad Jumps"], Reps: "5"},
				},
			},
			{
				BlockType:  "EMOM",
				BlockNotes: "4 Minutes: Finisher: Multiple Static Holds: 1 minute each.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Hold"], Tips: "Hold", Duration: "1 min"},
					{Order: 2, ExerciseID: exIDs["Wall Sits"], Tips: "Hold", Duration: "1 min"},
					{Order: 3, ExerciseID: exIDs["Hollow Hold"], Tips: "Hold", Duration: "1 min"},
					{Order: 4, ExerciseID: exIDs["Plank Hold"], Tips: "Hold", Duration: "1 min"},
				},
			},
		},
	})

	// --- DAY 25 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   25,
		Title:       "Plyo Push",
		Description: "Workout 1: Plyometric for time workout. Workout 2: Plank challenge finisher.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "For Time",
				BlockRounds: 4,
				BlockNotes:  "Complete 4 rounds of all 8 exercises for the given number of reps. Complete the workout as fast as you can.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Plyo Press Ups"], Reps: "12"},
					{Order: 2, ExerciseID: exIDs["Tuck Jumps"], Reps: "20"},
					{Order: 3, ExerciseID: exIDs["H.O.G. Press Ups"], Reps: "12"},
					{Order: 4, ExerciseID: exIDs["Explosive Starjumps"], Reps: "20"},
					{Order: 5, ExerciseID: exIDs["Moving Press Ups"], Reps: "12"},
					{Order: 6, ExerciseID: exIDs["Ski Jumps"], Reps: "20"},
					{Order: 7, ExerciseID: exIDs["Oblique Hops"], Reps: "12"},
					{Order: 8, ExerciseID: exIDs["Jump Lunge"], Reps: "20"},
				},
			},
			{
				BlockType:  "For Time",
				BlockNotes: "Finisher: Plank Challenge: Hold for as long as possible or accumulate 5-minute total hold.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Plank Hold"], Duration: "Max time"},
				},
			},
		},
	})

	// --- DAY 26 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   26,
		Title:       "Core Strength & Upper Body Finisher",
		Description: "Workout 1: Timed circuit. Workout 2: EMOM press up finisher. Choose any variation of press up you want for the given number of reps.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "For Time",
				BlockRounds: 4,
				BlockNotes:  "Complete 4 rounds of all 8 exercises for the given number of reps. Complete the workout as fast as you can.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Diamond Sit Ups"], Reps: "15"},
					{Order: 2, ExerciseID: exIDs["High Low Plank"], Reps: "15"},
					{Order: 3, ExerciseID: exIDs["Press Ups"], Reps: "15"},
					{Order: 4, ExerciseID: exIDs["Tricep Dips (with Chair)"], Reps: "15"},
					{Order: 5, ExerciseID: exIDs["Bicycles"], Reps: "15", Tips: "2 Bicycles = 1 rep"},
					{Order: 6, ExerciseID: exIDs["Sit Ups"], Reps: "15"},
					{Order: 7, ExerciseID: exIDs["H.O.G. Press Ups"], Reps: "15"},
					{Order: 8, ExerciseID: exIDs["Plank Jabs"], Reps: "15", Tips: "2 Jabs = 1 rep"},
				},
			},
			{
				BlockType:   "EMOM",
				BlockRounds: 1,
				BlockNotes:  "5 minutes: EMOM Finisher Press Up Variations (you chosoe your variation). 12 reps in minute 1, increase by 2 reps each minute.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Ups"], Reps: "12"},
					{Order: 2, ExerciseID: exIDs["Press Ups"], Reps: "14"},
					{Order: 3, ExerciseID: exIDs["Press Ups"], Reps: "16"},
					{Order: 4, ExerciseID: exIDs["Press Ups"], Reps: "18"},
					{Order: 5, ExerciseID: exIDs["Press Ups"], Reps: "20"},
				},
			},
		},
	})

	// --- DAY 27 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   27,
		Title:       "Full Body Cardio and Agility",
		Description: "Workout 1: Full Body Cardio circuit. Workout 2: EMOM burpee finisher. Complete 6 burpees in the first minute and increase the reps by 2 each minute.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 5,
				BlockNotes:  "Exercise for 50 seconds and then rest for 10 seconds. Repeat the circuit 5 times. Full duration 30 minutes.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpee Sprints"], Duration: "50s", Rest: "10s"},
					{Order: 2, ExerciseID: exIDs["Thrusters"], Duration: "50s", Rest: "10s"},
					{Order: 3, ExerciseID: exIDs["Squat Jumps"], Duration: "50s", Rest: "10s"},
					{Order: 4, ExerciseID: exIDs["T-Runs"], Duration: "50s", Rest: "10s"},
					{Order: 5, ExerciseID: exIDs["High Knees"], Duration: "50s", Rest: "10s"},
					{Order: 6, ExerciseID: exIDs["Mountain Climbers"], Duration: "50s", Rest: "10s"},
				},
			},
			{
				BlockType:  "EMOM",
				BlockNotes: "5 minutes: EMOM Finisher Death by Burpees. 6 reps in minute 1, increase by 2 reps each minute.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpees"], Reps: "6"},
					{Order: 2, ExerciseID: exIDs["Burpees"], Reps: "8"},
					{Order: 3, ExerciseID: exIDs["Burpees"], Reps: "10"},
					{Order: 4, ExerciseID: exIDs["Burpees"], Reps: "12"},
					{Order: 5, ExerciseID: exIDs["Burpees"], Reps: "14"},
				},
			},
		},
	})

	// --- DAY 28 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   28,
		Title:       "Endurance Workout",
		Description: "Workout 1: 5 min squat jumps. Workout 2: 4 min press ups. Workouut 3: 3 min mountain climbers. Workout 4: 2 min burpees. Workout 5: 1 min tuck jumps. Workout 6: full body ladder finisher.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "AMRAP",
				BlockNotes: "5 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Jumps"], Reps: "Max"},
				},
			},
			{
				BlockType:  "AMRAP",
				BlockNotes: "4 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 2, ExerciseID: exIDs["Press Ups"], Reps: "Max"},
				},
			},
			{
				BlockType:  "AMRAP",
				BlockNotes: "3 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 3, ExerciseID: exIDs["Mountain Climbers"], Reps: "Max"},
				},
			},
			{
				BlockType:  "AMRAP",
				BlockNotes: "2 minutes",
				Exercises: []models.WorkoutExercise{
					{Order: 4, ExerciseID: exIDs["Burpees"], Reps: "Max"},
				},
			},
			{
				BlockType:  "AMRAP",
				BlockNotes: "1 minute",
				Exercises: []models.WorkoutExercise{
					{Order: 5, ExerciseID: exIDs["Tuck Jumps"], Reps: "Max"},
				},
			},
			{
				BlockType:  "For Time",
				BlockNotes: "A descending ladder workout for full body endurance. Complete all exercises from 10 reps down to 1 as quickly as possible.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Leg Circles"], Reps: "10, 9, 8, 7, 6, 5, 4, 3, 2, 1"},
					{Order: 2, ExerciseID: exIDs["Burpee Sprints"], Reps: "10, 9, 8, 7, 6, 5, 4, 3, 2, 1"},
				},
			},
		},
	})

	// --- DAY 29 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   29,
		Title:       "Recovery day & Optional workout",
		Description: "Mobility. Optional Workout: A final finisher. Complete for time. Reps can be broken up or done in any order.",
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
				BlockType:  "For Time",
				BlockNotes: "Complete all 5 exercises for the given number of reps as quickly as possible.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Jumps"], Reps: "200"},
					{Order: 2, ExerciseID: exIDs["Mountain Climbers"], Reps: "150"},
					{Order: 3, ExerciseID: exIDs["Burpees"], Reps: "100"},
					{Order: 4, ExerciseID: exIDs["Sit Ups"], Reps: "75"},
					{Order: 5, ExerciseID: exIDs["Tuck Jumps"], Reps: "50"},
				},
			},
		},
	})

	// --- DAY 30 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   30,
		Title:       "FINAL FITNESS ASSESSMENT",
		Description: "Complete this fitness assessment one more time and compare the results from Day 1.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "Fitness Assessment",
				BlockNotes: "Push yourself as hard as you did on day 1 and note your improvements.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Ups"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 2, ExerciseID: exIDs["Straddle Sit Ups"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 3, ExerciseID: exIDs["Plank Hold"], Reps: "Max Effort", Duration: "Max Time", Rest: "2 mins"},
					{Order: 4, ExerciseID: exIDs["Squat Jumps"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 5, ExerciseID: exIDs["Burpees"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 6, ExerciseID: exIDs["Jump Lunge"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins", Tips: "2 Lunges = 1 rep"},
					{Order: 7, ExerciseID: exIDs["Explosive Starjumps"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
					{Order: 8, ExerciseID: exIDs["Thrusters"], Reps: "Max Effort", Duration: "1 min", Rest: "2 mins"},
				},
			},
		},
	})

	log.Println("Successfully seeded data for the entire 30-Day Advanced Programme.")
}
