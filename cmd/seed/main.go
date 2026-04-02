package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/joho/godotenv"
	wr "github.com/mroth/weightedrand/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Location represents a geographic location.
type Location struct {
	City     string `bson:"city"`
	State    string `bson:"state"`
	Country  string `bson:"country"`
	Timezone string `bson:"timezone"`
	Region   string `bson:"region"`
}

// Call represents a recruiter call to a candidate.
type Call struct {
	Date     time.Time `bson:"date"`
	Duration int       `bson:"duration"`
	Outcome  string    `bson:"outcome"`
	Notes    string    `bson:"notes"`
}

// Recruiter represents a recruiter record.
type Recruiter struct {
	ID         string `bson:"id"`
	Name       string `bson:"name"`
	Email      string `bson:"email"`
	Phone      string `bson:"phone"`
	Department string `bson:"department"`
	HireDate   string `bson:"hire_date"`
}

// RecruiterRef is the embedded recruiter reference in a candidate.
type RecruiterRef struct {
	ID   string `bson:"id"`
	Name string `bson:"name"`
}

// Candidate represents a candidate record.
type Candidate struct {
	Name           string       `bson:"name"`
	Email          string       `bson:"email"`
	Phone          string       `bson:"phone"`
	Location       Location     `bson:"location"`
	Industry       string       `bson:"industry"`
	RoleCategory   string       `bson:"role_category"`
	RoleTitle      string       `bson:"role_title"`
	Seniority      string       `bson:"seniority"`
	YearsExp       int          `bson:"years_experience"`
	Status         string       `bson:"status"`
	Recruiter      RecruiterRef `bson:"recruiter"`
	Calls          []Call       `bson:"calls"`
	CreatedAt      time.Time    `bson:"created_at"`
	FirstContactAt *time.Time   `bson:"first_contact_at,omitempty"`
	ScreeningDate  *time.Time   `bson:"screening_date,omitempty"`
	QualifiedDate  *time.Time   `bson:"qualified_date,omitempty"`
	SalaryExpMin   int          `bson:"salary_expectation_min"`
	SalaryExpMax   int          `bson:"salary_expectation_max"`
	Source         string       `bson:"source"`
	Tags           []string     `bson:"tags"`
}

var recruiters = []Recruiter{
	{ID: "rec_001", Name: "Sarah Chen", Email: "sarah.chen@company.com", Phone: "+1-555-0101", Department: "Engineering Recruiting", HireDate: "2022-03-15"},
	{ID: "rec_002", Name: "James Wilson", Email: "james.wilson@company.com", Phone: "+1-555-0102", Department: "General Recruiting", HireDate: "2021-07-20"},
	{ID: "rec_003", Name: "Maria Garcia", Email: "maria.garcia@company.com", Phone: "+1-555-0103", Department: "Technical Recruiting", HireDate: "2023-01-10"},
	{ID: "rec_004", Name: "David Kim", Email: "david.kim@company.com", Phone: "+1-555-0104", Department: "Engineering Recruiting", HireDate: "2022-09-01"},
	{ID: "rec_005", Name: "Emily Brown", Email: "emily.brown@company.com", Phone: "+1-555-0105", Department: "Executive Recruiting", HireDate: "2020-11-05"},
	{ID: "rec_006", Name: "Alex Patel", Email: "alex.patel@company.com", Phone: "+1-555-0106", Department: "Data Science Recruiting", HireDate: "2023-06-15"},
	{ID: "rec_007", Name: "Lisa Thompson", Email: "lisa.thompson@company.com", Phone: "+1-555-0107", Department: "General Recruiting", HireDate: "2021-04-22"},
	{ID: "rec_008", Name: "Ryan Martinez", Email: "ryan.martinez@company.com", Phone: "+1-555-0108", Department: "Technical Recruiting", HireDate: "2022-12-01"},
}

