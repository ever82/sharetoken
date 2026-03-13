package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"sharetoken/x/crowdfunding/types"
)

// RegisterInvariants registers the crowdfunding module invariants
func (k *Keeper) RegisterInvariants(ir sdk.InvariantRegistry) {
	ir.RegisterRoute(types.ModuleName, "campaign-status",
		k.CampaignStatusInvariant())
	ir.RegisterRoute(types.ModuleName, "campaign-funding",
		k.CampaignFundingInvariant())
	ir.RegisterRoute(types.ModuleName, "campaign-backers",
		k.CampaignBackersInvariant())
	ir.RegisterRoute(types.ModuleName, "campaign-timeline",
		k.CampaignTimelineInvariant())
}

// CampaignStatusInvariant checks that campaign statuses are valid
func (k *Keeper) CampaignStatusInvariant() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var invalidStatuses []string

		k.mutex.RLock()
		defer k.mutex.RUnlock()

		for _, campaign := range k.campaigns {
			switch campaign.Status {
			case types.CampaignStatusDraft,
				types.CampaignStatusActive,
				types.CampaignStatusFunded,
				types.CampaignStatusExpired,
				types.CampaignStatusCancelled,
				types.CampaignStatusClosed:
				// Valid status
			default:
				invalidStatuses = append(invalidStatuses, fmt.Sprintf("%s:%s", campaign.ID, campaign.Status))
			}
		}

		if len(invalidStatuses) > 0 {
			return sdk.FormatInvariant(
				types.ModuleName,
				"campaign-status",
				fmt.Sprintf("found %d campaigns with invalid statuses: %v", len(invalidStatuses), invalidStatuses),
			), true
		}

		return sdk.FormatInvariant(
			types.ModuleName,
			"campaign-status",
			"all campaigns have valid status",
		), false
	}
}

// CampaignFundingInvariant checks that funding amounts are consistent
func (k *Keeper) CampaignFundingInvariant() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var invalidFunding []string

		k.mutex.RLock()
		defer k.mutex.RUnlock()

		for _, campaign := range k.campaigns {
			// Raised amount should not exceed goal for funded campaigns
			if campaign.Status == types.CampaignStatusFunded && campaign.RaisedAmount < campaign.GoalAmount {
				invalidFunding = append(invalidFunding, fmt.Sprintf("%s:funded but raised(%d)<goal(%d)", campaign.ID, campaign.RaisedAmount, campaign.GoalAmount))
			}

			// Active campaigns should have valid goal
			if campaign.Status == types.CampaignStatusActive && campaign.GoalAmount == 0 {
				invalidFunding = append(invalidFunding, fmt.Sprintf("%s:active with zero goal", campaign.ID))
			}

			// Raised amount should never be negative (implied by uint64)
			// But we check consistency with backer count
			if campaign.BackerCount > 0 && campaign.RaisedAmount == 0 {
				invalidFunding = append(invalidFunding, fmt.Sprintf("%s:has backers(%d) but zero raised", campaign.ID, campaign.BackerCount))
			}

			// Campaign type specific validations
			switch campaign.Type {
			case types.CampaignTypeInvestment:
				if campaign.EquityOffered <= 0 || campaign.EquityOffered > 100 {
					invalidFunding = append(invalidFunding, fmt.Sprintf("%s:invalid equity %f", campaign.ID, campaign.EquityOffered))
				}
			case types.CampaignTypeLending:
				if campaign.InterestRate < 0 {
					invalidFunding = append(invalidFunding, fmt.Sprintf("%s:negative interest rate", campaign.ID))
				}
			case types.CampaignTypeDonation:
				// No additional validation for donation
			default:
				invalidFunding = append(invalidFunding, fmt.Sprintf("%s:invalid type %s", campaign.ID, campaign.Type))
			}
		}

		if len(invalidFunding) > 0 {
			return sdk.FormatInvariant(
				types.ModuleName,
				"campaign-funding",
				fmt.Sprintf("found %d campaigns with invalid funding: %v", len(invalidFunding), invalidFunding),
			), true
		}

		return sdk.FormatInvariant(
			types.ModuleName,
			"campaign-funding",
			"all campaigns have valid funding",
		), false
	}
}

// CampaignBackersInvariant checks that backer data is consistent
func (k *Keeper) CampaignBackersInvariant() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var invalidBackers []string

		k.mutex.RLock()
		defer k.mutex.RUnlock()

		// Check that all backers reference valid campaigns
		for _, backer := range k.backers {
			if _, exists := k.campaigns[backer.CampaignID]; !exists {
				invalidBackers = append(invalidBackers, fmt.Sprintf("backer %s references invalid campaign %s", backer.ID, backer.CampaignID))
			}

			// Backer amount should be positive
			if backer.Amount == 0 {
				invalidBackers = append(invalidBackers, fmt.Sprintf("backer %s has zero amount", backer.ID))
			}

			// If refunded, refund amount should not exceed contribution
			if backer.Refunded && backer.RefundAmount > backer.Amount {
				invalidBackers = append(invalidBackers, fmt.Sprintf("backer %s refund(%d)>amount(%d)", backer.ID, backer.RefundAmount, backer.Amount))
			}
		}

		if len(invalidBackers) > 0 {
			return sdk.FormatInvariant(
				types.ModuleName,
				"campaign-backers",
				fmt.Sprintf("found %d invalid backers: %v", len(invalidBackers), invalidBackers),
			), true
		}

		return sdk.FormatInvariant(
			types.ModuleName,
			"campaign-backers",
			"all backers are valid",
		), false
	}
}

// CampaignTimelineInvariant checks that campaign timelines are valid
func (k *Keeper) CampaignTimelineInvariant() sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var invalidTimelines []string

		k.mutex.RLock()
		defer k.mutex.RUnlock()

		for _, campaign := range k.campaigns {
			// End time should be after start time for active/completed campaigns
			if campaign.StartTime > 0 && campaign.EndTime > 0 {
				if campaign.EndTime <= campaign.StartTime {
					invalidTimelines = append(invalidTimelines, fmt.Sprintf("%s:end(%d)<=start(%d)", campaign.ID, campaign.EndTime, campaign.StartTime))
				}
			}

			// Created at should be set
			if campaign.CreatedAt == 0 {
				invalidTimelines = append(invalidTimelines, fmt.Sprintf("%s:no created time", campaign.ID))
			}

			// Updated at should be >= created at
			if campaign.UpdatedAt > 0 && campaign.UpdatedAt < campaign.CreatedAt {
				invalidTimelines = append(invalidTimelines, fmt.Sprintf("%s:updated(%d)<created(%d)", campaign.ID, campaign.UpdatedAt, campaign.CreatedAt))
			}
		}

		if len(invalidTimelines) > 0 {
			return sdk.FormatInvariant(
				types.ModuleName,
				"campaign-timeline",
				fmt.Sprintf("found %d campaigns with invalid timelines: %v", len(invalidTimelines), invalidTimelines),
			), true
		}

		return sdk.FormatInvariant(
			types.ModuleName,
			"campaign-timeline",
			"all campaigns have valid timelines",
		), false
	}
}

// AllInvariants runs all crowdfunding invariants
func (k *Keeper) AllInvariants(ctx sdk.Context) (string, bool) {
	res, stop := k.CampaignStatusInvariant()(ctx)
	if stop {
		return res, stop
	}

	res, stop = k.CampaignFundingInvariant()(ctx)
	if stop {
		return res, stop
	}

	res, stop = k.CampaignBackersInvariant()(ctx)
	if stop {
		return res, stop
	}

	return k.CampaignTimelineInvariant()(ctx)
}
