package model

import (
	"errors"
	"strings"
)

// Rule for filter
type Rule struct {
	From     string
	To       string
	Subject  string
	Channels []string
}

// Filter has an array of rules
type Filter []*Rule

// Fix remove spaces and convert to lower case the rules
func (r *Rule) Fix() {
	r.From = strings.TrimSpace(strings.ToLower(r.From))
	r.To = strings.TrimSpace(strings.ToLower(r.To))
	r.Subject = strings.TrimSpace(strings.ToLower(r.Subject))

	for i, channel := range r.Channels {
		channel = strings.TrimSpace(channel)
		channel = strings.ToLower(channel)

		if !strings.HasPrefix(channel, "#") && !strings.HasPrefix(channel, "@") {
			channel = "#" + channel
		}
		r.Channels[i] = channel
	}
}

// Validate check if this rule is valid
func (r *Rule) Validate() error {
	if len(r.From) == 0 && len(r.Subject) == 0 && len(r.To) == 0 {
		return errors.New("Need to set From, or To, or Subject")
	}

	if len(r.Channels) == 0 {
		return errors.New("Need to set at least one channel or user for destination")
	}

	for _, channel := range r.Channels {
		if channel != "" && !validateChannel(channel) {
			return errors.New("Need to set #channel or @user")
		}
	}

	return nil
}

func (r *Rule) matchFrom(from string) bool {
	from = strings.ToLower(from)
	if len(r.From) == 0 {
		return true
	}
	return strings.Contains(from, r.From)
}

func (r *Rule) matchTo(to string) bool {
	to = strings.ToLower(to)
	if len(r.To) == 0 {
		return true
	}
	return strings.Contains(to, r.To)
}

func (r *Rule) matchSubject(subject string) bool {
	subject = strings.ToLower(subject)
	if len(r.Subject) == 0 {
		return true
	}
	return strings.Contains(subject, r.Subject)
}

// Match check if from and subject meets this rule
func (r *Rule) Match(from, to, subject string) bool {
	return r.matchFrom(from) && r.matchTo(to) && r.matchSubject(subject)
}

// GetChannels return the first channels with attempt the rules
func (f *Filter) GetChannels(from, subject string) []string {
	for _, r := range *f {
		if r.Match(from, to, subject) {
			return r.Channels
		}
	}
	return []string{""}
}

// Validate check if all rules is valid
func (f *Filter) Validate() error {
	if len(*f) == 0 {
		return errors.New("Filter need to be at least one rule to be valid")
	}

	for _, r := range *f {
		if err := r.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Fix all rules
func (f *Filter) Fix() {
	for _, r := range *f {
		r.Fix()
	}
}
