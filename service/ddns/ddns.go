// Package ddns keeps a Route 53 A record pointed at the machine's current
// public IP address. It is meant for hosts that sit behind a dynamic public IP
// (e.g. a home connection) so that a stable hostname always resolves to the box.
package ddns

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	log "github.com/sirupsen/logrus"
)

// ipEcho endpoints return the caller's public IP as a plain text body. We query
// them in order and use the first that answers so a single provider outage does
// not stall the updater.
var ipEchoURLs = []string{
	"https://checkip.amazonaws.com",
	"https://api.ipify.org",
	"https://ifconfig.me/ip",
}

// Config holds everything the updater needs. The zero value is not usable;
// HostedZoneID and RecordName are required.
type Config struct {
	// HostedZoneID is the Route 53 hosted zone that owns RecordName.
	HostedZoneID string
	// RecordName is the fully qualified record to keep updated, e.g.
	// "kesmarki.godevltd.com".
	RecordName string
	// TTL is the A record TTL in seconds. Defaults to 300 when zero.
	TTL int64
	// Interval is how often the public IP is re-checked. Defaults to 1h when
	// zero.
	Interval time.Duration
}

// Updater synchronises a Route 53 A record with the host's public IP.
type Updater struct {
	cfg    Config
	client *route53.Client
	http   *http.Client

	// lastIP caches the IP we last pushed to Route 53 so we skip the API call
	// when nothing changed. It is only touched from Run's goroutine.
	lastIP string
}

// New builds an Updater. It loads AWS credentials/region from the default chain
// (environment variables, shared config, instance profile, ...).
func New(ctx context.Context, cfg Config) (*Updater, error) {
	if cfg.HostedZoneID == "" || cfg.RecordName == "" {
		return nil, fmt.Errorf("ddns: HostedZoneID and RecordName are required")
	}
	if cfg.TTL == 0 {
		cfg.TTL = 300
	}
	if cfg.Interval == 0 {
		cfg.Interval = time.Hour
	}

	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("ddns: load aws config: %w", err)
	}

	return &Updater{
		cfg:    cfg,
		client: route53.NewFromConfig(awsCfg),
		http:   &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Run performs an initial sync, then re-checks on cfg.Interval until ctx is
// cancelled. It blocks, so callers typically run it in its own goroutine.
func (u *Updater) Run(ctx context.Context) {
	log.Infof("ddns: keeping %s in sync every %s", u.cfg.RecordName, u.cfg.Interval)

	if err := u.syncOnce(ctx); err != nil {
		log.Errorf("ddns: initial sync failed: %s", err)
	}

	ticker := time.NewTicker(u.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("ddns: stopped")
			return
		case <-ticker.C:
			if err := u.syncOnce(ctx); err != nil {
				log.Errorf("ddns: sync failed: %s", err)
			}
		}
	}
}

// syncOnce fetches the current public IP and updates Route 53 if it differs
// from the last value we pushed.
func (u *Updater) syncOnce(ctx context.Context) error {
	ip, err := u.publicIP(ctx)
	if err != nil {
		return err
	}

	if ip == u.lastIP {
		log.Debugf("ddns: public IP unchanged (%s), skipping update", ip)
		return nil
	}

	if err := u.upsert(ctx, ip); err != nil {
		return err
	}

	u.lastIP = ip
	log.Infof("ddns: %s -> %s", u.cfg.RecordName, ip)
	return nil
}

// publicIP returns the host's current public IPv4 address, trying each echo
// endpoint until one responds.
func (u *Updater) publicIP(ctx context.Context) (string, error) {
	var lastErr error
	for _, url := range ipEchoURLs {
		ip, err := u.fetchIP(ctx, url)
		if err != nil {
			lastErr = err
			log.Debugf("ddns: %s failed: %s", url, err)
			continue
		}
		return ip, nil
	}
	return "", fmt.Errorf("ddns: could not determine public IP: %w", lastErr)
}

func (u *Updater) fetchIP(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := u.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 64))
	if err != nil {
		return "", err
	}

	ip := strings.TrimSpace(string(body))
	if ip == "" {
		return "", fmt.Errorf("empty response")
	}
	return ip, nil
}

// upsert points the A record at ip via a Route 53 UPSERT change.
func (u *Updater) upsert(ctx context.Context, ip string) error {
	_, err := u.client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(u.cfg.HostedZoneID),
		ChangeBatch: &types.ChangeBatch{
			Comment: aws.String("kesmarki ddns"),
			Changes: []types.Change{{
				Action: types.ChangeActionUpsert,
				ResourceRecordSet: &types.ResourceRecordSet{
					Name: aws.String(u.cfg.RecordName),
					Type: types.RRTypeA,
					TTL:  aws.Int64(u.cfg.TTL),
					ResourceRecords: []types.ResourceRecord{{
						Value: aws.String(ip),
					}},
				},
			}},
		},
	})
	if err != nil {
		return fmt.Errorf("ddns: route53 upsert: %w", err)
	}
	return nil
}
