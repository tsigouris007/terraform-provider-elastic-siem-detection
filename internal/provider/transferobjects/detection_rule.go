package transferobjects

import "time"

type ThreatItem struct {
	Framework string `json:"framework,omitempty"`
	Tactic    struct {
		ID        string `json:"id,omitempty"`
		Name      string `json:"name,omitempty"`
		Reference string `json:"reference,omitempty"`
	} `json:"tactic,omitempty"`
	Technique []struct {
		ID           string `json:"id,omitempty"`
		Name         string `json:"name,omitempty"`
		Reference    string `json:"reference,omitempty"`
		Subtechnique []struct {
			ID        string `json:"id,omitempty"`
			Name      string `json:"name,omitempty"`
			Reference string `json:"reference,omitempty"`
		} `json:"subtechnique,omitempty"`
	} `json:"technique,omitempty"`
}

type ThreatMapping struct {
	Entries []struct {
		Field string `json:"field,omitempty"`
		Type  string `json:"type,omitempty"`
		Value string `json:"value,omitempty"`
	} `json:"entries,omitempty"`
}

type RuleThreshold struct {
	Cardinality []struct {
		Field string `json:"field,omitempty"`
		Value int    `json:"value,omitempty"`
	} `json:"cardinality,omitempty"`
	Field []string `json:"field,omitempty"`
	Value int      `json:"value,omitempty"`
}

type ExecutionHistoryItem struct {
	LastExecution struct {
		Date        time.Time `json:"date,omitempty"`
		Status      string    `json:"status,omitempty"`
		StatusOrder int       `json:"status_order,omitempty"`
		Message     string    `json:"message,omitempty"`
		Metrics     struct {
			TotalSearchDurationMs     int `json:"total_search_duration_ms,omitempty"`
			TotalIndexingDurationMs   int `json:"total_indexing_duration_ms,omitempty"`
			TotalEnrichmentDurationMs int `json:"total_enrichment_duration_ms,omitempty"`
		} `json:"metrics,omitempty"`
	} `json:"last_execution,omitempty"`
}

type RiskScoreMapping struct {
	Field    string `json:"field,omitempty"`
	Operator string `json:"operator,omitempty"`
	Value    string `json:"value,omitempty"`
}

type SeverityMapping struct {
	Field    string `json:"field,omitempty"`
	Value    string `json:"value,omitempty"`
	Operator string `json:"operator,omitempty"`
	Severity string `json:"severity,omitempty"`
}

type ActionItem struct {
	ActionTypeID string `json:"action_type_id,omitempty"`
	Group        string `json:"group,omitempty"`
	ID           string `json:"id,omitempty"`
	Params       struct {
		Documents   []string `json:"documents,omitempty"`
		Message     string   `json:"message,omitempty"`
		To          string   `json:"to,omitempty"`
		Cc          string   `json:"cc,omitempty"`
		Bcc         string   `json:"bcc,omitempty"`
		Subject     string   `json:"subject,omitempty"`
		Body        string   `json:"body,omitempty"`
		Severity    string   `json:"severity,omitempty"`
		EventAction string   `json:"eventAction,omitempty"`
		DedupKey    string   `json:"dedupKey,omitempty"`
		Timestamp   string   `json:"timestamp,omitempty"`
		Component   string   `json:"component,omitempty"`
		Group       string   `json:"group,omitempty"`
		Source      string   `json:"source,omitempty"`
		Summary     string   `json:"summary,omitempty"`
		Class       string   `json:"class,omitempty"`
	} `json:"params,omitempty"`
}

type MetaItem struct {
	From             string `json:"from,omitempty"`
	KibanaSiemAppURL string `json:"kibana_siem_app_url,omitempty"`
}

type ExceptionListItem struct {
	ID            string `json:"id,omitempty"`
	ListID        string `json:"list_id,omitempty"`
	NamespaceType string `json:"namespace_type,omitempty"`
	Type          string `json:"type,omitempty"`
}

type DetectionRuleResponse struct {
	DetectionRule
	CreatedAt        time.Time            `json:"created_at,omitempty"`
	CreatedBy        string               `json:"created_by,omitempty"`
	ExecutionSummary ExecutionHistoryItem `json:"execution_summary,omitempty"`
	Meta             MetaItem             `json:"meta,omitempty"`
	UpdatedAt        time.Time            `json:"updated_at,omitempty"`
}

type DetectionRule struct {
	Actions             []ActionItem        `json:"actions,omitempty"`
	AnomalyThreshold    int                 `json:"anomaly_threshold,omitempty"`
	Author              []string            `json:"author,omitempty"`
	BuildingBlockTYpe   string              `json:"building_block_type,omitempty"`
	Description         string              `json:"description,omitempty"`
	Enabled             *bool               `json:"enabled,omitempty"` // bool values need to be pointers to include false
	EventCategoryField  string              `json:"event_category_field,omitempty"`
	ExceptionsList      []ExceptionListItem `json:"exceptions_list"`
	FalsePositives      []interface{}       `json:"false_positives,omitempty"`
	Filters             []interface{}       `json:"filters,omitempty"`
	From                string              `json:"from,omitempty"`
	ID                  string              `json:"id,omitempty"`
	Immutable           *bool               `json:"immutable,omitempty"` // bool values need to be pointers to include false
	Index               []string            `json:"index,omitempty"`
	Interval            string              `json:"interval,omitempty"`
	Language            string              `json:"language,omitempty"`
	License             string              `json:"license,omitempty"`
	MachineLeanJID      []string            `json:"machine_learning_job_id,omitempty"`
	MaxSignals          int                 `json:"max_signals,omitempty"`
	Name                string              `json:"name,omitempty"`
	Note                string              `json:"note,omitempty"`
	OutputIndex         string              `json:"output_index,omitempty"`
	Query               string              `json:"query,omitempty"`
	References          []interface{}       `json:"references,omitempty"`
	RelatedIntegrations []interface{}       `json:"related_integrations,omitempty"`
	RequiredFields      []interface{}       `json:"required_fields,omitempty"`
	RiskScore           int                 `json:"risk_score,omitempty"`
	RiskScoreMapping    []RiskScoreMapping  `json:"risk_score_mapping,omitempty"`
	RuleID              string              `json:"rule_id,omitempty"`
	RuleNameOverride    string              `json:"rule_name_override,omitempty"`
	SaveId              string              `json:"setup,omitempty"`
	Setup               string              `json:"saved_id,omitempty"`
	Severity            string              `json:"severity,omitempty"`
	SeverityMapping     []SeverityMapping   `json:"severity_mapping,omitempty"`
	Tags                []string            `json:"tags,omitempty"`
	Threat              []ThreatItem        `json:"threat,omitempty"`
	ThreatFilters       []interface{}       `json:"threat_filters,omitempty"`
	ThreatIndex         []string            `json:"threat_index,omitempty"`
	ThreatIndicatorPath string              `json:"threat_indicator_path,omitempty"`
	ThreatQuery         string              `json:"threat_query,omitempty"`
	ThreatMapping       []ThreatMapping     `json:"threat_mapping,omitempty"`
	Threshold           RuleThreshold       `json:"threshold,omitempty"`
	Throttle            string              `json:"throttle,omitempty"`
	TiebreakerField     string              `json:"tiebreaker_field,omitempty"`
	TimestampField      string              `json:"timestamp_field,omitempty"`
	TimeStampOverride   string              `json:"timestamp_override,omitempty"`
	To                  string              `json:"to,omitempty"`
	Type                string              `json:"type,omitempty"`
	UpdatedBy           string              `json:"updated_by,omitempty"`
	Version             int                 `json:"version,omitempty"`
}
