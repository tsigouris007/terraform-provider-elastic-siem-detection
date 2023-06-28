package transferobjects

type PrivilegesResponse struct {
	Username        string `json:"username,omitempty"`
	HasAllRequested bool   `json:"has_all_requested,omitempty"`
	Cluster         struct {
		MonitorMl             bool `json:"monitor_ml,omitempty"`
		ManageCcr             bool `json:"manage_ccr,omitempty"`
		ManageIndexTemplates  bool `json:"manage_index_templates,omitempty"`
		MonitorWatcher        bool `json:"monitor_watcher,omitempty"`
		MonitorTransform      bool `json:"monitor_transform,omitempty"`
		ReadIlm               bool `json:"read_ilm,omitempty"`
		ManageAPIKey          bool `json:"manage_api_key,omitempty"`
		ManageSecurity        bool `json:"manage_security,omitempty"`
		ManageOwnAPIKey       bool `json:"manage_own_api_key,omitempty"`
		ManageSaml            bool `json:"manage_saml,omitempty"`
		All                   bool `json:"all,omitempty"`
		ManageIlm             bool `json:"manage_ilm,omitempty"`
		ManageIngestPipelines bool `json:"manage_ingest_pipelines,omitempty"`
		ReadCcr               bool `json:"read_ccr,omitempty"`
		ManageRollup          bool `json:"manage_rollup,omitempty"`
		Monitor               bool `json:"monitor,omitempty"`
		ManageWatcher         bool `json:"manage_watcher,omitempty"`
		Manage                bool `json:"manage,omitempty"`
		ManageTransform       bool `json:"manage_transform,omitempty"`
		ManageToken           bool `json:"manage_token,omitempty"`
		ManageMl              bool `json:"manage_ml,omitempty"`
		ManagePipeline        bool `json:"manage_pipeline,omitempty"`
		MonitorRollup         bool `json:"monitor_rollup,omitempty"`
		TransportClient       bool `json:"transport_client,omitempty"`
		CreateSnapshot        bool `json:"create_snapshot,omitempty"`
	} `json:"cluster,omitempty"`
	Index struct {
		AlertsSecurityAlertsDefault struct {
			All               bool `json:"all,omitempty"`
			Create            bool `json:"create,omitempty"`
			CreateDoc         bool `json:"create_doc,omitempty"`
			CreateIndex       bool `json:"create_index,omitempty"`
			Delete            bool `json:"delete,omitempty"`
			DeleteIndex       bool `json:"delete_index,omitempty"`
			Index             bool `json:"index,omitempty"`
			Maintenance       bool `json:"maintenance,omitempty"`
			Manage            bool `json:"manage,omitempty"`
			ManageFollowIndex bool `json:"manage_follow_index,omitempty"`
			ManageIlm         bool `json:"manage_ilm,omitempty"`
			ManageLeaderIndex bool `json:"manage_leader_index,omitempty"`
			Monitor           bool `json:"monitor,omitempty"`
			Read              bool `json:"read,omitempty"`
			ReadCrossCluster  bool `json:"read_cross_cluster,omitempty"`
			ViewIndexMetadata bool `json:"view_index_metadata,omitempty"`
			Write             bool `json:"write,omitempty"`
		} `json:".alerts-security.alerts-default,omitempty"`
	} `json:"index,omitempty"`
	Application struct {
	} `json:"application,omitempty"`
	IsAuthenticated  bool `json:"is_authenticated,omitempty"`
	HasEncryptionKey bool `json:"has_encryption_key,omitempty"`
}