var locationPools = map[string][]Location{
	"us_west": {
		{City: "San Francisco", State: "CA", Country: "USA", Timezone: "America/Los_Angeles"},
		{City: "Seattle", State: "WA", Country: "USA", Timezone: "America/Los_Angeles"},
		{City: "Los Angeles", State: "CA", Country: "USA", Timezone: "America/Los_Angeles"},
		{City: "Portland", State: "OR", Country: "USA", Timezone: "America/Los_Angeles"},
		{City: "San Jose", State: "CA", Country: "USA", Timezone: "America/Los_Angeles"},
	},
	"us_east": {
		{City: "New York", State: "NY", Country: "USA", Timezone: "America/New_York"},
		{City: "Boston", State: "MA", Country: "USA", Timezone: "America/New_York"},
		{City: "Washington", State: "DC", Country: "USA", Timezone: "America/New_York"},
		{City: "Philadelphia", State: "PA", Country: "USA", Timezone: "America/New_York"},
		{City: "Atlanta", State: "GA", Country: "USA", Timezone: "America/New_York"},
	},
	"us_central": {
		{City: "Chicago", State: "IL", Country: "USA", Timezone: "America/Chicago"},
		{City: "Austin", State: "TX", Country: "USA", Timezone: "America/Chicago"},
		{City: "Dallas", State: "TX", Country: "USA", Timezone: "America/Chicago"},
		{City: "Denver", State: "CO", Country: "USA", Timezone: "America/Denver"},
		{City: "Minneapolis", State: "MN", Country: "USA", Timezone: "America/Chicago"},
	},
	"canada": {
		{City: "Toronto", State: "ON", Country: "Canada", Timezone: "America/Toronto"},
		{City: "Vancouver", State: "BC", Country: "Canada", Timezone: "America/Vancouver"},
		{City: "Montreal", State: "QC", Country: "Canada", Timezone: "America/Toronto"},
		{City: "Ottawa", State: "ON", Country: "Canada", Timezone: "America/Toronto"},
	},
	"uk": {
		{City: "London", State: "England", Country: "UK", Timezone: "Europe/London"},
		{City: "Manchester", State: "England", Country: "UK", Timezone: "Europe/London"},
		{City: "Edinburgh", State: "Scotland", Country: "UK", Timezone: "Europe/London"},
		{City: "Bristol", State: "England", Country: "UK", Timezone: "Europe/London"},
	},
	"europe": {
		{City: "Berlin", State: "Berlin", Country: "Germany", Timezone: "Europe/Berlin"},
		{City: "Amsterdam", State: "North Holland", Country: "Netherlands", Timezone: "Europe/Amsterdam"},
		{City: "Paris", State: "Île-de-France", Country: "France", Timezone: "Europe/Paris"},
		{City: "Dublin", State: "Leinster", Country: "Ireland", Timezone: "Europe/Dublin"},
		{City: "Stockholm", State: "Stockholm", Country: "Sweden", Timezone: "Europe/Stockholm"},
	},
	"india": {
		{City: "Bangalore", State: "Karnataka", Country: "India", Timezone: "Asia/Kolkata"},
		{City: "Hyderabad", State: "Telangana", Country: "India", Timezone: "Asia/Kolkata"},
		{City: "Mumbai", State: "Maharashtra", Country: "India", Timezone: "Asia/Kolkata"},
		{City: "Pune", State: "Maharashtra", Country: "India", Timezone: "Asia/Kolkata"},
		{City: "Chennai", State: "Tamil Nadu", Country: "India", Timezone: "Asia/Kolkata"},
		{City: "Delhi", State: "Delhi", Country: "India", Timezone: "Asia/Kolkata"},
	},
	"southeast_asia": {
		{City: "Singapore", State: "Singapore", Country: "Singapore", Timezone: "Asia/Singapore"},
		{City: "Ho Chi Minh City", State: "Ho Chi Minh", Country: "Vietnam", Timezone: "Asia/Ho_Chi_Minh"},
		{City: "Bangkok", State: "Bangkok", Country: "Thailand", Timezone: "Asia/Bangkok"},
		{City: "Jakarta", State: "Jakarta", Country: "Indonesia", Timezone: "Asia/Jakarta"},
		{City: "Manila", State: "Metro Manila", Country: "Philippines", Timezone: "Asia/Manila"},
	},
	"latam": {
		{City: "São Paulo", State: "São Paulo", Country: "Brazil", Timezone: "America/Sao_Paulo"},
		{City: "Mexico City", State: "CDMX", Country: "Mexico", Timezone: "America/Mexico_City"},
		{City: "Buenos Aires", State: "Buenos Aires", Country: "Argentina", Timezone: "America/Argentina/Buenos_Aires"},
		{City: "Bogotá", State: "Cundinamarca", Country: "Colombia", Timezone: "America/Bogota"},
		{City: "Santiago", State: "Santiago", Country: "Chile", Timezone: "America/Santiago"},
	},
}

