package database

import (
	"log"

	"github.com/88warren/lmw-fitness-backend/models"
)

func AdvancedWorkoutDaySeed() {
	log.Println("Seeding workout day, block, and exercise data for the 30-Day Advanced Programme...")

	programID, err := getProgramIDByName(DB, "30-Day Advanced Program")
	if err != nil {
		log.Fatalf("Failed to find '30-Day Advanced Program': %v", err)
	}

	// Fetch all necessary exercise IDs upfront for efficiency
	exIDs := make(map[string]uint)
	exercises := []string{
		"Press Ups", "Straddle Sit-ups", "Plank Hold", "Squats", "Burpees", "Jump Lunges",
		"Pike Jumps", "Plank Shoulder Taps", "Plyo Press-ups", "Squat Jumps", "Lunges",
		"Broad Jumps", "Glute Bridges", "Bicycle Legs", "V Press-ups", "Ab Twists",
		"Flutter Kicks", "Mountain Climbers", "Diamond Sit-ups", "High Knees", "Starjumps",
		"Sprints", "Thrusters", "Belt Kicks", "Heel Flicks", "Plank Jabs", "Scissor Kicks",
		"Squat Twists", "Knee to Chest", "Heel Taps", "Tuck Jumps", "High Low Plank",
		"Walkaways", "Jack Knife", "Burpee Sprints", "Cross Jacks", "Leg Raises", "Bear Crawls",
		"Switch Kicks", "T-Runs", "Ski Jumps", "Wide Arm Press-ups", "Inchworms", "Burpee Tuck Jumps",
		"Oblique Press-ups", "Half Sit Ups", "Y-shaped Lunges", "Calf Jumps", "HOG Press-ups",
		"Moving Press-ups", "Oblique Hops", "Burpees", "Hollow Rock Hold",
	}

	for _, name := range exercises {
		id, err := getExerciseIDByName(DB, name)
		if err != nil {
			log.Printf("Failed to find exercise '%s'. Please seed this exercise first.", name)
		}
		exIDs[name] = id
	}

	// Helper function to create a workout day and handle errors
	createWorkoutDay := func(day models.WorkoutDay) {
		var existingDay models.WorkoutDay
		if err := DB.Where("program_id = ? AND day_number = ?", day.ProgramID, day.DayNumber).First(&existingDay).Error; err == nil {
			log.Printf("Day %d already exists, skipping creation.", day.DayNumber)
			return
		}
		if err := DB.Create(&day).Error; err != nil {
			log.Printf("Failed to create Day %d: %v", day.DayNumber, err)
		} else {
			log.Printf("Successfully created Day %d: %s", day.DayNumber, day.Title)
		}
	}

	// --- DAY 1 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   1,
		Title:       "FITNESS ASSESSMENT",
		Description: "Complete for time and record results to compare against Day 30.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "Fitness Assessment",
				BlockNotes: "Complete each exercise for max reps/time. Record results.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Ups"], Reps: "Max Reps", Duration: "1 min"},
					{Order: 2, ExerciseID: exIDs["Straddle Sit-ups"], Reps: "Max Reps", Duration: "1 min"},
					{Order: 3, ExerciseID: exIDs["Plank Hold"], Reps: "Max Time"},
					{Order: 4, ExerciseID: exIDs["Squats"], Reps: "Max Reps", Duration: "1 min"},
					{Order: 5, ExerciseID: exIDs["Burpees"], Reps: "Max Reps", Duration: "3 min"},
					{Order: 6, ExerciseID: exIDs["Jump Lunges"], Reps: "Max Reps", Duration: "1 min"},
				},
			},
		},
	})

	// --- DAY 2 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   2,
		Title:       "Upper Body Power",
		Description: "Every Minute on the Minute (EMOM) for 20 minutes.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "EMOM",
				BlockRounds: 4, // 5 minutes x 4 rounds
				BlockNotes:  "Repeat for 4 rounds. Minute 5 is rest.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Wide Arm Press-ups"], Reps: "8", Instructions: "Minute 1"},
					{Order: 2, ExerciseID: exIDs["Pike Jumps"], Reps: "10", Instructions: "Minute 2"},
					{Order: 3, ExerciseID: exIDs["Plank Shoulder Taps"], Reps: "12", Instructions: "Minute 3"},
					{Order: 4, ExerciseID: exIDs["Plyo Press-ups"], Reps: "8", Instructions: "Minute 4"},
					{Order: 5, Reps: "Rest", Instructions: "Minute 5"},
				},
			},
		},
	})

	// --- DAY 3 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   3,
		Title:       "Lower Body Strength",
		Description: "A circuit focusing on lower body strength with sets and reps.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 4,
				BlockNotes:  "Rest 90s between rounds. Add 1 rep to each exercise/5s to the glute bridge each round.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Jumps"], Reps: "15"},
					{Order: 2, ExerciseID: exIDs["Lunges"], Reps: "20 (10 each leg)"},
					{Order: 3, ExerciseID: exIDs["Broad Jumps"], Reps: "12"},
					{Order: 4, ExerciseID: exIDs["Burpees"], Reps: "10"},
					{Order: 5, ExerciseID: exIDs["Glute Bridges"], Duration: "30s"},
				},
			},
		},
	})

	// --- DAY 4 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   4,
		Title:       "Core Circuit",
		Description: "A timed circuit focused on core strength.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 3,
				BlockNotes:  "45s work, 15s rest. Rest 60s between rounds.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Plank Shoulder Taps"], Duration: "45s", RestDuration: "15s"},
					{Order: 2, ExerciseID: exIDs["Bicycle Legs"], Duration: "45s", RestDuration: "15s"},
					{Order: 3, ExerciseID: exIDs["V Press-ups"], Duration: "45s", RestDuration: "15s"},
					{Order: 4, ExerciseID: exIDs["Ab Twists"], Duration: "45s", RestDuration: "15s"},
					{Order: 5, ExerciseID: exIDs["Flutter Kicks"], Duration: "45s", RestDuration: "15s"},
					{Order: 6, ExerciseID: exIDs["Mountain Climbers"], Duration: "45s", RestDuration: "15s"},
					{Order: 7, ExerciseID: exIDs["Diamond Sit-ups"], Duration: "45s", RestDuration: "15s"},
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
				BlockType:   "Intervals",
				BlockRounds: 3,
				BlockNotes:  "40s work, 20s rest. 3 rounds total.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["High Knees"], Duration: "40s", RestDuration: "20s"},
					{Order: 2, ExerciseID: exIDs["Burpees"], Duration: "40s", RestDuration: "20s"},
					{Order: 3, ExerciseID: exIDs["Starjumps"], Duration: "40s", RestDuration: "20s"},
					{Order: 4, ExerciseID: exIDs["Mountain Climbers"], Duration: "40s", RestDuration: "20s"},
					{Order: 5, ExerciseID: exIDs["Sprints"], Duration: "40s", RestDuration: "20s"},
					{Order: 6, ExerciseID: exIDs["Thrusters"], Duration: "40s", RestDuration: "20s"},
					{Order: 7, ExerciseID: exIDs["Belt Kicks"], Duration: "40s", RestDuration: "20s"},
					{Order: 8, ExerciseID: exIDs["Heel Flicks"], Duration: "40s", RestDuration: "20s"},
				},
			},
		},
	})

	// --- DAY 6 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   6,
		Title:       "Full Body Tabata & Core",
		Description: "Multiple Tabata blocks with rest in between.",
		WorkoutBlocks: []models.WorkoutBlock{
			{BlockType: "Tabata", BlockNotes: "20s work / 10s rest x 8 rounds", Exercises: []models.WorkoutExercise{{Order: 1, ExerciseID: exIDs["Burpees"]}}},
			{BlockType: "Rest", BlockNotes: "Rest 60 seconds"},
			{BlockType: "Tabata", BlockNotes: "20s work / 10s rest x 8 rounds", Exercises: []models.WorkoutExercise{{Order: 1, ExerciseID: exIDs["Plank Jabs"]}}},
			{BlockType: "Rest", BlockNotes: "Rest 60 seconds"},
			{BlockType: "Tabata", BlockNotes: "20s work / 10s rest x 8 rounds", Exercises: []models.WorkoutExercise{{Order: 1, ExerciseID: exIDs["High Knees"]}}},
			{BlockType: "Rest", BlockNotes: "Rest 60 seconds"},
			{BlockType: "Tabata", BlockNotes: "20s work / 10s rest x 8 rounds", Exercises: []models.WorkoutExercise{{Order: 1, ExerciseID: exIDs["Scissor Kicks"]}}},
			{BlockType: "Rest", BlockNotes: "Rest 60 seconds"},
			{BlockType: "Tabata", BlockNotes: "20s work / 10s rest x 8 rounds", Exercises: []models.WorkoutExercise{{Order: 1, ExerciseID: exIDs["Lunges"]}}},
			{BlockType: "Rest", BlockNotes: "Rest 60 seconds"},
			{BlockType: "Tabata", BlockNotes: "20s work / 10s rest x 8 rounds", Exercises: []models.WorkoutExercise{{Order: 1, ExerciseID: exIDs["Squat Jumps"]}}},
		},
	})

	// --- DAY 7 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   7,
		Title:       "Full Body Flow",
		Description: "As Many Rounds as Possible (AMRAP) in 25 minutes.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "AMRAP",
				BlockNotes: "25 minutes total. Complete as many rounds as possible.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Thrusters"], Reps: "10"},
					{Order: 2, ExerciseID: exIDs["Press Ups"], Reps: "15"},
					{Order: 3, ExerciseID: exIDs["Squat Twists"], Reps: "20"},
					{Order: 4, ExerciseID: exIDs["Knee to Chest"], Reps: "25"},
					{Order: 5, ExerciseID: exIDs["Heel Taps"], Reps: "30"},
				},
			},
		},
	})

	// --- DAY 8 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   8,
		Title:       "OPTIONAL MOBILITY DAY",
		Description: "Mobility work or EMOM workout.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "EMOM",
				BlockNotes: "15-20 minutes total. Perform reps at the top of each minute.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, Instructions: "Minute 1", Reps: "10 Squats + 5 Press-ups"},
					{Order: 2, Instructions: "Minute 2", Reps: "10 Lunges (per leg)"},
					{Order: 3, Instructions: "Minute 3", Reps: "10 Burpees Modified"},
					{Order: 4, Instructions: "Minute 4", Reps: "20 Crunches"},
					{Order: 5, Instructions: "Minute 5", Reps: "15 Tricep Dips or Diamond Press Ups"},
				},
			},
		},
	})

	// --- DAY 9 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   9,
		Title:       "Plyometric Power (Pyramid)",
		Description: "A pyramid-style workout focused on explosive movements.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "Pyramid",
				BlockNotes: "Pyramid up and down: 2-4-6-8-10-12-14-16-14-12-10-8-6-4-2 reps. Rest when needed.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Tuck Jumps"]},
					{Order: 2, ExerciseID: exIDs["Plyo Press-ups"]},
					{Order: 3, ExerciseID: exIDs["Squat Jumps"]},
					{Order: 4, ExerciseID: exIDs["Explosive Starjumps"]},
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
				BlockType:  "Time Challenge",
				BlockNotes: "Complete for time. Break up reps as needed.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Ups"], Reps: "100"},
					{Order: 2, ExerciseID: exIDs["Tricep Dips"], Reps: "75"},
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
		Title:       "Core Focus & Full Body Bursts",
		Description: "A core EMOM followed by a full body circuit.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "EMOM",
				BlockNotes: "12 minutes total.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, Instructions: "Minute 1", ExerciseID: exIDs["Flutter Kicks"], Reps: "30"},
					{Order: 2, Instructions: "Minute 2", ExerciseID: exIDs["Leg Raises"], Reps: "15"},
					{Order: 3, Instructions: "Minute 3", ExerciseID: exIDs["Jack Knife"], Reps: "10"},
				},
			},
			{
				BlockType:   "Circuit",
				BlockRounds: 4,
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpee Sprints"], Reps: "12"},
					{Order: 2, ExerciseID: exIDs["Cross Jacks"], Reps: "12"},
					{Order: 3, ExerciseID: exIDs["Explosive Starjumps"], Reps: "12", Tips: "Explosive"},
				},
			},
		},
	})

	// --- DAY 12 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   12,
		Title:       "Core Domination",
		Description: "An EMOM workout focused on core strength.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "EMOM",
				BlockNotes: "Every 2 minutes for 20 minutes.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Crunches"], Reps: "20"},
					{Order: 2, ExerciseID: exIDs["Leg Raises"], Reps: "15"},
					{Order: 3, ExerciseID: exIDs["Plank Hold"], Duration: "30s"},
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
				BlockNotes:  "40s work, 20s rest.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpee Sprints"], Duration: "40s", RestDuration: "20s"},
					{Order: 2, ExerciseID: exIDs["Switch Kicks"], Duration: "40s", RestDuration: "20s"},
					{Order: 3, ExerciseID: exIDs["Bear Crawls"], Duration: "40s", RestDuration: "20s"},
					{Order: 4, ExerciseID: exIDs["Sprawls"], Duration: "40s", RestDuration: "20s"},
					{Order: 5, ExerciseID: exIDs["T-Runs"], Duration: "40s", RestDuration: "20s"},
					{Order: 6, ExerciseID: exIDs["Ski Jumps"], Duration: "40s", RestDuration: "20s"},
				},
			},
		},
	})

	// --- DAY 14 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   14,
		Title:       "Full Body AMRAP",
		Description: "An AMRAP circuit with rounds and rest.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "AMRAP",
				BlockRounds: 3,
				BlockNotes:  "60s AMRAP for each exercise, then 90s rest between rounds. Record reps. Each round try to do more reps than the last",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Jumps"], Duration: "60s", Reps: "AMRAP"},
					{Order: 2, ExerciseID: exIDs["Wide Arm Press-ups"], Duration: "60s", Reps: "AMRAP"},
					{Order: 3, ExerciseID: exIDs["Mountain Climbers"], Duration: "60s", Reps: "AMRAP"},
					{Order: 4, ExerciseID: exIDs["Glute Bridges"], Duration: "60s", Reps: "AMRAP"},
					{Order: 5, ExerciseID: exIDs["Burpees"], Duration: "60s", Reps: "AMRAP"},
				},
			},
		},
	})

	// --- DAY 15 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   15,
		Title:       "OPTIONAL MOBILITY DAY And/OR Full Body AMRAP & Ladders",
		Description: "Two distinct parts: AMRAP and a ladder workout.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "AMRAP",
				BlockNotes: "12 minutes total. Complete as many rounds as possible.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squat Jumps"], Reps: "10"},
					{Order: 2, ExerciseID: exIDs["Plyo Press-ups"], Reps: "8"},
					{Order: 3, ExerciseID: exIDs["Burpee Tuck Jumps"], Reps: "6"},
				},
			},
			{
				BlockType:  "Rest",
				BlockNotes: "2 minutes rest.",
			},
			{
				BlockType:  "Ladder",
				BlockNotes: "Pyramid up and down: 5-6-7-8-9-10-9-8-7-6-5 reps.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Inchworms"]},
					{Order: 2, ExerciseID: exIDs["Pike Jumps"]},
				},
			},
		},
	})

	// --- DAY 16 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   16,
		Title:       "Lower Body Power & Endurance",
		Description: "A timed circuit followed by a descending ladder.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 3,
				BlockNotes:  "50s work, 10s rest.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Jump Lunges"], Duration: "50s", RestDuration: "10s"},
					{Order: 2, ExerciseID: exIDs["Squat Kicks"], Duration: "50s", RestDuration: "10s"},
					{Order: 3, ExerciseID: exIDs["Ski Jumps"], Duration: "50s", RestDuration: "10s"},
				},
			},
			{
				BlockType:  "Ladder",
				BlockNotes: "Descending rep ladder for Broad Jumps and Thrusters.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Broad Jumps"], Reps: "2, 4, 6, 8, 10"},
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
		Description: "A 24-minute EMOM workout.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "EMOM",
				BlockRounds: 4, // 6 minutes per round
				BlockNotes:  "24 minutes total. Repeat the sequence 4 times.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, Instructions: "Minute 1", ExerciseID: exIDs["Diamond Sit-ups"], Reps: "15"},
					{Order: 2, Instructions: "Minute 2", ExerciseID: exIDs["Press Ups"], Reps: "15"},
					{Order: 3, Instructions: "Minute 3", ExerciseID: exIDs["Tricep Dips"], Reps: "15"},
					{Order: 4, Instructions: "Minute 4", ExerciseID: exIDs["Overhead Jabs (Fast)"], Reps: "20 per arm"},
					{Order: 5, Instructions: "Minute 5", ExerciseID: exIDs["Plank Hold"], Duration: "Max hold (40s)"},
					{Order: 6, Instructions: "Minute 6", Reps: "Rest"},
				},
			},
		},
	})

	// --- DAY 18 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   18,
		Title:       "Cardio Chaos",
		Description: "Density training: AMRAP in 20 minutes.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "AMRAP",
				BlockNotes: "20 minutes total. Complete as many rounds as possible.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpee Tucks"], Reps: "5"},
					{Order: 2, ExerciseID: exIDs["Jump Lunges"], Reps: "10 each leg"},
					{Order: 3, ExerciseID: exIDs["Mountain Climbers"], Reps: "15 each leg"},
					{Order: 4, ExerciseID: exIDs["High Knees"], Reps: "20 each leg"},
					{Order: 5, ExerciseID: exIDs["Heel Flicks"], Reps: "25 each leg"},
				},
			},
		},
	})

	// --- DAY 19 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   19,
		Title:       "Core & Cardio Tabata Challenge",
		Description: "Multiple Tabata blocks alternating between two exercises.",
		WorkoutBlocks: []models.WorkoutBlock{
			{BlockType: "Tabata", BlockNotes: "20s work / 10s rest x 10 rounds", Exercises: []models.WorkoutExercise{{Order: 1, ExerciseID: exIDs["Sprints"]}, {Order: 2, ExerciseID: exIDs["Knee to Chest"]}}},
			{BlockType: "Rest", BlockNotes: "Rest 60 seconds"},
			{BlockType: "Tabata", BlockNotes: "20s work / 10s rest x 10 rounds", Exercises: []models.WorkoutExercise{{Order: 1, ExerciseID: exIDs["Starjumps"]}, {Order: 2, ExerciseID: exIDs["Bicycle Legs"]}}},
			{BlockType: "Rest", BlockNotes: "Rest 60 seconds"},
			{BlockType: "Tabata", BlockNotes: "20s work / 10s rest x 10 rounds", Exercises: []models.WorkoutExercise{{Order: 1, ExerciseID: exIDs["High Knees"]}, {Order: 2, ExerciseID: exIDs["Diamond Sit-ups"]}}},
			{BlockType: "Rest", BlockNotes: "Rest 60 seconds"},
			{BlockType: "Tabata", BlockNotes: "20s work / 10s rest x 10 rounds", Exercises: []models.WorkoutExercise{{Order: 1, ExerciseID: exIDs["Thrusters"]}, {Order: 2, ExerciseID: exIDs["Ab Twists"]}}},
			{BlockType: "Rest", BlockNotes: "Rest 60 seconds"},
			{BlockType: "Tabata", BlockNotes: "20s work / 10s rest x 10 rounds", Exercises: []models.WorkoutExercise{{Order: 1, ExerciseID: exIDs["Belt Kicks"]}, {Order: 2, ExerciseID: exIDs["Leg Raises"]}}},
		},
	})

	// --- DAY 20 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   20,
		Title:       "Full Body Fusion",
		Description: "Complex training circuit.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 5,
				BlockNotes:  "Rest 2 minutes between rounds.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Thrusters"], Reps: "8"},
					{Order: 2, ExerciseID: exIDs["Burpees"], Reps: "10"},
					{Order: 3, ExerciseID: exIDs["Squat Twists"], Reps: "12"},
					{Order: 4, ExerciseID: exIDs["Press Ups"], Reps: "14", Tips: "With Twists"},
					{Order: 5, ExerciseID: exIDs["Starjumps"], Reps: "16"},
				},
			},
		},
	})

	// --- DAY 21 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   21,
		Title:       "Endurance Test",
		Description: "Complete for time. Reps can be broken up as needed.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "Endurance Test",
				BlockNotes: "Complete for time. Record total time.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squats"], Reps: "200"},
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
		Title:       "OPTIONAL MOBILITY DAY AND/OR Full Body Pyramids & AMRAP",
		Description: "A pyramid workout followed by a short finisher.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "Pyramid",
				BlockNotes: "Pyramid up and down: 5-10-15-10-5 reps for Burpees and 10-20-30-20-10 reps for Squat Jumps.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpees"], Reps: "5, 10, 15, 10, 5"},
					{Order: 2, ExerciseID: exIDs["Squat Jumps"], Reps: "10, 20, 30, 20, 10"},
				},
			},
			{BlockType: "Rest", BlockNotes: "Rest 2 minutes"},
			{
				BlockType:  "AMRAP",
				BlockNotes: "5 minutes total.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Tuck Jumps"], Reps: "5"},
					{Order: 2, ExerciseID: exIDs["Wide Arm Press-ups"], Reps: "10"},
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
		Description: "A timed circuit, a pyramid workout, and a finisher.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 3,
				BlockNotes:  "50s work, 10s rest.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Up Twists"], Tips: "With Twists", Duration: "50s", RestDuration: "10s"},
					{Order: 2, ExerciseID: exIDs["Oblique Press-ups"], Duration: "50s", RestDuration: "10s"},
					{Order: 3, ExerciseID: exIDs["Plank Leg Raises"], Duration: "50s", RestDuration: "10s"},
					{Order: 4, ExerciseID: exIDs["Plank Hold"], Duration: "50s", RestDuration: "10s"},
				},
			},
			{
				BlockType:  "Pyramid",
				BlockNotes: "Pyramid up and down: 15-10-5-10-15 V Press-ups and 25-20-15-20-25 Half Sits.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["V Press-ups"], Reps: "15, 10, 5, 10, 15"},
					{Order: 2, ExerciseID: exIDs["Half Sit Ups"], Reps: "25, 20, 15, 20, 25"},
				},
			},
			{
				BlockType:  "Finisher",
				BlockNotes: "Complete after the main workout.",
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
		Description: "A timed circuit, an AMRAP, and a static hold finisher.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 4,
				BlockNotes:  "50s work, 10s rest.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["T-Runs"], Duration: "50s", RestDuration: "10s"},
					{Order: 2, ExerciseID: exIDs["Y-shaped Lunges"], Duration: "50s", RestDuration: "10s"},
					{Order: 3, ExerciseID: exIDs["Squat Twists"], Duration: "50s", RestDuration: "10s"},
					{Order: 4, ExerciseID: exIDs["Calf Jumps"], Duration: "50s", RestDuration: "10s"},
				},
			},
			{
				BlockType:  "AMRAP",
				BlockNotes: "10 minutes total.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Switch Kicks"], Reps: "10 per leg"},
					{Order: 2, ExerciseID: exIDs["Thrusters"], Reps: "10"},
					{Order: 3, ExerciseID: exIDs["Broad Jumps"], Reps: "5"},
				},
			},
			{
				BlockType:  "Finisher",
				BlockNotes: "Static Holds: 1 minute each.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squats"], Tips: "Hold", Duration: "1 min"},
					{Order: 2, ExerciseID: exIDs["Wall Sit"], Tips: "Hold", Duration: "1 min"},
					{Order: 3, ExerciseID: exIDs["Hollow Rock Holds"], Tips: "Hold", Duration: "1 min"},
					{Order: 4, ExerciseID: exIDs["Plack Hold"], Tips: "Hold", Duration: "1 min"},
				},
			},
		},
	})

	// --- DAY 25 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   25,
		Title:       "Plyo Push",
		Description: "A plyometric workout with a plank challenge finisher.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 4,
				BlockNotes:  "Rest 60 seconds between rounds.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Plyo Press-ups"], Reps: "12"},
					{Order: 2, ExerciseID: exIDs["Tuck Jumps"], Reps: "20"},
					{Order: 3, ExerciseID: exIDs["HOG Press-ups"], Reps: "12"},
					{Order: 4, ExerciseID: exIDs["Starjumps"], Reps: "20"},
					{Order: 5, ExerciseID: exIDs["Moving Press-ups"], Reps: "12"},
					{Order: 6, ExerciseID: exIDs["Ski Jumps"], Reps: "20"},
					{Order: 7, ExerciseID: exIDs["Oblique Hops"], Reps: "12"},
					{Order: 8, ExerciseID: exIDs["Jump Lunges"], Reps: "20"},
				},
			},
			{
				BlockType:  "Finisher",
				BlockNotes: "Plank Challenge: Max time or 5-minute total hold.",
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
		Description: "A sets-and-reps workout with a press-up finisher.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 4,
				BlockNotes:  "Complete 15 reps of each. Rest 30-60s between sets.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Diamond Sit-ups"], Reps: "15"},
					{Order: 2, ExerciseID: exIDs["High Low Plank"], Reps: "15"},
					{Order: 3, ExerciseID: exIDs["Press Ups"], Reps: "15"},
					{Order: 4, ExerciseID: exIDs["Tricep Dips"], Reps: "15"},
					{Order: 5, ExerciseID: exIDs["Bicycles"], Reps: "15"},
					{Order: 6, ExerciseID: exIDs["Sit ups"], Reps: "15"},
					{Order: 7, ExerciseID: exIDs["HOG Press-ups"], Reps: "15"},
					{Order: 8, ExerciseID: exIDs["Plank Jabs"], Reps: "15"},
				},
			},
			{
				BlockType:  "Finisher",
				BlockNotes: "EMOM Finisher: Push-Up Variations. 12 reps in minute 1, increase by 2 reps each minute for 5 minutes.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Ups"], Reps: "EMOM increasing reps"},
				},
			},
		},
	})

	// --- DAY 27 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   27,
		Title:       "Full Body Cardio and Agility",
		Description: "A timed circuit followed by a 'Death by Burpees' finisher.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:   "Circuit",
				BlockRounds: 5,
				BlockNotes:  "50s work, 10s rest.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Burpee Sprints"], Duration: "50s", RestDuration: "10s"},
					{Order: 2, ExerciseID: exIDs["Thrusters"], Duration: "50s", RestDuration: "10s"},
					{Order: 3, ExerciseID: exIDs["Squat Jumps"], Duration: "50s", RestDuration: "10s"},
					{Order: 4, ExerciseID: exIDs["T-Runs"], Duration: "50s", RestDuration: "10s"},
					{Order: 5, ExerciseID: exIDs["High Knees"], Duration: "50s", RestDuration: "10s"},
					{Order: 6, ExerciseID: exIDs["Mountain Climbers"], Duration: "50s", RestDuration: "10s"},
				},
			},
			{
				BlockType:  "Finisher",
				BlockNotes: "Death by Burpees: 6 reps first minute, 8 reps second minute, etc., until failure.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Death by Burpees"], Reps: "Increasing Reps until Failure"},
				},
			},
		},
	})

	// --- DAY 28 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   28,
		Title:       "Endurance Workout",
		Description: "A timed max effort workout with a ladder finisher.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "Timed Max Effort",
				BlockNotes: "1 minute rest between exercises.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squats"], Duration: "5 min", Reps: "Max"},
					{Order: 2, ExerciseID: exIDs["Press Ups"], Duration: "4 min", Reps: "Max"},
					{Order: 3, ExerciseID: exIDs["Mountain Climbers"], Duration: "3 min", Reps: "Max"},
					{Order: 4, ExerciseID: exIDs["Burpees"], Duration: "2 min", Reps: "Max"},
					{Order: 5, ExerciseID: exIDs["Tuck Jumps"], Duration: "1 min", Reps: "Max"},
				},
			},
			{
				BlockType:  "Finisher",
				BlockNotes: "10-9-8-7-6-5-4-3-2-1 ladder of any two exercises of your choice.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, Tips: "User's Choice"},
					{Order: 2, Tips: "User's Choice"},
				},
			},
		},
	})

	// --- DAY 29 ---
	createWorkoutDay(models.WorkoutDay{
		ProgramID:   programID,
		DayNumber:   29,
		Title:       "Peak Performance (Grand Finisher)",
		Description: "A final, grueling workout with three options for completion.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "Grand Finisher",
				BlockNotes: "Total reps: 200 Squats, 150 M-Climbers, 100 Burpees, 75 Sit-ups, 50 Tuck Jumps.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Squats"], Reps: "200"},
					{Order: 2, ExerciseID: exIDs["Mountain Climbers"], Reps: "150"},
					{Order: 3, ExerciseID: exIDs["Burpees"], Reps: "100"},
					{Order: 4, ExerciseID: exIDs["Sit ups"], Reps: "75"},
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
		Description: "Repeat Day 1 assessment to track progress.",
		WorkoutBlocks: []models.WorkoutBlock{
			{
				BlockType:  "Fitness Assessment",
				BlockNotes: "Compare results with Day 1 to track progress.",
				Exercises: []models.WorkoutExercise{
					{Order: 1, ExerciseID: exIDs["Press Ups"], Reps: "Max Reps", Duration: "1 min"},
					{Order: 2, ExerciseID: exIDs["Straddle Sit-ups"], Reps: "Max Reps", Duration: "1 min"},
					{Order: 3, ExerciseID: exIDs["Plank Hold"], Reps: "Max Time"},
					{Order: 4, ExerciseID: exIDs["Squats"], Reps: "Max Reps", Duration: "1 min"},
					{Order: 5, ExerciseID: exIDs["Burpees"], Reps: "Max Reps", Duration: "3 min"},
					{Order: 6, ExerciseID: exIDs["Jump Lunges"], Reps: "Max Reps", Duration: "1 min"},
				},
			},
		},
	})

	log.Println("Successfully seeded data for the entire 30-Day Advanced Programme.")
}
