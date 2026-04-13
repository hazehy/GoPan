package types

type AdminOverviewResponse struct {
	TotalUsers     int64           `json:"total_users"`
	ActiveUsers    int64           `json:"active_users"`
	DisabledUsers  int64           `json:"disabled_users"`
	TotalFiles     int64           `json:"total_files"`
	TotalFolders   int64           `json:"total_folders"`
	TotalFileSize  int64           `json:"total_file_size"`
	TodayUploads   int64           `json:"today_uploads"`
	TodayRegisters int64           `json:"today_registers"`
	ExtStats       []*AdminExtStat `json:"ext_stats"`
}

type AdminExtStat struct {
	Ext   string `json:"ext"`
	Count int64  `json:"count"`
}

type AdminUserListRequest struct {
	Page    int    `json:"page,optional" form:"page,optional"`
	Size    int    `json:"size,optional" form:"size,optional"`
	Keyword string `json:"keyword,optional" form:"keyword,optional"`
}

type AdminUserItem struct {
	Identity           string `json:"identity"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	Status             int    `json:"status"`
	Role               int    `json:"role"`
	UploadPermission   int    `json:"upload_permission"`
	DownloadPermission int    `json:"download_permission"`
	SharePermission    int    `json:"share_permission"`
	LastLoginAt        string `json:"last_login_at"`
	CreatedAt          string `json:"created_at"`
}

type AdminUserListResponse struct {
	List  []*AdminUserItem `json:"list"`
	Count int64            `json:"count"`
}

type AdminUserStatusUpdateRequest struct {
	Identity           string `json:"identity"`
	Status             *int   `json:"status,optional" form:"status,optional"`
	UploadPermission   *int   `json:"upload_permission,optional" form:"upload_permission,optional"`
	DownloadPermission *int   `json:"download_permission,optional" form:"download_permission,optional"`
	SharePermission    *int   `json:"share_permission,optional" form:"share_permission,optional"`
}

type AdminFileListRequest struct {
	Page     int    `json:"page,optional" form:"page,optional"`
	Size     int    `json:"size,optional" form:"size,optional"`
	Keyword  string `json:"keyword,optional" form:"keyword,optional"`
	UserName string `json:"user_name,optional" form:"user_name,optional"`
}

type AdminFileItem struct {
	Identity           string `json:"identity"`
	ParentId           int64  `json:"parent_id"`
	UserIdentity       string `json:"user_identity"`
	UserName           string `json:"user_name"`
	RepositoryIdentity string `json:"repository_identity"`
	Name               string `json:"name"`
	Ext                string `json:"ext"`
	Path               string `json:"path"`
	Size               int64  `json:"size"`
	UpdatedAt          string `json:"updated_at"`
}

type AdminFileListResponse struct {
	List  []*AdminFileItem `json:"list"`
	Count int64            `json:"count"`
}

type AdminFileDeleteRequest struct {
	Identity string `json:"identity"`
}

type AdminLogListRequest struct {
	Page      int    `json:"page,optional" form:"page,optional"`
	Size      int    `json:"size,optional" form:"size,optional"`
	Keyword   string `json:"keyword,optional" form:"keyword,optional"`
	Action    string `json:"action,optional" form:"action,optional"`
	ActorName string `json:"actor_name,optional" form:"actor_name,optional"`
	FileExt   string `json:"file_ext,optional" form:"file_ext,optional"`
	Sharer    string `json:"sharer_name,optional" form:"sharer_name,optional"`
	Saver     string `json:"saver_name,optional" form:"saver_name,optional"`
	Day       string `json:"day,optional" form:"day,optional"`
}

type AdminLogItem struct {
	Identity       string `json:"identity"`
	ActorIdentity  string `json:"actor_identity"`
	ActorName      string `json:"actor_name"`
	ActorRole      int    `json:"actor_role"`
	Action         string `json:"action"`
	TargetType     string `json:"target_type"`
	TargetIdentity string `json:"target_identity"`
	Detail         string `json:"detail"`
	CreatedAt      string `json:"created_at"`
}

type AdminLogListResponse struct {
	List  []*AdminLogItem `json:"list"`
	Count int64           `json:"count"`
}
