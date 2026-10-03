package segments

import (
	"context"
	"math/bits"
	"sort"

	"github.com/SuperCoolPencil/cue/internal/domain"
)

// Chromaprint algorithm 1: 1365 samples per frame at 11025 Hz. The
// filtering/classifier delay spans ~2.6s. Trim matched boundaries inward;
// suggestions are deliberately manual-only until labeled-media validation.
const FrameSeconds = 1365.0 / 11025.0
const boundaryMargin = 3.0

type match struct{ start, end int }

func diverse(fp []uint32) bool {
	values := map[uint32]bool{}
	for _, v := range fp {
		values[v] = true
	}
	return len(values) >= 16
}

// pairMatches scans offset diagonals, allowing short fingerprint interruptions.
// At least 90% of each candidate must match; silence fails the diversity check.
func pairMatches(ctx context.Context, a, b []uint32, minFrames int) ([]match, error) {
	var matches []match
	for offset := -len(b) + minFrames; offset <= len(a)-minFrames; offset++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		i, j := offset, 0
		if i < 0 {
			j = -i
			i = 0
		}
		start, lastGood, good := -1, -1, 0
		flush := func(end int) {
			end = lastGood + 1
			if start >= 0 && end-start >= minFrames && good*10 >= (end-start)*9 && diverse(a[start:end]) {
				matches = append(matches, match{start, end})
			}
			start, lastGood, good = -1, -1, 0
		}
		for i < len(a) && j < len(b) {
			if bits.OnesCount32(a[i]^b[j]) <= 3 {
				if start < 0 {
					start = i
				}
				lastGood = i
				good++
			} else if start >= 0 && i-lastGood > 8 {
				flush(i)
			}
			i++
			j++
		}
		flush(i)
	}
	return matches, nil
}

// Detect requires the same bounded passage to match two other episodes.
// Each episode is aligned independently, allowing shifted cold opens and
// different opening clusters within a season. Conflicting candidates abstain.
func Detect(ctx context.Context, fingerprints [][]uint32, durations []int64, offsets []int64, kind string) ([][]domain.SkipSegment, error) {
	out := make([][]domain.SkipSegment, len(fingerprints))
	minFrames := 202
	for i, a := range fingerprints {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if i >= len(durations) || i >= len(offsets) {
			continue
		}
		var candidates []match
		support := make([][]match, len(fingerprints))
		for j, b := range fingerprints {
			if i == j {
				continue
			}
			m, err := pairMatches(ctx, a, b, minFrames)
			if err != nil {
				return nil, err
			}
			support[j] = m
			candidates = append(candidates, m...)
		}
		var accepted []match
		for _, candidate := range candidates {
			count := 0
			start, end := candidate.start, candidate.end
			for _, other := range support {
				for _, m := range other {
					lo, hi := max(candidate.start, m.start), min(candidate.end, m.end)
					if hi-lo >= minFrames && hi-lo >= (candidate.end-candidate.start)*8/10 {
						count++
						start = max(start, lo)
						end = min(end, hi)
						break
					}
				}
			}
			if count >= 2 && end-start >= minFrames {
				accepted = append(accepted, match{start, end})
			}
		}
		sort.Slice(accepted, func(i, j int) bool {
			if accepted[i].start == accepted[j].start {
				return accepted[i].end > accepted[j].end
			}
			return accepted[i].start < accepted[j].start
		})
		var unique []match
		for _, m := range accepted {
			if len(unique) > 0 && m.start < unique[len(unique)-1].end {
				// Multiple overlapping estimates: only retain their common interior.
				u := &unique[len(unique)-1]
				u.start = max(u.start, m.start)
				u.end = min(u.end, m.end)
			} else {
				unique = append(unique, m)
			}
		}
		for _, m := range unique {
			start := offsets[i] + int64((float64(m.start)*FrameSeconds+boundaryMargin)*1000)
			end := offsets[i] + int64(float64(m.end)*FrameSeconds*1000)
			if end-start < 20000 {
				continue
			}
			out[i] = append(out[i], domain.SkipSegment{Kind: kind, StartMs: start, EndMs: end, Origin: Version, ManualOnly: true})
		}
		out[i] = domain.ValidSkipSegments(out[i], durations[i])
	}
	return out, nil
}