var roleTitlePools = map[string][]string{
	"engineering": {
		"Software Engineer", "Backend Engineer", "Frontend Engineer",
		"Full Stack Engineer", "Mobile Engineer", "Platform Engineer",
		"Infrastructure Engineer", "Site Reliability Engineer",
		"Security Engineer", "QA Engineer",
	},
	"data_science": {
		"Data Scientist", "Machine Learning Engineer", "Data Engineer",
		"Data Analyst", "ML Research Scientist", "AI Engineer",
		"NLP Engineer", "Computer Vision Engineer",
	},
	"devops": {
		"DevOps Engineer", "Cloud Engineer", "Systems Engineer",
		"Release Engineer", "Infrastructure Engineer", "Platform Engineer",
	},
	"product": {
		"Product Manager", "Technical Product Manager", "Product Analyst",
		"Product Owner", "Growth Product Manager",
	},
	"design": {
		"UX Designer", "UI Designer", "Product Designer",
		"UX Researcher", "Design Systems Engineer",
	},
	"management": {
		"Engineering Manager", "VP of Engineering", "Director of Engineering",
		"CTO", "Technical Lead", "Team Lead",
	},
}

var sources = []string{
	"linkedin", "referral", "job_board", "direct_apply",
	"recruiter_outreach", "career_fair", "github", "stack_overflow",
}

var tagPool = []string{
	"remote_ok", "visa_required", "relocation", "contract_to_hire",
	"urgent", "backfill", "new_headcount", "diversity_candidate",
	"internal_referral", "boomerang", "top_school", "faang_experience",
}

