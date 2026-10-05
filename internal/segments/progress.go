package segments

import "sync"

// AnalysisProgress describes work running in this process, rather than inferring
// activity from cache files that may have been left by an interrupted run.
type AnalysisProgress struct {
	Active, Matching bool
	Completed, Total int
}

type progressKey struct{ server, user, show, season string }

var analysisProgress = struct {
	sync.RWMutex
	values map[progressKey]AnalysisProgress
}{values: make(map[progressKey]AnalysisProgress)}

func SeasonAnalysisProgress(server, user, show, season string) AnalysisProgress {
	analysisProgress.RLock()
	defer analysisProgress.RUnlock()
	return analysisProgress.values[progressKey{server, user, show, season}]
}

func setAnalysisProgress(key progressKey, progress AnalysisProgress) {
	analysisProgress.Lock()
	defer analysisProgress.Unlock()
	if !progress.Active {
		delete(analysisProgress.values, key)
	} else {
		analysisProgress.values[key] = progress
	}
}
