package api

import (
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// filterRoundForViewer trims participations and guess rows based on role and round status
func filterRoundForViewer(round *models.Round, viewerID string, siteAdmin, groupAdmin bool) {
	if round == nil || len(round.Participations) == 0 {
		return
	}

	revealed := dao.RoundIsRevealed(round.Status)
	canSeeAll := revealed || siteAdmin || groupAdmin

	filtered := make([]*models.RoundParticipation, 0, len(round.Participations))
	for _, p := range round.Participations {
		if p.UserID == round.PickerUserID {
			continue
		}

		if !canSeeAll && p.UserID != viewerID {
			continue
		}

		if !revealed && !siteAdmin && !groupAdmin && p.UserID != viewerID {
			p.Guesses = nil
		}

		filtered = append(filtered, p)
	}

	round.Participations = filtered
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// viewerCanSeeRoundWord reports whether the caller may see the picked word
func viewerCanSeeRoundWord(round *models.Round, viewerID string, siteAdmin, groupAdmin bool) bool {
	if round == nil || round.WordPlain == nil {
		return false
	}

	if dao.RoundIsRevealed(round.Status) {
		return true
	}

	return siteAdmin || groupAdmin || round.PickerUserID == viewerID
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// yourRoleInRound returns picker or guesser for the viewer
func yourRoleInRound(round *models.Round, viewerID string) string {
	if round.PickerUserID == viewerID {
		return "picker"
	}

	return "guesser"
}
