package exercise

type Summary struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	IsCompound    bool   `json:"is_compound"`
	PrimaryMuscle string `json:"primary_muscle"`
}
