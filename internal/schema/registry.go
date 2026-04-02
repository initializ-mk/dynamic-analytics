package schema

// SchemaInfo describes the MongoDB schema for LLM context and the /api/schema endpoint.
type FieldInfo struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Values      []string `json:"values,omitempty"`
}

type CollectionSchema struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Fields      []FieldInfo `json:"fields"`
}

func GetCandidatesSchema() CollectionSchema {
	return CollectionSchema{
		Name:        "candidates",
		Description: "AI recruitment candidates with call history, status tracking, and recruiter assignments",
		Fields: []FieldInfo{
			{Name: "name", Type: "string", Description: "Candidate full name"},
			{Name: "email", Type: "string", Description: "Candidate email address"},
			{Name: "phone", Type: "string", Description: "Candidate phone number"},
			{Name: "location.city", Type: "string", Description: "City name"},
			{Name: "location.state", Type: "string", Description: "State or province"},
			{Name: "location.country", Type: "string", Description: "Country name"},
			{Name: "location.region", Type: "string", Description: "Geographic region",
				Values: []string{"us_west", "us_east", "us_central", "canada", "uk", "europe", "india", "southeast_asia", "latam"}},
			{Name: "location.timezone", Type: "string", Description: "IANA timezone"},
			{Name: "industry", Type: "string", Description: "Target industry",
				Values: []string{"saas", "ai_ml", "fintech", "cloud_infra", "healthtech", "ecommerce", "cybersecurity", "martech"}},
			{Name: "role_category", Type: "string", Description: "Role category",
				Values: []string{"engineering", "data_science", "devops", "product", "design", "management"}},
			{Name: "role_title", Type: "string", Description: "Specific job title"},
			{Name: "seniority", Type: "string", Description: "Seniority level",
				Values: []string{"junior", "mid", "senior", "staff", "principal", "director", "vp"}},
			{Name: "years_experience", Type: "int", Description: "Years of professional experience"},
			{Name: "status", Type: "string", Description: "Current pipeline status",
				Values: []string{"dialed", "no_answer", "connected", "screened", "qualified", "rejected", "hired", "callback_scheduled"}},
			{Name: "recruiter.id", Type: "string", Description: "Recruiter ID (rec_001 through rec_008)"},
			{Name: "recruiter.name", Type: "string", Description: "Recruiter full name"},
			{Name: "calls", Type: "array", Description: "Array of call records with fields: date (datetime), duration (int, seconds), outcome (string: no_answer/voicemail/connected/screened/qualified), notes (string)"},
			{Name: "created_at", Type: "datetime", Description: "When the candidate was added to the system"},
			{Name: "first_contact_at", Type: "datetime", Description: "When first successful contact was made (nullable)"},
			{Name: "screening_date", Type: "datetime", Description: "When screening call occurred (nullable)"},
			{Name: "qualified_date", Type: "datetime", Description: "When candidate was qualified (nullable)"},
			{Name: "salary_expectation_min", Type: "int", Description: "Minimum salary expectation in USD"},
			{Name: "salary_expectation_max", Type: "int", Description: "Maximum salary expectation in USD"},
			{Name: "source", Type: "string", Description: "How the candidate was sourced",
				Values: []string{"linkedin", "referral", "job_board", "direct_apply", "recruiter_outreach", "career_fair", "github", "stack_overflow"}},
			{Name: "tags", Type: "array", Description: "Tags like remote_ok, visa_required, urgent, etc."},
		},
	}
}

func GetRecruitersSchema() CollectionSchema {
	return CollectionSchema{
		Name:        "recruiters",
		Description: "Recruiter team members",
		Fields: []FieldInfo{
			{Name: "id", Type: "string", Description: "Recruiter ID (rec_001 through rec_008)"},
			{Name: "name", Type: "string", Description: "Recruiter full name"},
			{Name: "email", Type: "string", Description: "Recruiter email"},
			{Name: "phone", Type: "string", Description: "Recruiter phone"},
			{Name: "department", Type: "string", Description: "Department name"},
			{Name: "hire_date", Type: "string", Description: "Date recruiter was hired"},
		},
	}
}