func main() {
	godotenv.Load()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	fmt.Println("Connected to MongoDB")

	db := client.Database("recruitment")

	// Drop existing collections
	db.Collection("candidates").Drop(ctx)
	db.Collection("recruiters").Drop(ctx)
	fmt.Println("Dropped existing collections")

	// Insert recruiters
	recDocs := make([]interface{}, len(recruiters))
	for i, r := range recruiters {
		recDocs[i] = r
	}
	_, err = db.Collection("recruiters").InsertMany(ctx, recDocs)
	if err != nil {
		log.Fatalf("Failed to insert recruiters: %v", err)
	}
	fmt.Printf("Inserted %d recruiters\n", len(recruiters))

	// Build weighted choosers
	regionChooser, _ := wr.NewChooser(
		wr.NewChoice("us_west", 15),
		wr.NewChoice("us_east", 13),
		wr.NewChoice("us_central", 12),
		wr.NewChoice("canada", 10),
		wr.NewChoice("uk", 8),
		wr.NewChoice("europe", 7),
		wr.NewChoice("india", 15),
		wr.NewChoice("southeast_asia", 10),
		wr.NewChoice("latam", 10),
	)

	industryChooser, _ := wr.NewChooser(
		wr.NewChoice("saas", 25),
		wr.NewChoice("ai_ml", 20),
		wr.NewChoice("fintech", 15),
		wr.NewChoice("cloud_infra", 15),
		wr.NewChoice("healthtech", 10),
		wr.NewChoice("ecommerce", 5),
		wr.NewChoice("cybersecurity", 5),
		wr.NewChoice("martech", 5),
	)

	roleCatChooser, _ := wr.NewChooser(
		wr.NewChoice("engineering", 35),
		wr.NewChoice("data_science", 25),
		wr.NewChoice("devops", 15),
		wr.NewChoice("product", 10),
		wr.NewChoice("design", 8),
		wr.NewChoice("management", 7),
	)

	seniorityChooser, _ := wr.NewChooser(
		wr.NewChoice("junior", 15),
		wr.NewChoice("mid", 30),
		wr.NewChoice("senior", 30),
		wr.NewChoice("staff", 12),
		wr.NewChoice("principal", 5),
		wr.NewChoice("director", 5),
		wr.NewChoice("vp", 3),
	)

	statusChooser, _ := wr.NewChooser(
		wr.NewChoice("dialed", 20),
		wr.NewChoice("no_answer", 15),
		wr.NewChoice("connected", 20),
		wr.NewChoice("screened", 18),
		wr.NewChoice("qualified", 12),
		wr.NewChoice("rejected", 8),
		wr.NewChoice("hired", 4),
		wr.NewChoice("callback_scheduled", 3),
	)

	now := time.Now()
	candidates := make([]interface{}, 500)

	for i := 0; i < 500; i++ {
		region := regionChooser.Pick()
		industry := industryChooser.Pick()
		roleCat := roleCatChooser.Pick()
		seniority := seniorityChooser.Pick()
		status := statusChooser.Pick()

		// Pick location from region pool
		locs := locationPools[region]
		loc := locs[rand.Intn(len(locs))]
		loc.Region = region

		// Pick role title from category pool
		titles := roleTitlePools[roleCat]
		roleTitle := titles[rand.Intn(len(titles))]

		// Assign recruiter
		rec := assignRecruiter(roleCat)

		// Years of experience based on seniority
		yearsExp := yearsForSeniority(seniority)

		// Salary expectations based on seniority
		salaryMin, salaryMax := salaryForSeniority(seniority)

		// Created at: random within last 90 days
		createdAt := now.Add(-time.Duration(rand.Intn(90*24)) * time.Hour)
		createdAt = createdAt.Add(-time.Duration(rand.Intn(60)) * time.Minute)

		// Generate calls
		calls := generateCalls(status, createdAt, now)

		// Date logic
		var firstContactAt, screeningDate, qualifiedDate *time.Time
		if hasConnectedCall(calls) {
			t := createdAt.Add(time.Duration(rand.Intn(3*24)) * time.Hour)
			firstContactAt = &t
		}
		if status == "screened" || status == "qualified" || status == "hired" || status == "rejected" {
			if firstContactAt != nil {
				t := firstContactAt.Add(time.Duration(1+rand.Intn(13)) * 24 * time.Hour)
				screeningDate = &t
			}
		}
		if status == "qualified" || status == "hired" {
			if screeningDate != nil {
				t := screeningDate.Add(time.Duration(1+rand.Intn(6)) * 24 * time.Hour)
				qualifiedDate = &t
			}
		}

		// Source
		source := sources[rand.Intn(len(sources))]

		// Tags (1-3 random)
		numTags := 1 + rand.Intn(3)
		tags := pickRandomTags(numTags)

		candidate := Candidate{
			Name:           gofakeit.Name(),
			Email:          gofakeit.Email(),
			Phone:          gofakeit.Phone(),
			Location:       loc,
			Industry:       industry,
			RoleCategory:   roleCat,
			RoleTitle:      roleTitle,
			Seniority:      seniority,
			YearsExp:       yearsExp,
			Status:         status,
			Recruiter:      rec,
			Calls:          calls,
			CreatedAt:      createdAt,
			FirstContactAt: firstContactAt,
			ScreeningDate:  screeningDate,
			QualifiedDate:  qualifiedDate,
			SalaryExpMin:   salaryMin,
			SalaryExpMax:   salaryMax,
			Source:         source,
			Tags:           tags,
		}
		candidates[i] = candidate
	}

	_, err = db.Collection("candidates").InsertMany(ctx, candidates)
	if err != nil {
		log.Fatalf("Failed to insert candidates: %v", err)
	}
	fmt.Printf("Inserted %d candidates\n", len(candidates))

	// Create indexes
	coll := db.Collection("candidates")
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "recruiter.id", Value: 1}}},
		{Keys: bson.D{{Key: "industry", Value: 1}}},
		{Keys: bson.D{{Key: "role_category", Value: 1}}},
		{Keys: bson.D{{Key: "location.region", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: 1}}},
	}
	_, err = coll.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}
	fmt.Println("Created indexes")

	// Print summary
	count, _ := coll.CountDocuments(ctx, bson.D{})
	recCount, _ := db.Collection("recruiters").CountDocuments(ctx, bson.D{})
	fmt.Printf("\nSummary:\n  Candidates: %d\n  Recruiters: %d\n", count, recCount)

	// Status breakdown
	cursor, _ := coll.Aggregate(ctx, mongo.Pipeline{
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$status"}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
	})
	var results []bson.M
	cursor.All(ctx, &results)
	fmt.Println("  Status breakdown:")
	for _, r := range results {
		fmt.Printf("    %s: %v\n", r["_id"], r["count"])
	}
}

