package gui

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"codeberg.org/reiver/go-activitypub"
	"codeberg.org/reiver/go-fedifinger"
)

type fetchResult struct {
	name        string
	summaryHTML string
	iconURL     string
	bannerURL   string
	profileURL  string
	err         error
}

const fetchTimeout = 15 * time.Second

func fetchActor(fediID string) fetchResult {
	ctx, cancel := context.WithTimeout(context.Background(), fetchTimeout)
	defer cancel()

	data, err := fedifinger.Get(ctx, fediID)
	if nil != err {
		return fetchResult{err: fmt.Errorf("fetching profile: %w", err)}
	}

	var actor activitypub.AnyActor
	err = json.Unmarshal(data, &actor)
	if nil != err {
		return fetchResult{err: fmt.Errorf("parsing actor: %w", err)}
	}

	var result fetchResult

	if name, ok := actor.Name.Get(); ok {
		result.name = name
	}

	if summary, ok := actor.Summary.Get(); ok {
		result.summaryHTML = summary
	}

	result.iconURL = firstIRI(actor.Icon)
	result.bannerURL = firstIRI(actor.Image)

	for _, protoLink := range actor.URL {
		if nil == protoLink {
			continue
		}
		iris, err := activitypub.ReturnIRIs(protoLink)
		if nil != err || 0 == len(iris) {
			continue
		}
		result.profileURL = iris[0]
		break
	}

	return result
}

func firstIRI(items []activitypub.ProtoImageOrProtoLink) string {
	for _, item := range items {
		if nil == item {
			continue
		}

		protoNode, ok := item.(activitypub.ProtoNode)
		if !ok {
			continue
		}

		iris, err := activitypub.ReturnIRIs(protoNode)
		if nil != err || 0 == len(iris) {
			continue
		}
		return iris[0]
	}
	return ""
}
