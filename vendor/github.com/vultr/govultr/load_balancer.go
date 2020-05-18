package govultr

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// LoadBalancerService is the interface to interact with the server endpoints on the Vultr API
// Link: https://www.vultr.com/api/#loadbalancer
type LoadBalancerService interface {
	List(ctx context.Context) ([]LoadBalancers, error)
	Delete(ctx context.Context, ID int) error
	SetLabel(ctx context.Context, ID int, label string) error
	AttachedInstances(ctx context.Context, ID int) (*InstanceList, error)
	AttachInstance(ctx context.Context, ID, backendNode int) error
	DetachInstance(ctx context.Context, ID, backendNode int) error
	GetHealthCheck(ctx context.Context, ID int) (*HealthCheck, error)
	SetHealthCheck(ctx context.Context, ID int, healthConfig *HealthCheck) error
	GetGenericInfo(ctx context.Context, ID int) (*GenericInfo, error)
	ListForwardingRules(ctx context.Context, ID int) (*ForwardingRules, error)
	DeleteForwardingRule(ctx context.Context, ID int, RuleID string) error
	CreateForwardingRule(ctx context.Context, ID int, rule *ForwardingRule) (*ForwardingRule, error)
	GetFullConfig(ctx context.Context, ID int) (*LBConfig, error)
	HasSSL(ctx context.Context, ID int) (*struct {
		SSLInfo bool `json:"has_ssl"`
	}, error)
	Create(ctx context.Context, region int, label string, genericInfo *GenericInfo, healthCheck *HealthCheck, rules []ForwardingRule, ssl *SSL, instances *InstanceList) (*LoadBalancers, error)
	UpdateGenericInfo(ctx context.Context, ID int, label string, genericInfo *GenericInfo) error
	AddSSL(ctx context.Context, ID int, ssl *SSL) error
	RemoveSSL(ctx context.Context, ID int) error
}

// LoadBalancerHandler handles interaction with the server methods for the Vultr API
type LoadBalancerHandler struct {
	client *Client
}

// LoadBalancers represent a basic structure of a load balancer
type LoadBalancers struct {
	ID          int    `json:"SUBID,omitempty"`
	DateCreated string `json:"date_created,omitempty"`
	RegionID    int    `json:"DCID,omitempty"`
	Location    string `json:"location,omitempty"`
	Label       string `json:"label,omitempty"`
	Status      string `json:"status,omitempty"`
	IPV4        string `json:"ipv4,omitempty"`
	IPV6        string `json:"ipv6,omitempty"`
}

// InstanceList represents instances that attached to your load balancer
type InstanceList struct {
	InstanceList []int `json:"instance_list"`
}

// HealthCheck represents your health check configuration for your load balancer.
type HealthCheck struct {
	Protocol           string `json:"protocol,omitempty"`
	Port               int    `json:"port,omitempty"`
	Path               string `json:"path,omitempty"`
	CheckInterval      int    `json:"check_interval,omitempty"`
	ResponseTimeout    int    `json:"response_timeout,omitempty"`
	UnhealthyThreshold int    `json:"unhealthy_threshold,omitempty"`
	HealthyThreshold   int    `json:"healthy_threshold,omitempty"`
}

// GenericInfo represents generic configuration of your load balancer
type GenericInfo struct {
	BalancingAlgorithm string          `json:"balancing_algorithm"`
	SSLRedirect        *bool           `json:"ssl_redirect,omitempty"`
	StickySessions     *StickySessions `json:"sticky_sessions"`
	ProxyProtocol      *bool           `json:"proxy_protocol"`
}

// CookieName represents cookie for your load balancer
type StickySessions struct {
	StickySessionsEnabled string `json:"sticky_sessions"`
	CookieName            string `json:"cookie_name"`
}

// ForwardingRules represent a list of forwarding rules
type ForwardingRules struct {
	ForwardRuleList []ForwardingRule `json:"forward_rule_list"`
}

// ForwardingRule represent a single forwarding rule
type ForwardingRule struct {
	RuleID           string `json:"RULEID,omitempty"`
	FrontendProtocol string `json:"frontend_protocol,omitempty"`
	FrontendPort     int    `json:"frontend_port,omitempty"`
	BackendProtocol  string `json:"backend_protocol,omitempty"`
	BackendPort      int    `json:"backend_port,omitempty"`
}

// LBConfig represents the full config with all components of a load balancer
type LBConfig struct {
	GenericInfo `json:"generic_info"`
	HealthCheck `json:"health_checks_info"`
	SSLInfo     bool `json:"has_ssl"`
	ForwardingRules
	InstanceList
}