func assignRecruiter(roleCat string) RecruiterRef {
	if roleCat == "engineering" {
		r := rand.Float64()
		if r < 0.30 {
			return RecruiterRef{ID: "rec_001", Name: "Sarah Chen"}
		} else if r < 0.60 {
			return RecruiterRef{ID: "rec_004", Name: "David Kim"}
		}
	}
	// Even distribution across all 8
	rec := recruiters[rand.Intn(len(recruiters))]
	return RecruiterRef{ID: rec.ID, Name: rec.Name}
}

func yearsForSeniority(seniority string) int {
	switch seniority {
	case "junior":
		return rand.Intn(3)
	case "mid":
		return 2 + rand.Intn(4)
	case "senior":
		return 5 + rand.Intn(6)
	case "staff":
		return 8 + rand.Intn(7)
	case "principal":
		return 12 + rand.Intn(8)
	case "director":
		return 10 + rand.Intn(10)
	case "vp":
		return 15 + rand.Intn(10)
	default:
		return 3
	}
}

func salaryForSeniority(seniority string) (int, int) {
	switch seniority {
	case "junior":
		base := 60000 + rand.Intn(30000)
		return base, base + 15000 + rand.Intn(10000)
	case "mid":
		base := 90000 + rand.Intn(40000)
		return base, base + 20000 + rand.Intn(15000)
	case "senior":
		base := 140000 + rand.Intn(50000)
		return base, base + 25000 + rand.Intn(20000)
	case "staff":
		base := 180000 + rand.Intn(60000)
		return base, base + 30000 + rand.Intn(25000)
	case "principal":
		base := 220000 + rand.Intn(80000)
		return base, base + 40000 + rand.Intn(30000)
	case "director":
		base := 200000 + rand.Intn(80000)
		return base, base + 40000 + rand.Intn(30000)
	case "vp":
		base := 250000 + rand.Intn(100000)
		return base, base + 50000 + rand.Intn(40000)
	default:
		return 80000, 100000
	}
}

