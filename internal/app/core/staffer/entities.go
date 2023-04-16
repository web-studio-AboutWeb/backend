package staffer

type StafferPosition string

const (
	StafferPositionFrontendDev = "frontend"
	StafferPositionBackendDev  = "backend"
	StafferPositionTeamLead    = "teamlead"
	StafferPositionManager     = "manager"
	StafferPositionMarketer    = "marketer"
)

type Staffer struct {
	Id        int16           `json:"id"`
	UserId    int16           `json:"user_id"`
	ProjectId int16           `json:"project_id"`
	Position  StafferPosition `json:"position"`
}