// SSL represents valid SSL config
type SSL struct {
	PrivateKey  string `json:"ssl_private_key"`
	Certificate string `json:"ssl_certificate"`
	Chain       string `json:"chain,omitempty"`
}

// List all load balancer subscriptions on the current account.
func (l *LoadBalancerHandler) List(ctx context.Context) ([]LoadBalancers, error) {
	uri := "/v1/loadbalancer/list"

	req, err := l.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	var lbList []LoadBalancers

	err = l.client.DoWithContext(ctx, req, &lbList)
	if err != nil {
		return nil, err
	}

	return lbList, nil
}

// Delete a load balancer subscription.
func (l *LoadBalancerHandler) Delete(ctx context.Context, ID int) error {
	uri := "/v1/loadbalancer/destroy"

	values := url.Values{
		"SUBID": {strconv.Itoa(ID)},
	}

	req, err := l.client.NewRequest(ctx, http.MethodPost, uri, values)
	if err != nil {
		return err
	}

	err = l.client.DoWithContext(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

// SetLabel sets the label for your load balancer subscription.
func (l *LoadBalancerHandler) SetLabel(ctx context.Context, ID int, label string) error {
	uri := "/v1/loadbalancer/label_set"

	values := url.Values{
		"SUBID": {strconv.Itoa(ID)},
		"label": {label},
	}

	req, err := l.client.NewRequest(ctx, http.MethodPost, uri, values)
	if err != nil {
		return err
	}

	err = l.client.DoWithContext(ctx, req, nil)
	if err != nil {
		return nil
	}

	return nil
}

// AttachedInstances lists the instances that are currently attached to a load balancer subscription.
func (l *LoadBalancerHandler) AttachedInstances(ctx context.Context, ID int) (*InstanceList, error) {
	uri := "/v1/loadbalancer/instance_list"

	req, err := l.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("SUBID", strconv.Itoa(ID))
	req.URL.RawQuery = q.Encode()

	var instances InstanceList

	err = l.client.DoWithContext(ctx, req, &instances)
	if err != nil {
		return nil, err
	}

	return &instances, nil
}

// AttachInstance attaches a backend node to your load balancer subscription
func (l *LoadBalancerHandler) AttachInstance(ctx context.Context, ID, backendNode int) error {
	uri := "/v1/loadbalancer/instance_attach"

	values := url.Values{
		"SUBID":       {strconv.Itoa(ID)},
		"backendNode": {strconv.Itoa(backendNode)},
	}

	req, err := l.client.NewRequest(ctx, http.MethodPost, uri, values)
	if err != nil {
		return err
	}

	err = l.client.DoWithContext(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

// DetachInstance detaches a backend node to your load balancer subscription
func (l *LoadBalancerHandler) DetachInstance(ctx context.Context, ID, backendNode int) error {
	uri := "/v1/loadbalancer/instance_detach"

	values := url.Values{
		"SUBID":       {strconv.Itoa(ID)},
		"backendNode": {strconv.Itoa(backendNode)},
	}

	req, err := l.client.NewRequest(ctx, http.MethodPost, uri, values)
	if err != nil {
		return err
	}

	err = l.client.DoWithContext(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetHealthCheck retrieves the health check configuration for your load balancer subscription.
func (l *LoadBalancerHandler) GetHealthCheck(ctx context.Context, ID int) (*HealthCheck, error) {
	uri := "/v1/loadbalancer/health_check_info"

	req, err := l.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("SUBID", strconv.Itoa(ID))
	req.URL.RawQuery = q.Encode()

	var healthCheck HealthCheck
	err = l.client.DoWithContext(ctx, req, &healthCheck)
	if err != nil {
		return nil, err
	}

	return &healthCheck, nil
}

// SetHealthCheck sets your health check configuration for your load balancer
func (l *LoadBalancerHandler) SetHealthCheck(ctx context.Context, ID int, healthConfig *HealthCheck) error {
	uri := "/v1/loadbalancer/health_check_update"

	values := url.Values{
		"SUBID": {strconv.Itoa(ID)},
	}

	if healthConfig != nil {
		if healthConfig.Protocol != "" {
			values.Add("protocol", healthConfig.Protocol)
		}

		if healthConfig.Port != 0 {
			values.Add("port", strconv.Itoa(healthConfig.Port))
		}

		if healthConfig.CheckInterval != 0 {
			values.Add("check_interval", strconv.Itoa(healthConfig.CheckInterval))
		}

		if healthConfig.ResponseTimeout != 0 {
			values.Add("response_timeout", strconv.Itoa(healthConfig.ResponseTimeout))
		}

		if healthConfig.UnhealthyThreshold != 0 {
			values.Add("unhealthy_threshold", strconv.Itoa(healthConfig.UnhealthyThreshold))
		}

		if healthConfig.HealthyThreshold != 0 {
			values.Add("healthy_threshold", strconv.Itoa(healthConfig.HealthyThreshold))
		}

		if healthConfig.Path != "" {
			values.Add("path", healthConfig.Path)
		}
	}

	req, err := l.client.NewRequest(ctx, http.MethodPost, uri, values)
	if err != nil {
		return err
	}

	err = l.client.DoWithContext(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetGenericInfo is the generic configuration of a load balancer subscription
func (l *LoadBalancerHandler) GetGenericInfo(ctx context.Context, ID int) (*GenericInfo, error) {
	uri := "/v1/loadbalancer/generic_info"

	req, err := l.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("SUBID", strconv.Itoa(ID))
	req.URL.RawQuery = q.Encode()

	var info GenericInfo

	err = l.client.DoWithContext(ctx, req, &info)
	if err != nil {
		return nil, err
	}

	return &info, err
}

// ListForwardingRules lists all forwarding rules for a load balancer subscription
func (l *LoadBalancerHandler) ListForwardingRules(ctx context.Context, ID int) (*ForwardingRules, error) {
	uri := "v1/loadbalancer/forward_rule_list"

	req, err := l.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("SUBID", strconv.Itoa(ID))
	req.URL.RawQuery = q.Encode()

	var frList ForwardingRules

	err = l.client.DoWithContext(ctx, req, &frList)
	if err != nil {
		return nil, err
	}

	return &frList, nil
}

// DeleteForwardingRule removes a forwarding rule from a load balancer subscription
func (l *LoadBalancerHandler) DeleteForwardingRule(ctx context.Context, ID int, RuleID string) error {
	uri := "/v1/loadbalancer/forward_rule_delete"

	values := url.Values{
		"SUBID":  {strconv.Itoa(ID)},
		"RULEID": {RuleID},
	}

	req, err := l.client.NewRequest(ctx, http.MethodPost, uri, values)
	if err != nil {
		return err
	}

	err = l.client.DoWithContext(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

// CreateForwardingRule will create a new forwarding rule for your load balancer subscription.
// Note the RuleID will be returned in the ForwardingRule struct
func (l *LoadBalancerHandler) CreateForwardingRule(ctx context.Context, ID int, rule *ForwardingRule) (*ForwardingRule, error) {
	uri := "/v1/loadbalancer/forward_rule_create"

	values := url.Values{
		"SUBID":             {strconv.Itoa(ID)},
		"frontend_protocol": {rule.FrontendProtocol},
		"backend_protocol":  {rule.BackendProtocol},
		"frontend_port":     {strconv.Itoa(rule.FrontendPort)},
		"backend_port":      {strconv.Itoa(rule.BackendPort)},
	}

	req, err := l.client.NewRequest(ctx, http.MethodPost, uri, values)
	if err != nil {
		return nil, err
	}

	var fr ForwardingRule
	err = l.client.DoWithContext(ctx, req, &fr)
	if err != nil {
		return nil, err
	}

	return &fr, nil
}

// GetFullConfig retrieves the entire configuration of a load balancer subscription.
func (l *LoadBalancerHandler) GetFullConfig(ctx context.Context, ID int) (*LBConfig, error) {
	uri := "/v1/loadbalancer/conf_info"

	req, err := l.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("SUBID", strconv.Itoa(ID))
	req.URL.RawQuery = q.Encode()

	var lbConfig LBConfig
	err = l.client.DoWithContext(ctx, req, &lbConfig)
	if err != nil {
		return nil, err
	}

	return &lbConfig, nil
}

// HasSSL retrieves whether or not your load balancer subscription has an SSL cert attached.
func (l *LoadBalancerHandler) HasSSL(ctx context.Context, ID int) (*struct {
	SSLInfo bool `json:"has_ssl"`
}, error) {
	uri := "/v1/loadbalancer/ssl_info"

	req, err := l.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("SUBID", strconv.Itoa(ID))
	req.URL.RawQuery = q.Encode()

	ssl := &struct {
		SSLInfo bool `json:"has_ssl"`
	}{}
	err = l.client.DoWithContext(ctx, req, ssl)
	if err != nil {
		return nil, err
	}

	return ssl, nil
}

// Create a load balancer
func (l *LoadBalancerHandler) Create(ctx context.Context, region int, label string, genericInfo *GenericInfo, healthCheck *HealthCheck, rules []ForwardingRule, ssl *SSL, instances *InstanceList) (*LoadBalancers, error) {
	uri := "/v1/loadbalancer/create"

	values := url.Values{
		"DCID": {strconv.Itoa(region)},
	}

	if label != "" {
		values.Add("label", label)
	}

	// Check generic info struct
	if genericInfo != nil {
		if genericInfo.SSLRedirect != nil {
			if strconv.FormatBool(*genericInfo.SSLRedirect) == "true" {
				values.Add("config_ssl_redirect", "true")
			}
		}

		if genericInfo.BalancingAlgorithm != "" {
			values.Add("balancing_algorithm", genericInfo.BalancingAlgorithm)
		}

		if genericInfo.StickySessions != nil {
			if genericInfo.StickySessions.StickySessionsEnabled == "on" {
				values.Add("sticky_sessions", genericInfo.StickySessions.StickySessionsEnabled)
				values.Add("cookie_name", genericInfo.StickySessions.CookieName)
			}
		}

		if genericInfo.ProxyProtocol != nil {
			value := "off"
			if strconv.FormatBool(*genericInfo.ProxyProtocol) == "true" {
				value = "on"
			}
			values.Add("proxy_protocol", value)
		}
	}

	if healthCheck != nil {
		t, _ := json.Marshal(healthCheck)
		values.Add("health_check", string(t))
	}

	if rules != nil {
		t, e := json.Marshal(rules)
		if e != nil {
			panic(e)
		}
		values.Add("forwarding_rules", string(t))
	}

	if ssl != nil {
		values.Add("ssl_private_key", ssl.PrivateKey)
		values.Add("ssl_certificate", ssl.Certificate)

		if ssl.Chain != "" {
			values.Add("ssl_chain", ssl.Chain)
		}
	}

	if instances != nil {
		t, e := json.Marshal(instances.InstanceList)
		if e != nil {
			panic(e)
		}
		values.Add("attached_nodes", string(t))
	}

	req, err := l.client.NewRequest(ctx, http.MethodPost, uri, values)
	if err != nil {
		return nil, err
	}

	var lb LoadBalancers
	err = l.client.DoWithContext(ctx, req, &lb)
	if err != nil {
		return nil, err
	}

	return &lb, nil
}

// UpdateGenericInfo will update portions of your generic info section
func (l *LoadBalancerHandler) UpdateGenericInfo(ctx context.Context, ID int, label string, genericInfo *GenericInfo) error {
	uri := "/v1/loadbalancer/generic_update"

	values := url.Values{
		"SUBID": {strconv.Itoa(ID)},
	}

	if label != "" {
		values.Add("label", label)
	}

	if genericInfo != nil {
		if genericInfo.StickySessions != nil {
			values.Add("sticky_sessions", genericInfo.StickySessions.StickySessionsEnabled)
			values.Add("cookie_name", genericInfo.StickySessions.CookieName)
		}

		if genericInfo.SSLRedirect != nil {
			values.Add("ssl_redirect", strconv.FormatBool(*genericInfo.SSLRedirect))
		}

		if genericInfo.BalancingAlgorithm != "" {
			values.Add("balancing_algorithm", genericInfo.BalancingAlgorithm)
		}

		if genericInfo.ProxyProtocol != nil {
			value := "off"
			if strconv.FormatBool(*genericInfo.ProxyProtocol) == "true" {
				value = "on"
			}
			values.Add("proxy_protocol", value)
		}
	}

	req, err := l.client.NewRequest(ctx, http.MethodPost, uri, values)
	if err != nil {
		return err
	}

	err = l.client.DoWithContext(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

// AddSSL will attach an SSL certificate to a given load balancer
func (l *LoadBalancerHandler) AddSSL(ctx context.Context, ID int, ssl *SSL) error {
	uri := "/v1/loadbalancer/ssl_add"

	values := url.Values{
		"SUBID":           {strconv.Itoa(ID)},
		"ssl_private_key": {ssl.PrivateKey},
		"ssl_certificate": {ssl.Certificate},
	}

	if ssl.Chain != "" {
		values.Add("ssl_chain", ssl.Chain)
	}

	req, err := l.client.NewRequest(ctx, http.MethodPost, uri, values)
	if err != nil {
		return err
	}

	err = l.client.DoWithContext(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

// RemoveSSL will remove an SSL certificate from a load balancer
func (l *LoadBalancerHandler) RemoveSSL(ctx context.Context, ID int) error {
	uri := "/v1/loadbalancer/ssl_remove"

	values := url.Values{
		"SUBID": {strconv.Itoa(ID)},
	}

	req, err := l.client.NewRequest(ctx, http.MethodPost, uri, values)
	if err != nil {
		return err
	}

	err = l.client.DoWithContext(ctx, req, nil)
	if err != nil {
		return err
	}
	return nil
}