func generateCalls(status string, createdAt, now time.Time) []Call {
	var calls []Call
	spread := now.Sub(createdAt)

	switch status {
	case "no_answer":
		n := 1 + rand.Intn(3)
		for i := 0; i < n; i++ {
			outcome := "no_answer"
			if rand.Float64() < 0.3 {
				outcome = "voicemail"
			}
			calls = append(calls, Call{
				Date:     createdAt.Add(time.Duration(rand.Int63n(int64(spread + 1)))),
				Duration: 0,
				Outcome:  outcome,
				Notes:    "",
			})
		}

	case "dialed":
		n := 1 + rand.Intn(5)
		for i := 0; i < n; i++ {
			outcome := "no_answer"
			dur := 0
			if i == n-1 && rand.Float64() < 0.4 {
				outcome = "connected"
				dur = rand.Intn(60)
			}
			calls = append(calls, Call{
				Date:     createdAt.Add(time.Duration(rand.Int63n(int64(spread + 1)))),
				Duration: dur,
				Outcome:  outcome,
				Notes:    "",
			})
		}

	case "connected":
		n := 2 + rand.Intn(5)
		for i := 0; i < n; i++ {
			outcome := "no_answer"
			dur := 0
			if i == 1 {
				outcome = "connected"
				dur = 60 + rand.Intn(240)
			} else if rand.Float64() < 0.3 {
				outcome = "connected"
				dur = 30 + rand.Intn(120)
			}
			calls = append(calls, Call{
				Date:     createdAt.Add(time.Duration(rand.Int63n(int64(spread + 1)))),
				Duration: dur,
				Outcome:  outcome,
				Notes:    "",
			})
		}

	case "screened":
		n := 3 + rand.Intn(6)
		for i := 0; i < n; i++ {
			outcome := "no_answer"
			dur := 0
			if i == 1 {
				outcome = "connected"
				dur = 60 + rand.Intn(240)
			} else if i == 2 {
				outcome = "screened"
				dur = 600 + rand.Intn(1200)
			} else if rand.Float64() < 0.2 {
				outcome = "connected"
				dur = 30 + rand.Intn(120)
			}
			calls = append(calls, Call{
				Date:     createdAt.Add(time.Duration(rand.Int63n(int64(spread + 1)))),
				Duration: dur,
				Outcome:  outcome,
				Notes:    "",
			})
		}

	case "qualified", "hired":
		n := 4 + rand.Intn(7)
		for i := 0; i < n; i++ {
			outcome := "no_answer"
			dur := 0
			if i == 1 {
				outcome = "connected"
				dur = 60 + rand.Intn(240)
			} else if i == 2 {
				outcome = "screened"
				dur = 900 + rand.Intn(1500)
			} else if i == 3 {
				outcome = "qualified"
				dur = 300 + rand.Intn(600)
			} else if rand.Float64() < 0.2 {
				outcome = "connected"
				dur = 30 + rand.Intn(120)
			}
			calls = append(calls, Call{
				Date:     createdAt.Add(time.Duration(rand.Int63n(int64(spread + 1)))),
				Duration: dur,
				Outcome:  outcome,
				Notes:    "",
			})
		}

	case "callback_scheduled":
		n := 1 + rand.Intn(4)
		for i := 0; i < n; i++ {
			outcome := "no_answer"
			dur := 0
			if rand.Float64() < 0.4 {
				outcome = "connected"
				dur = 30 + rand.Intn(90)
			}
			calls = append(calls, Call{
				Date:     createdAt.Add(time.Duration(rand.Int63n(int64(spread + 1)))),
				Duration: dur,
				Outcome:  outcome,
				Notes:    "",
			})
		}

	case "rejected":
		n := 2 + rand.Intn(5)
		for i := 0; i < n; i++ {
			outcome := "no_answer"
			dur := 0
			if i == 1 {
				outcome = "connected"
				dur = 60 + rand.Intn(180)
			} else if i == 2 {
				outcome = "screened"
				dur = 300 + rand.Intn(600)
			} else if rand.Float64() < 0.2 {
				outcome = "connected"
				dur = 30 + rand.Intn(90)
			}
			calls = append(calls, Call{
				Date:     createdAt.Add(time.Duration(rand.Int63n(int64(spread + 1)))),
				Duration: dur,
				Outcome:  outcome,
				Notes:    "",
			})
		}
	}

	return calls
}

func hasConnectedCall(calls []Call) bool {
	for _, c := range calls {
		if c.Outcome == "connected" || c.Outcome == "screened" || c.Outcome == "qualified" {
			return true
		}
	}
	return false
}

func pickRandomTags(n int) []string {
	perm := rand.Perm(len(tagPool))
	tags := make([]string, n)
	for i := 0; i < n; i++ {
		tags[i] = tagPool[perm[i]]
	}
	return tags
}
